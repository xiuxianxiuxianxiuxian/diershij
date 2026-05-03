package repository

import (
	"context"
	"time"

	"github.com/cultivation-world/shared/types"
)

type SpellRepository interface {
	Create(ctx context.Context, spell *types.Spell) error
	GetByID(ctx context.Context, id types.SpellID) (*types.Spell, error)
	GetByName(ctx context.Context, name string) (*types.Spell, error)
	ListByType(ctx context.Context, spellType types.SpellType) ([]*types.Spell, error)
	ListByElement(ctx context.Context, element types.SpellElement) ([]*types.Spell, error)
	Update(ctx context.Context, spell *types.Spell) error
	Delete(ctx context.Context, id types.SpellID) error

	// Entity spell methods
	LearnSpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error
	GetEntitySpells(ctx context.Context, entityID types.EntityID) ([]*types.EntitySpell, error)
	GetEntitySpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) (*types.EntitySpell, error)
	UpdateSpellCastTime(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error
}

type PostgresSpellRepository struct {
	db *Database
}

func NewPostgresSpellRepository(db *Database) *PostgresSpellRepository {
	return &PostgresSpellRepository{db: db}
}

func (r *PostgresSpellRepository) Create(ctx context.Context, spell *types.Spell) error {
	if spell.CreatedAt.IsZero() {
		spell.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO spells (id, name, type, element, cost, base_damage, base_heal, duration, cooldown, description, realm_requirement, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Pool().Exec(ctx, query,
		spell.ID, spell.Name, spell.Type, spell.Element, spell.Cost, spell.BaseDamage, spell.BaseHeal,
		spell.Duration, spell.Cooldown, spell.Description, spell.RealmRequirement, spell.CreatedAt,
	)
	return err
}

func (r *PostgresSpellRepository) GetByID(ctx context.Context, id types.SpellID) (*types.Spell, error) {
	query := `
		SELECT id, name, type, element, cost, base_damage, base_heal, duration, cooldown, description, realm_requirement, created_at
		FROM spells WHERE id = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, id)
	return r.scanSpell(row)
}

func (r *PostgresSpellRepository) GetByName(ctx context.Context, name string) (*types.Spell, error) {
	query := `
		SELECT id, name, type, element, cost, base_damage, base_heal, duration, cooldown, description, realm_requirement, created_at
		FROM spells WHERE name = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, name)
	return r.scanSpell(row)
}

func (r *PostgresSpellRepository) ListByType(ctx context.Context, spellType types.SpellType) ([]*types.Spell, error) {
	query := `
		SELECT id, name, type, element, cost, base_damage, base_heal, duration, cooldown, description, realm_requirement, created_at
		FROM spells WHERE type = $1
	`
	rows, err := r.db.Pool().Query(ctx, query, spellType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*types.Spell
	for rows.Next() {
		spell, err := r.scanSpell(rows)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}
	return spells, rows.Err()
}

func (r *PostgresSpellRepository) ListByElement(ctx context.Context, element types.SpellElement) ([]*types.Spell, error) {
	query := `
		SELECT id, name, type, element, cost, base_damage, base_heal, duration, cooldown, description, realm_requirement, created_at
		FROM spells WHERE element = $1
	`
	rows, err := r.db.Pool().Query(ctx, query, element)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spells []*types.Spell
	for rows.Next() {
		spell, err := r.scanSpell(rows)
		if err != nil {
			return nil, err
		}
		spells = append(spells, spell)
	}
	return spells, rows.Err()
}

func (r *PostgresSpellRepository) Update(ctx context.Context, spell *types.Spell) error {
	query := `
		UPDATE spells SET
			name = $2, type = $3, element = $4, cost = $5, base_damage = $6, base_heal = $7,
			duration = $8, cooldown = $9, description = $10, realm_requirement = $11
		WHERE id = $1
	`
	_, err := r.db.Pool().Exec(ctx, query,
		spell.ID, spell.Name, spell.Type, spell.Element, spell.Cost, spell.BaseDamage, spell.BaseHeal,
		spell.Duration, spell.Cooldown, spell.Description, spell.RealmRequirement,
	)
	return err
}

func (r *PostgresSpellRepository) Delete(ctx context.Context, id types.SpellID) error {
	query := `DELETE FROM spells WHERE id = $1`
	_, err := r.db.Pool().Exec(ctx, query, id)
	return err
}

func (r *PostgresSpellRepository) LearnSpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error {
	query := `
		INSERT INTO entity_spells (entity_id, spell_id, learned_at, proficiency)
		VALUES ($1, $2, $3, 0)
		ON CONFLICT (entity_id, spell_id) DO NOTHING
	`
	_, err := r.db.Pool().Exec(ctx, query, entityID, spellID, time.Now())
	return err
}

func (r *PostgresSpellRepository) GetEntitySpells(ctx context.Context, entityID types.EntityID) ([]*types.EntitySpell, error) {
	query := `
		SELECT es.entity_id, es.spell_id, es.learned_at, es.proficiency, es.last_cast_at,
			   s.id, s.name, s.type, s.element, s.cost, s.base_damage, s.base_heal, s.duration, s.cooldown, s.description, s.realm_requirement
		FROM entity_spells es
		JOIN spells s ON es.spell_id = s.id
		WHERE es.entity_id = $1
	`
	rows, err := r.db.Pool().Query(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entitySpells []*types.EntitySpell
	for rows.Next() {
		es, err := r.scanEntitySpell(rows)
		if err != nil {
			return nil, err
		}
		entitySpells = append(entitySpells, es)
	}
	return entitySpells, rows.Err()
}

func (r *PostgresSpellRepository) GetEntitySpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) (*types.EntitySpell, error) {
	query := `
		SELECT es.entity_id, es.spell_id, es.learned_at, es.proficiency, es.last_cast_at,
			   s.id, s.name, s.type, s.element, s.cost, s.base_damage, s.base_heal, s.duration, s.cooldown, s.description, s.realm_requirement
		FROM entity_spells es
		JOIN spells s ON es.spell_id = s.id
		WHERE es.entity_id = $1 AND es.spell_id = $2
	`
	row := r.db.Pool().QueryRow(ctx, query, entityID, spellID)
	return r.scanEntitySpell(row)
}

func (r *PostgresSpellRepository) UpdateSpellCastTime(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error {
	query := `UPDATE entity_spells SET last_cast_at = $3 WHERE entity_id = $1 AND spell_id = $2`
	_, err := r.db.Pool().Exec(ctx, query, entityID, spellID, time.Now())
	return err
}

func (r *PostgresSpellRepository) scanSpell(scanner interface {
	Scan(dest ...interface{}) error
}) (*types.Spell, error) {
	spell := &types.Spell{}
	var createdAt time.Time

	err := scanner.Scan(
		&spell.ID, &spell.Name, &spell.Type, &spell.Element, &spell.Cost, &spell.BaseDamage, &spell.BaseHeal,
		&spell.Duration, &spell.Cooldown, &spell.Description, &spell.RealmRequirement, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	spell.CreatedAt = createdAt
	return spell, nil
}

func (r *PostgresSpellRepository) scanEntitySpell(scanner interface {
	Scan(dest ...interface{}) error
}) (*types.EntitySpell, error) {
	es := &types.EntitySpell{}
	spell := &types.Spell{}
	var learnedAt time.Time
	var lastCastAt *time.Time

	err := scanner.Scan(
		&es.EntityID, &es.SpellID, &learnedAt, &es.Proficiency, &lastCastAt,
		&spell.ID, &spell.Name, &spell.Type, &spell.Element, &spell.Cost, &spell.BaseDamage, &spell.BaseHeal,
		&spell.Duration, &spell.Cooldown, &spell.Description, &spell.RealmRequirement,
	)
	if err != nil {
		return nil, err
	}

	es.Spell = spell
	es.LearnedAt = learnedAt
	es.LastCastAt = lastCastAt

	return es, nil
}
