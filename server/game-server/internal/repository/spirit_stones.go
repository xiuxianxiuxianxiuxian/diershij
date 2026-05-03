package repository

import (
	"context"

	"github.com/cultivation-world/shared/types"
	"github.com/jackc/pgx/v5"
)

type SpiritStonesRepository struct {
	db *Database
}

func NewSpiritStonesRepository(db *Database) *SpiritStonesRepository {
	return &SpiritStonesRepository{db: db}
}

func (r *SpiritStonesRepository) GetByEntityID(ctx context.Context, entityID types.EntityID) (*types.SpiritStones, error) {
	query := `SELECT low_grade, medium_grade, high_grade, premium_grade FROM spirit_stones WHERE entity_id = $1`

	var stones types.SpiritStones
	err := r.db.Pool().QueryRow(ctx, query, entityID).Scan(
		&stones.LowGrade, &stones.MediumGrade, &stones.HighGrade, &stones.PremiumGrade,
	)
	if err == pgx.ErrNoRows {
		return &types.SpiritStones{}, nil
	}
	return &stones, err
}

func (r *SpiritStonesRepository) Upsert(ctx context.Context, entityID types.EntityID, stones *types.SpiritStones) error {
	query := `
		INSERT INTO spirit_stones (entity_id, low_grade, medium_grade, high_grade, premium_grade)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (entity_id) DO UPDATE SET
			low_grade = EXCLUDED.low_grade,
			medium_grade = EXCLUDED.medium_grade,
			high_grade = EXCLUDED.high_grade,
			premium_grade = EXCLUDED.premium_grade
	`

	_, err := r.db.Pool().Exec(ctx, query,
		entityID, stones.LowGrade, stones.MediumGrade, stones.HighGrade, stones.PremiumGrade,
	)
	return err
}
