package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Method struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Quality               string  `json:"quality"`
	RealmRequirement      string  `json:"realm_requirement"`
	ElementAffinity       string  `json:"element_affinity"`
	CultivationMultiplier float64 `json:"cultivation_multiplier"`
	BreakthroughBonus     float64 `json:"breakthrough_bonus"`
	Description           string  `json:"description"`
}

type EntityMethod struct {
	EntityID     string    `json:"entity_id"`
	MethodID     string    `json:"method_id"`
	MasteryLevel float64   `json:"mastery_level"`
	IsMainMethod bool      `json:"is_main_method"`
	LearnedAt    time.Time `json:"learned_at"`
}

type MethodRepository struct {
	db *Database
}

func NewMethodRepository(db *Database) *MethodRepository {
	return &MethodRepository{db: db}
}

func (r *MethodRepository) GetByID(ctx context.Context, id string) (*Method, error) {
	query := `SELECT id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description FROM methods WHERE id = $1`

	var m Method
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Quality, &m.RealmRequirement,
		&m.ElementAffinity, &m.CultivationMultiplier, &m.BreakthroughBonus, &m.Description,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &m, err
}

func (r *MethodRepository) GetByRealm(ctx context.Context, realm string) ([]*Method, error) {
	query := `SELECT id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description FROM methods WHERE realm_requirement = $1 ORDER BY quality, name`

	rows, err := r.db.Pool().Query(ctx, query, realm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []*Method
	for rows.Next() {
		var m Method
		if err := rows.Scan(
			&m.ID, &m.Name, &m.Quality, &m.RealmRequirement,
			&m.ElementAffinity, &m.CultivationMultiplier, &m.BreakthroughBonus, &m.Description,
		); err != nil {
			return nil, err
		}
		methods = append(methods, &m)
	}
	return methods, nil
}

func (r *MethodRepository) GetEntityMethods(ctx context.Context, entityID string) ([]*EntityMethod, error) {
	query := `SELECT entity_id, method_id, mastery_level, is_main_method, learned_at FROM entity_methods WHERE entity_id = $1`

	rows, err := r.db.Pool().Query(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []*EntityMethod
	for rows.Next() {
		var em EntityMethod
		if err := rows.Scan(&em.EntityID, &em.MethodID, &em.MasteryLevel, &em.IsMainMethod, &em.LearnedAt); err != nil {
			return nil, err
		}
		methods = append(methods, &em)
	}
	return methods, nil
}

func (r *MethodRepository) LearnMethod(ctx context.Context, entityID string, methodID string) error {
	query := `INSERT INTO entity_methods (entity_id, method_id, mastery_level, is_main_method, learned_at) VALUES ($1, $2, 0, false, NOW())`
	_, err := r.db.Pool().Exec(ctx, query, entityID, methodID)
	return err
}

func (r *MethodRepository) SetMainMethod(ctx context.Context, entityID string, methodID string) error {
	tx, err := r.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Unset all main methods for this entity
	_, err = tx.Exec(ctx, `UPDATE entity_methods SET is_main_method = false WHERE entity_id = $1`, entityID)
	if err != nil {
		return err
	}

	// Set the new main method
	_, err = tx.Exec(ctx, `UPDATE entity_methods SET is_main_method = true WHERE entity_id = $1 AND method_id = $2`, entityID, methodID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *MethodRepository) GetMainMethod(ctx context.Context, entityID string) (*EntityMethod, error) {
	query := `SELECT entity_id, method_id, mastery_level, is_main_method, learned_at FROM entity_methods WHERE entity_id = $1 AND is_main_method = true LIMIT 1`

	var em EntityMethod
	err := r.db.Pool().QueryRow(ctx, query, entityID).Scan(
		&em.EntityID, &em.MethodID, &em.MasteryLevel, &em.IsMainMethod, &em.LearnedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &em, err
}
