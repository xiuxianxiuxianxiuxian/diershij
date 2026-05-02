package repository

import (
    "context"
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
    db    *Database
    redis *redis.Client
}

func NewEntityRepository(db *Database, redis *redis.Client) *EntityRepository {
    return &EntityRepository{db: db, redis: redis}
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

    return err
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

    return &entity, nil
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
               remaining_lifespan, max_lifespan
        FROM base_attributes WHERE entity_id = $1
    `

    var attr types.Attributes
    err := r.db.Pool().QueryRow(ctx, query, entityID).Scan(
        &attr.Qi, &attr.MaxQi, &attr.SpiritualPower, &attr.MaxSpiritualPower,
        &attr.DivineSense, &attr.Comprehension, &attr.Constitution, &attr.Luck,
        &attr.CultivationProgress, &attr.AttackPower, &attr.Defense, &attr.Speed,
        &attr.MentalStability, &attr.RemainingLifespan, &attr.MaxLifespan,
    )

    if err == pgx.ErrNoRows {
        return &types.Attributes{}, nil
    }

    return &attr, err
}

func (r *EntityRepository) UpdateAttributes(ctx context.Context, entityID types.EntityID, attr *types.Attributes) error {
    query := `
        INSERT INTO base_attributes (
            entity_id, qi, max_qi, spiritual_power, max_spiritual_power,
            divine_sense, comprehension, constitution, luck, cultivation_progress,
            attack_power, defense, speed, mental_stability,
            remaining_lifespan, max_lifespan
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        ON CONFLICT (entity_id) DO UPDATE SET
            qi = EXCLUDED.qi, max_qi = EXCLUDED.max_qi,
            spiritual_power = EXCLUDED.spiritual_power, max_spiritual_power = EXCLUDED.max_spiritual_power,
            divine_sense = EXCLUDED.divine_sense, comprehension = EXCLUDED.comprehension,
            constitution = EXCLUDED.constitution, luck = EXCLUDED.luck,
            cultivation_progress = EXCLUDED.cultivation_progress, attack_power = EXCLUDED.attack_power,
            defense = EXCLUDED.defense, speed = EXCLUDED.speed,
            mental_stability = EXCLUDED.mental_stability, remaining_lifespan = EXCLUDED.remaining_lifespan,
            max_lifespan = EXCLUDED.max_lifespan
    `

    _, err := r.db.Pool().Exec(ctx, query,
        entityID, attr.Qi, attr.MaxQi, attr.SpiritualPower, attr.MaxSpiritualPower,
        attr.DivineSense, attr.Comprehension, attr.Constitution, attr.Luck,
        attr.CultivationProgress, attr.AttackPower, attr.Defense, attr.Speed,
        attr.MentalStability, attr.RemainingLifespan, attr.MaxLifespan,
    )

    return err
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
