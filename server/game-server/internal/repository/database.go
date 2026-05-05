package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/cultivation-world/shared/config"
    "github.com/cultivation-world/shared/types"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

type Database struct {
    pool *pgxpool.Pool
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }

    return &Database{pool: pool}, nil
}

func (db *Database) Close() {
    db.pool.Close()
}

func (db *Database) Pool() *pgxpool.Pool {
    return db.pool
}

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })
}

type EntityRepository struct {
    db          *Database
    redis       *redis.Client
    spiritStones *SpiritStonesRepository
    karma       *KarmaRepository
}

func NewEntityRepository(db *Database, redis *redis.Client, spiritStones *SpiritStonesRepository, karma *KarmaRepository) *EntityRepository {
    return &EntityRepository{
        db:          db,
        redis:       redis,
        spiritStones: spiritStones,
        karma:       karma,
    }
}

func (r *EntityRepository) GetByName(ctx context.Context, name string) (*types.Entity, error) {
    query := `SELECT id, entity_type, name, realm, region_id, x, y, status, created_at, updated_at 
               FROM entities WHERE name = $1`
    
    var entity types.Entity
    var pos struct {
        RegionID string
        X, Y     float64
    }
    
    err := r.db.Pool().QueryRow(ctx, query, name).Scan(
        &entity.ID, &entity.EntityType, &entity.Name, &entity.Realm,
        &pos.RegionID, &pos.X, &pos.Y, &entity.Status,
        &entity.CreatedAt, &entity.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    entity.Position = types.WorldPosition{
        RegionID: pos.RegionID,
        X:        pos.X,
        Y:        pos.Y,
    }

    attributes, _ := r.GetAttributes(ctx, entity.ID)
    if attributes != nil {
        entity.Attributes = *attributes
    }

    karma, _ := r.karma.GetByEntityID(ctx, entity.ID)
    if karma != nil {
        entity.Karma = *karma
    }

    return &entity, nil
}

func (r *EntityRepository) Create(ctx context.Context, entity *types.Entity) error {
    query := `
        INSERT INTO entities (id, entity_type, name, realm, region_id, x, y, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

    _, err := r.db.Pool().Exec(ctx, query,
        entity.ID, entity.EntityType, entity.Name, entity.Realm,
        entity.Position.RegionID, entity.Position.X, entity.Position.Y,
        entity.Status, entity.CreatedAt, entity.UpdatedAt,
    )

    if err != nil {
        return err
    }

    // 初始化基础属性（含灵石）
    if err := r.UpdateAttributes(ctx, entity.ID, &entity.Attributes); err != nil {
        return err
    }

    // 初始化业力
    return r.karma.Upsert(ctx, entity.ID, &entity.Karma)
}

func (r *EntityRepository) GetByID(ctx context.Context, id types.EntityID) (*types.Entity, error) {
    query := `
        SELECT id, entity_type, name, realm, region_id, x, y, status, created_at, updated_at
        FROM entities WHERE id = $1
    `

    var entity types.Entity
    var pos struct {
        RegionID string
        X, Y     float64
    }

    err := r.db.Pool().QueryRow(ctx, query, id).Scan(
        &entity.ID, &entity.EntityType, &entity.Name, &entity.Realm,
        &pos.RegionID, &pos.X, &pos.Y, &entity.Status,
        &entity.CreatedAt, &entity.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    entity.Position = types.WorldPosition{
        RegionID: pos.RegionID,
        X:        pos.X,
        Y:        pos.Y,
    }

    // 加载业力
    karma, _ := r.karma.GetByEntityID(ctx, entity.ID)
    if karma != nil {
        entity.Karma = *karma
    }

    // 加载属性（包含灵石）
    attributes, _ := r.GetAttributes(ctx, entity.ID)
    if attributes != nil {
        entity.Attributes = *attributes
    }

    return &entity, nil
}

func (r *EntityRepository) SetPasswordHash(ctx context.Context, entityID types.EntityID, hash string) error {
	_, err := r.db.Pool().Exec(ctx,
		"UPDATE entities SET password_hash = $2 WHERE id = $1",
		entityID, hash,
	)
	return err
}

func (r *EntityRepository) GetPasswordHash(ctx context.Context, entityID types.EntityID) (string, error) {
	var hash string
	err := r.db.Pool().QueryRow(ctx,
		"SELECT password_hash FROM entities WHERE id = $1", entityID,
	).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (r *EntityRepository) Update(ctx context.Context, entity *types.Entity) error {
    query := `
        UPDATE entities SET
            name = $2, realm = $3, region_id = $4, x = $5, y = $6,
            status = $7, updated_at = $8
        WHERE id = $1
    `

    _, err := r.db.Pool().Exec(ctx, query,
        entity.ID, entity.Name, entity.Realm,
        entity.Position.RegionID, entity.Position.X, entity.Position.Y,
        entity.Status, entity.UpdatedAt,
    )

    return err
}

func (r *EntityRepository) GetAttributes(ctx context.Context, entityID types.EntityID) (*types.Attributes, error) {
	query := `
		SELECT qi, max_qi, spiritual_power, max_spiritual_power, divine_sense,
		       comprehension, constitution, luck, cultivation_progress,
		       attack_power, defense, speed, mental_stability,
		       remaining_lifespan, max_lifespan,
		       crit_rate, crit_damage, dodge_rate, hit_rate, penetration, damage_reduction,
		       alchemy_level, artificing_level, mining_level, herb_level,
		       talisman_level, formation_level, fire_control, beast_taming,
		       reputation, sect_contribution, dao_heart, enlightenment,
		       property_value, business_income, root_purity, poison_level, curse_level,
		       root_awakened, mutated_root,
		       appearance, charisma, obsession_count, inner_demon_resistance, aging_penalty,
		       bloodline, bloodline_purity, physique, physique_awakened, destiny, world_favor,
		       law_resonance, domain_power, domain_range, law_suppression,
		       dao_seed_type, dao_seed_level, dao_seed_growth, dao_marks, dao_heart_comprehension, destiny_path,
		       injuries, buffs, debuffs, laws, faction_standings,
		       real_estate, disciple_ids, sworn_siblings, enemies, lovers,
		       mentor_id
		FROM base_attributes WHERE entity_id = $1
	`

	var attr types.Attributes
	var mentorID *string
	var injuriesJSON, buffsJSON, debuffsJSON json.RawMessage
	var lawsJSON, factionJSON json.RawMessage
	var realEstateJSON, discipleIDsJSON, swornJSON, enemiesJSON, loversJSON json.RawMessage

	err := r.db.Pool().QueryRow(ctx, query, entityID).Scan(
		// Basic cultivation (15)
		&attr.Qi, &attr.MaxQi, &attr.SpiritualPower, &attr.MaxSpiritualPower,
		&attr.DivineSense, &attr.Comprehension, &attr.Constitution, &attr.Luck,
		&attr.CultivationProgress, &attr.AttackPower, &attr.Defense, &attr.Speed,
		&attr.MentalStability, &attr.RemainingLifespan, &attr.MaxLifespan,
		// Combat (6)
		&attr.CritRate, &attr.CritDamage, &attr.DodgeRate, &attr.HitRate,
		&attr.Penetration, &attr.DamageReduction,
		// Life skills (8)
		&attr.AlchemyLevel, &attr.ArtificingLevel, &attr.MiningSkill, &attr.HerbKnowledge,
		&attr.TalismanSkill, &attr.FormationLevel, &attr.FireControl, &attr.BeastTaming,
		// Social + Mental (4)
		&attr.Reputation, &attr.SectContribution, &attr.DaoHeart, &attr.Enlightenment,
		// Wealth + special scalars (5)
		&attr.PropertyValue, &attr.BusinessIncome, &attr.RootPurity, &attr.PoisonLevel, &attr.CurseLevel,
		// Spiritual roots flags (2)
		&attr.RootAwakened, &attr.MutatedRoot,
		// Appearance (2)
		&attr.Appearance, &attr.Charisma,
		// Mental extended (2)
		&attr.ObsessionCount, &attr.InnerDemonResistance,
		// Lifespan extended (1)
		&attr.AgingPenalty,
		// Special attributes (6)
		&attr.Bloodline, &attr.BloodlinePurity, &attr.Physique, &attr.PhysiqueAwakened,
		&attr.Destiny, &attr.WorldFavor,
		// Law attributes (4)
		&attr.LawResonance, &attr.DomainPower, &attr.DomainRange, &attr.LawSuppression,
		// Dao attributes (6)
		&attr.DaoSeedType, &attr.DaoSeedLevel, &attr.DaoSeedGrowth, &attr.DaoMarks,
		&attr.DaoHeartComprehension, &attr.DestinyPath,
		// JSONB complex types (10)
		&injuriesJSON, &buffsJSON, &debuffsJSON, &lawsJSON, &factionJSON,
		&realEstateJSON, &discipleIDsJSON, &swornJSON, &enemiesJSON, &loversJSON,
		// Mentor
		&mentorID,
	)

	if err == pgx.ErrNoRows {
		return &types.Attributes{}, nil
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSONB fields
	if len(injuriesJSON) > 0 {
		json.Unmarshal(injuriesJSON, &attr.Injuries)
	}
	if len(buffsJSON) > 0 {
		json.Unmarshal(buffsJSON, &attr.Buffs)
	}
	if len(debuffsJSON) > 0 {
		json.Unmarshal(debuffsJSON, &attr.Debuffs)
	}
	if len(lawsJSON) > 0 {
		json.Unmarshal(lawsJSON, &attr.Laws)
	}
	if len(factionJSON) > 0 {
		json.Unmarshal(factionJSON, &attr.FactionStandings)
	}
	if len(realEstateJSON) > 0 {
		json.Unmarshal(realEstateJSON, &attr.RealEstate)
	}
	if len(discipleIDsJSON) > 0 {
		json.Unmarshal(discipleIDsJSON, &attr.DiscipleIDs)
	}
	if len(swornJSON) > 0 {
		json.Unmarshal(swornJSON, &attr.SwornSiblings)
	}
	if len(enemiesJSON) > 0 {
		json.Unmarshal(enemiesJSON, &attr.Enemies)
	}
	if len(loversJSON) > 0 {
		json.Unmarshal(loversJSON, &attr.Lovers)
	}
	if mentorID != nil {
		attr.MentorID = *mentorID
	}

	// 加载灵石
	stones, err := r.spiritStones.GetByEntityID(ctx, entityID)
	if err != nil {
		stones = &types.SpiritStones{}
	}
	attr.SpiritStones = *stones

	// 加载灵根
	roots, err := r.GetSpiritualRoots(ctx, entityID)
	if err == nil && len(roots) > 0 {
		attr.SpiritualRoots = roots
	}

	return &attr, nil
}

func (r *EntityRepository) UpdateAttributes(ctx context.Context, entityID types.EntityID, attr *types.Attributes) error {
	// Pre-marshal JSONB fields
	injuriesBytes, _ := json.Marshal(attr.Injuries)
	buffsBytes, _ := json.Marshal(attr.Buffs)
	debuffsBytes, _ := json.Marshal(attr.Debuffs)
	lawsBytes, _ := json.Marshal(attr.Laws)
	factionBytes, _ := json.Marshal(attr.FactionStandings)
	realEstateBytes, _ := json.Marshal(attr.RealEstate)
	discipleIDsBytes, _ := json.Marshal(attr.DiscipleIDs)
	swornBytes, _ := json.Marshal(attr.SwornSiblings)
	enemiesBytes, _ := json.Marshal(attr.Enemies)
	loversBytes, _ := json.Marshal(attr.Lovers)

	var mentorID *string
	if attr.MentorID != "" {
		mentorID = &attr.MentorID
	}

	query := `
		INSERT INTO base_attributes (
		    entity_id, qi, max_qi, spiritual_power, max_spiritual_power,
		    divine_sense, comprehension, constitution, luck, cultivation_progress,
		    attack_power, defense, speed, mental_stability,
		    remaining_lifespan, max_lifespan,
		    crit_rate, crit_damage, dodge_rate, hit_rate, penetration, damage_reduction,
		    alchemy_level, artificing_level, mining_level, herb_level,
		    talisman_level, formation_level, fire_control, beast_taming,
		    reputation, sect_contribution, dao_heart, enlightenment,
		    property_value, business_income, root_purity, poison_level, curse_level,
		    root_awakened, mutated_root,
		    appearance, charisma, obsession_count, inner_demon_resistance, aging_penalty,
		    bloodline, bloodline_purity, physique, physique_awakened, destiny, world_favor,
		    law_resonance, domain_power, domain_range, law_suppression,
		    dao_seed_type, dao_seed_level, dao_seed_growth, dao_marks, dao_heart_comprehension, destiny_path,
		    injuries, buffs, debuffs, laws, faction_standings,
		    real_estate, disciple_ids, sworn_siblings, enemies, lovers,
		    mentor_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		         $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
		         $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
		         $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57, $58,
		         $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73)
		ON CONFLICT (entity_id) DO UPDATE SET
		    qi = EXCLUDED.qi, max_qi = EXCLUDED.max_qi,
		    spiritual_power = EXCLUDED.spiritual_power, max_spiritual_power = EXCLUDED.max_spiritual_power,
		    divine_sense = EXCLUDED.divine_sense, comprehension = EXCLUDED.comprehension,
		    constitution = EXCLUDED.constitution, luck = EXCLUDED.luck,
		    cultivation_progress = EXCLUDED.cultivation_progress, attack_power = EXCLUDED.attack_power,
		    defense = EXCLUDED.defense, speed = EXCLUDED.speed,
		    mental_stability = EXCLUDED.mental_stability, remaining_lifespan = EXCLUDED.remaining_lifespan,
		    max_lifespan = EXCLUDED.max_lifespan,
		    crit_rate = EXCLUDED.crit_rate, crit_damage = EXCLUDED.crit_damage,
		    dodge_rate = EXCLUDED.dodge_rate, hit_rate = EXCLUDED.hit_rate,
		    penetration = EXCLUDED.penetration, damage_reduction = EXCLUDED.damage_reduction,
		    alchemy_level = EXCLUDED.alchemy_level, artificing_level = EXCLUDED.artificing_level,
		    mining_level = EXCLUDED.mining_level, herb_level = EXCLUDED.herb_level,
		    talisman_level = EXCLUDED.talisman_level, formation_level = EXCLUDED.formation_level,
		    fire_control = EXCLUDED.fire_control, beast_taming = EXCLUDED.beast_taming,
		    reputation = EXCLUDED.reputation, sect_contribution = EXCLUDED.sect_contribution,
		    dao_heart = EXCLUDED.dao_heart, enlightenment = EXCLUDED.enlightenment,
		    property_value = EXCLUDED.property_value, business_income = EXCLUDED.business_income,
		    root_purity = EXCLUDED.root_purity, poison_level = EXCLUDED.poison_level,
		    curse_level = EXCLUDED.curse_level,
		    root_awakened = EXCLUDED.root_awakened, mutated_root = EXCLUDED.mutated_root,
		    appearance = EXCLUDED.appearance, charisma = EXCLUDED.charisma,
		    obsession_count = EXCLUDED.obsession_count, inner_demon_resistance = EXCLUDED.inner_demon_resistance,
		    aging_penalty = EXCLUDED.aging_penalty,
		    bloodline = EXCLUDED.bloodline, bloodline_purity = EXCLUDED.bloodline_purity,
		    physique = EXCLUDED.physique, physique_awakened = EXCLUDED.physique_awakened,
		    destiny = EXCLUDED.destiny, world_favor = EXCLUDED.world_favor,
		    law_resonance = EXCLUDED.law_resonance, domain_power = EXCLUDED.domain_power,
		    domain_range = EXCLUDED.domain_range, law_suppression = EXCLUDED.law_suppression,
		    dao_seed_type = EXCLUDED.dao_seed_type, dao_seed_level = EXCLUDED.dao_seed_level,
		    dao_seed_growth = EXCLUDED.dao_seed_growth, dao_marks = EXCLUDED.dao_marks,
		    dao_heart_comprehension = EXCLUDED.dao_heart_comprehension, destiny_path = EXCLUDED.destiny_path,
		    injuries = EXCLUDED.injuries, buffs = EXCLUDED.buffs, debuffs = EXCLUDED.debuffs,
		    laws = EXCLUDED.laws, faction_standings = EXCLUDED.faction_standings,
		    real_estate = EXCLUDED.real_estate, disciple_ids = EXCLUDED.disciple_ids,
		    sworn_siblings = EXCLUDED.sworn_siblings, enemies = EXCLUDED.enemies, lovers = EXCLUDED.lovers,
		    mentor_id = EXCLUDED.mentor_id
	`

	_, err := r.db.Pool().Exec(ctx, query,
		entityID, attr.Qi, attr.MaxQi, attr.SpiritualPower, attr.MaxSpiritualPower,
		attr.DivineSense, attr.Comprehension, attr.Constitution, attr.Luck,
		attr.CultivationProgress, attr.AttackPower, attr.Defense, attr.Speed,
		attr.MentalStability, attr.RemainingLifespan, attr.MaxLifespan,
		attr.CritRate, attr.CritDamage, attr.DodgeRate, attr.HitRate,
		attr.Penetration, attr.DamageReduction,
		attr.AlchemyLevel, attr.ArtificingLevel, attr.MiningSkill, attr.HerbKnowledge,
		attr.TalismanSkill, attr.FormationLevel, attr.FireControl, attr.BeastTaming,
		attr.Reputation, attr.SectContribution, attr.DaoHeart, attr.Enlightenment,
		attr.PropertyValue, attr.BusinessIncome, attr.RootPurity, attr.PoisonLevel, attr.CurseLevel,
		attr.RootAwakened, attr.MutatedRoot,
		attr.Appearance, attr.Charisma, attr.ObsessionCount, attr.InnerDemonResistance, attr.AgingPenalty,
		attr.Bloodline, attr.BloodlinePurity, attr.Physique, attr.PhysiqueAwakened,
		attr.Destiny, attr.WorldFavor,
		attr.LawResonance, attr.DomainPower, attr.DomainRange, attr.LawSuppression,
		attr.DaoSeedType, attr.DaoSeedLevel, attr.DaoSeedGrowth, attr.DaoMarks,
		attr.DaoHeartComprehension, attr.DestinyPath,
		injuriesBytes, buffsBytes, debuffsBytes, lawsBytes, factionBytes,
		realEstateBytes, discipleIDsBytes, swornBytes, enemiesBytes, loversBytes,
		mentorID,
	)

	if err != nil {
		return err
	}
	return r.spiritStones.Upsert(ctx, entityID, &attr.SpiritStones)
}

func (r *EntityRepository) UpdateSpiritualRoots(ctx context.Context, entityID types.EntityID, roots []types.SpiritualRoot) error {
	if len(roots) == 0 {
		_, err := r.db.Pool().Exec(ctx,
			"DELETE FROM spiritual_roots WHERE entity_id = $1", entityID)
		return err
	}

	tx, err := r.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Replace all existing roots
	if _, err := tx.Exec(ctx, "DELETE FROM spiritual_roots WHERE entity_id = $1", entityID); err != nil {
		return err
	}

	for _, root := range roots {
		if _, err := tx.Exec(ctx,
			"INSERT INTO spiritual_roots (entity_id, element, purity) VALUES ($1, $2, $3)",
			entityID, root.Element, root.Purity,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *EntityRepository) GetSpiritualRoots(ctx context.Context, entityID types.EntityID) ([]types.SpiritualRoot, error) {
	rows, err := r.db.Pool().Query(ctx,
		"SELECT element, purity FROM spiritual_roots WHERE entity_id = $1 ORDER BY purity DESC",
		entityID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roots []types.SpiritualRoot
	for rows.Next() {
		var root types.SpiritualRoot
		if err := rows.Scan(&root.Element, &root.Purity); err != nil {
			return nil, err
		}
		roots = append(roots, root)
	}

	return roots, rows.Err()
}

func (r *EntityRepository) UpdateKarma(ctx context.Context, entityID types.EntityID, karma *types.Karma) error {
    return r.karma.Upsert(ctx, entityID, karma)
}

func (r *EntityRepository) CacheEntity(ctx context.Context, entity *types.Entity) error {
    return r.redis.Set(ctx, "entity:"+string(entity.ID), entity, 5*time.Minute).Err()
}

func (r *EntityRepository) GetCachedEntity(ctx context.Context, id types.EntityID) (*types.Entity, error) {
    var entity types.Entity
    err := r.redis.Get(ctx, "entity:"+string(id)).Scan(&entity)
    if err != nil {
        return nil, err
    }
    return &entity, nil
}
