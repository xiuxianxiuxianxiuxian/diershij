package repository

import (
	"context"

	"github.com/cultivation-world/shared/types"
	"github.com/jackc/pgx/v5"
)

type KarmaRepository struct {
	db *Database
}

func NewKarmaRepository(db *Database) *KarmaRepository {
	return &KarmaRepository{db: db}
}

func (r *KarmaRepository) GetByEntityID(ctx context.Context, entityID types.EntityID) (*types.Karma, error) {
	query := `SELECT karma_value, merit, karmic_debt, heavenly_mark FROM karma_attributes WHERE entity_id = $1`

	var karma types.Karma
	err := r.db.Pool().QueryRow(ctx, query, entityID).Scan(
		&karma.KarmaValue, &karma.Merit, &karma.KarmicDebt, &karma.HeavenlyMark,
	)
	if err == pgx.ErrNoRows {
		return &types.Karma{}, nil
	}
	return &karma, err
}

func (r *KarmaRepository) Upsert(ctx context.Context, entityID types.EntityID, karma *types.Karma) error {
	query := `
		INSERT INTO karma_attributes (entity_id, karma_value, merit, karmic_debt, heavenly_mark)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (entity_id) DO UPDATE SET
			karma_value = EXCLUDED.karma_value,
			merit = EXCLUDED.merit,
			karmic_debt = EXCLUDED.karmic_debt,
			heavenly_mark = EXCLUDED.heavenly_mark
	`

	_, err := r.db.Pool().Exec(ctx, query,
		entityID, karma.KarmaValue, karma.Merit, karma.KarmicDebt, karma.HeavenlyMark,
	)
	return err
}
