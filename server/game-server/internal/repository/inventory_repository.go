package repository

import (
	"context"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InventoryRepository interface {
	GetByEntityID(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error)
	GetItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID) (*types.InventoryItem, error)
	AddItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	RemoveItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	EquipItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, slot string) error
	UnequipItem(ctx context.Context, entityID types.EntityID, slot string) error
	GetEquippedItems(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error)
}

type PostgresInventoryRepository struct {
	db *Database
}

func NewPostgresInventoryRepository(db *Database) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db: db}
}

func (r *PostgresInventoryRepository) GetByEntityID(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error) {
	query := `
		SELECT i.id, i.entity_id, i.item_id, i.quantity, i.equipped, i.slot, i.durability, i.bound, i.acquired_at,
			   items.id, items.name, items.type, items.rarity, items.description, items.attributes,
			   items.stackable, items.max_stack, items.usable, items.level_requirement, items.realm_requirement
		FROM inventory i
		JOIN items ON i.item_id = items.id
		WHERE i.entity_id = $1
	`
	rows, err := r.db.Pool().Query(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*types.InventoryItem
	for rows.Next() {
		item, err := r.scanInventoryItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresInventoryRepository) GetItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID) (*types.InventoryItem, error) {
	query := `
		SELECT i.id, i.entity_id, i.item_id, i.quantity, i.equipped, i.slot, i.durability, i.bound, i.acquired_at,
			   items.id, items.name, items.type, items.rarity, items.description, items.attributes,
			   items.stackable, items.max_stack, items.usable, items.level_requirement, items.realm_requirement
		FROM inventory i
		JOIN items ON i.item_id = items.id
		WHERE i.entity_id = $1 AND i.item_id = $2
	`
	row := r.db.Pool().QueryRow(ctx, query, entityID, itemID)
	return r.scanInventoryItem(row)
}

func (r *PostgresInventoryRepository) AddItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error {
	// Check if item already exists in inventory
	existing, err := r.GetItem(ctx, entityID, itemID)
	if err == nil && existing != nil {
		// Update quantity
		query := `UPDATE inventory SET quantity = quantity + $3 WHERE entity_id = $1 AND item_id = $2`
		_, err = r.db.Pool().Exec(ctx, query, entityID, itemID, quantity)
		return err
	}

	// Insert new item
	query := `
		INSERT INTO inventory (id, entity_id, item_id, quantity, equipped, bound, acquired_at)
		VALUES ($1, $2, $3, $4, false, false, $5)
	`
	_, err = r.db.Pool().Exec(ctx, query, uuid.New().String(), entityID, itemID, quantity, time.Now())
	return err
}

func (r *PostgresInventoryRepository) RemoveItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error {
	// Get current quantity
	existing, err := r.GetItem(ctx, entityID, itemID)
	if err != nil {
		return err
	}

	if existing.Quantity < quantity {
		return pgx.ErrNoRows // Not enough items
	}

	if existing.Quantity == quantity {
		// Remove item completely
		query := `DELETE FROM inventory WHERE entity_id = $1 AND item_id = $2`
		_, err = r.db.Pool().Exec(ctx, query, entityID, itemID)
		return err
	}

	// Update quantity
	query := `UPDATE inventory SET quantity = quantity - $3 WHERE entity_id = $1 AND item_id = $2`
	_, err = r.db.Pool().Exec(ctx, query, entityID, itemID, quantity)
	return err
}

func (r *PostgresInventoryRepository) EquipItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, slot string) error {
	// First unequip any item in the same slot
	query := `UPDATE inventory SET equipped = false, slot = NULL WHERE entity_id = $1 AND slot = $2`
	_, err := r.db.Pool().Exec(ctx, query, entityID, slot)
	if err != nil {
		return err
	}

	// Equip the new item
	query = `UPDATE inventory SET equipped = true, slot = $3 WHERE entity_id = $1 AND item_id = $2`
	_, err = r.db.Pool().Exec(ctx, query, entityID, itemID, slot)
	return err
}

func (r *PostgresInventoryRepository) UnequipItem(ctx context.Context, entityID types.EntityID, slot string) error {
	query := `UPDATE inventory SET equipped = false, slot = NULL WHERE entity_id = $1 AND slot = $2`
	_, err := r.db.Pool().Exec(ctx, query, entityID, slot)
	return err
}

func (r *PostgresInventoryRepository) GetEquippedItems(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error) {
	query := `
		SELECT i.id, i.entity_id, i.item_id, i.quantity, i.equipped, i.slot, i.durability, i.bound, i.acquired_at,
			   items.id, items.name, items.type, items.rarity, items.description, items.attributes,
			   items.stackable, items.max_stack, items.usable, items.level_requirement, items.realm_requirement
		FROM inventory i
		JOIN items ON i.item_id = items.id
		WHERE i.entity_id = $1 AND i.equipped = true
	`
	rows, err := r.db.Pool().Query(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*types.InventoryItem
	for rows.Next() {
		item, err := r.scanInventoryItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresInventoryRepository) scanInventoryItem(scanner interface {
	Scan(dest ...interface{}) error
}) (*types.InventoryItem, error) {
	invItem := &types.InventoryItem{}
	item := &types.Item{}
	var acquiredAt time.Time

	err := scanner.Scan(
		&invItem.ID, &invItem.EntityID, &invItem.ItemID, &invItem.Quantity, &invItem.Equipped, &invItem.Slot,
		&invItem.Durability, &invItem.Bound, &acquiredAt,
		&item.ID, &item.Name, &item.Type, &item.Rarity, &item.Description, &item.Attributes,
		&item.Stackable, &item.MaxStack, &item.Usable, &item.LevelRequirement, &item.RealmRequirement,
	)
	if err != nil {
		return nil, err
	}

	invItem.Item = item
	invItem.AcquiredAt = acquiredAt.Unix()
	return invItem, nil
}
