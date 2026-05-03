package repository

import (
	"context"
	"time"
)

type Recipe struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	Difficulty        int                    `json:"difficulty"`
	Description       string                 `json:"description"`
	Materials         map[string]interface{} `json:"materials"`
	ResultItemID      string                 `json:"result_item_id"`
	ResultQuantity    int                    `json:"result_quantity"`
	SkillLevelRequired int                   `json:"skill_level_required"`
	CreatedAt         time.Time              `json:"created_at"`
}

type RecipeRepository struct {
	db *Database
}

func NewRecipeRepository(db *Database) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) GetByID(ctx context.Context, id string) (*Recipe, error) {
	query := `SELECT id, name, type, difficulty, description, materials, result_item_id, result_quantity, skill_level_required, created_at FROM recipes WHERE id = $1`

	var rec Recipe
	rec.Materials = make(map[string]interface{})
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&rec.ID, &rec.Name, &rec.Type, &rec.Difficulty, &rec.Description,
		&rec.Materials, &rec.ResultItemID, &rec.ResultQuantity,
		&rec.SkillLevelRequired, &rec.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecipeRepository) GetByType(ctx context.Context, recipeType string) ([]*Recipe, error) {
	query := `SELECT id, name, type, difficulty, description, materials, result_item_id, result_quantity, skill_level_required, created_at FROM recipes WHERE type = $1`

	rows, err := r.db.Pool().Query(ctx, query, recipeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*Recipe
	for rows.Next() {
		var rec Recipe
		rec.Materials = make(map[string]interface{})
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Type, &rec.Difficulty, &rec.Description,
			&rec.Materials, &rec.ResultItemID, &rec.ResultQuantity,
			&rec.SkillLevelRequired, &rec.CreatedAt); err != nil {
			return nil, err
		}
		recipes = append(recipes, &rec)
	}
	return recipes, nil
}
