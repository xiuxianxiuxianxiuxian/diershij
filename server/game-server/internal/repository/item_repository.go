package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/google/uuid"
)

type ItemRepository interface {
	Create(ctx context.Context, item *types.Item) error
	GetByID(ctx context.Context, id types.ItemID) (*types.Item, error)
	GetByName(ctx context.Context, name string) (*types.Item, error)
	ListByType(ctx context.Context, itemType types.ItemType) ([]*types.Item, error)
	Update(ctx context.Context, item *types.Item) error
	Delete(ctx context.Context, id types.ItemID) error
}

type PostgresItemRepository struct {
	db *Database
}

func NewPostgresItemRepository(db *Database) *PostgresItemRepository {
	return &PostgresItemRepository{db: db}
}

func (r *PostgresItemRepository) Create(ctx context.Context, item *types.Item) error {
	if item.ID == "" {
		item.ID = types.ItemID(uuid.New().String())
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}

	attributesJSON, err := json.Marshal(item.Attributes)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO items (id, name, type, rarity, description, attributes, stackable, max_stack, usable, level_requirement, realm_requirement, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = r.db.Pool().Exec(ctx, query,
		item.ID, item.Name, item.Type, item.Rarity, item.Description, attributesJSON,
		item.Stackable, item.MaxStack, item.Usable, item.LevelRequirement, item.RealmRequirement, item.CreatedAt,
	)
	return err
}

func (r *PostgresItemRepository) GetByID(ctx context.Context, id types.ItemID) (*types.Item, error) {
	query := `
		SELECT id, name, type, rarity, description, attributes, stackable, max_stack, usable, level_requirement, realm_requirement, created_at
		FROM items WHERE id = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, id)
	return r.scanItem(row)
}

func (r *PostgresItemRepository) GetByName(ctx context.Context, name string) (*types.Item, error) {
	query := `
		SELECT id, name, type, rarity, description, attributes, stackable, max_stack, usable, level_requirement, realm_requirement, created_at
		FROM items WHERE name = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, name)
	return r.scanItem(row)
}

func (r *PostgresItemRepository) ListByType(ctx context.Context, itemType types.ItemType) ([]*types.Item, error) {
	query := `
		SELECT id, name, type, rarity, description, attributes, stackable, max_stack, usable, level_requirement, realm_requirement, created_at
		FROM items WHERE type = $1
	`
	rows, err := r.db.Pool().Query(ctx, query, itemType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*types.Item
	for rows.Next() {
		item, err := r.scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresItemRepository) Update(ctx context.Context, item *types.Item) error {
	attributesJSON, err := json.Marshal(item.Attributes)
	if err != nil {
		return err
	}

	query := `
		UPDATE items SET
			name = $2, type = $3, rarity = $4, description = $5, attributes = $6,
			stackable = $7, max_stack = $8, usable = $9, level_requirement = $10, realm_requirement = $11
		WHERE id = $1
	`
	_, err = r.db.Pool().Exec(ctx, query,
		item.ID, item.Name, item.Type, item.Rarity, item.Description, attributesJSON,
		item.Stackable, item.MaxStack, item.Usable, item.LevelRequirement, item.RealmRequirement,
	)
	return err
}

func (r *PostgresItemRepository) Delete(ctx context.Context, id types.ItemID) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.Pool().Exec(ctx, query, id)
	return err
}

func (r *PostgresItemRepository) scanItem(scanner interface {
	Scan(dest ...interface{}) error
}) (*types.Item, error) {
	item := &types.Item{}
	var attributesJSON []byte
	var createdAt time.Time

	err := scanner.Scan(
		&item.ID, &item.Name, &item.Type, &item.Rarity, &item.Description, &attributesJSON,
		&item.Stackable, &item.MaxStack, &item.Usable, &item.LevelRequirement, &item.RealmRequirement, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	if len(attributesJSON) > 0 {
		json.Unmarshal(attributesJSON, &item.Attributes)
	}
	item.CreatedAt = createdAt

	return item, nil
}
