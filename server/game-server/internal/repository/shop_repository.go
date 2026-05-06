package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

// ShopInfo represents a NPC shop.
type ShopInfo struct {
	ID          string
	Name        string
	Description string
	RegionID    string
	ShopType    string
	NPCOwner    string
	MarkupRate  float64
	BuyRate     float64
}

// ShopItem represents a single item listing in a shop.
type ShopItem struct {
	ID           string
	ShopID       string
	ItemName     string
	ItemType     string
	Rarity       int
	Price        int64
	Quantity     int   // -1 = unlimited
	RefreshHours int
	MinRealm     string
}

// Auction represents a player auction listing.
type Auction struct {
	ID        string
	SellerID  string
	ItemID    string
	ItemName  string
	Quantity  int
	Price     int64
	Deposit   int64
	Status    string // active, sold, cancelled
	CreatedAt time.Time
	ExpiresAt *time.Time
	BuyerID   *string
	SoldAt    *time.Time
}

// ShopRepository handles shop and auction database operations.
type ShopRepository struct {
	db *Database
}

func NewShopRepository(db *Database) *ShopRepository {
	return &ShopRepository{db: db}
}

// ── Shops ──

func (r *ShopRepository) GetShopByID(ctx context.Context, shopID string) (*ShopInfo, error) {
	query := `SELECT id, name, COALESCE(description,''), COALESCE(region_id,''), shop_type, COALESCE(npc_owner,''), markup_rate, buy_rate FROM shops WHERE id = $1`
	row := r.db.Pool().QueryRow(ctx, query, shopID)
	var s ShopInfo
	err := row.Scan(&s.ID, &s.Name, &s.Description, &s.RegionID, &s.ShopType, &s.NPCOwner, &s.MarkupRate, &s.BuyRate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *ShopRepository) ListShops(ctx context.Context) ([]*ShopInfo, error) {
	query := `SELECT id, name, COALESCE(description,''), COALESCE(region_id,''), shop_type, COALESCE(npc_owner,''), markup_rate, buy_rate FROM shops ORDER BY name`
	rows, err := r.db.Pool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []*ShopInfo
	for rows.Next() {
		var s ShopInfo
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.RegionID, &s.ShopType, &s.NPCOwner, &s.MarkupRate, &s.BuyRate); err != nil {
			return nil, err
		}
		shops = append(shops, &s)
	}
	return shops, rows.Err()
}

func (r *ShopRepository) ListShopsByRegion(ctx context.Context, regionID string) ([]*ShopInfo, error) {
	query := `SELECT id, name, COALESCE(description,''), COALESCE(region_id,''), shop_type, COALESCE(npc_owner,''), markup_rate, buy_rate FROM shops WHERE region_id = $1 ORDER BY name`
	rows, err := r.db.Pool().Query(ctx, query, regionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []*ShopInfo
	for rows.Next() {
		var s ShopInfo
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.RegionID, &s.ShopType, &s.NPCOwner, &s.MarkupRate, &s.BuyRate); err != nil {
			return nil, err
		}
		shops = append(shops, &s)
	}
	return shops, rows.Err()
}

// ── Shop Inventory ──

func (r *ShopRepository) GetShopInventory(ctx context.Context, shopID string) ([]*ShopItem, error) {
	query := `SELECT id, shop_id, item_name, item_type, rarity, price, quantity, refresh_hours, min_realm FROM shop_inventory WHERE shop_id = $1 ORDER BY rarity, item_name`
	rows, err := r.db.Pool().Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ShopItem
	for rows.Next() {
		var si ShopItem
		if err := rows.Scan(&si.ID, &si.ShopID, &si.ItemName, &si.ItemType, &si.Rarity, &si.Price, &si.Quantity, &si.RefreshHours, &si.MinRealm); err != nil {
			return nil, err
		}
		items = append(items, &si)
	}
	return items, rows.Err()
}

func (r *ShopRepository) GetShopItemByName(ctx context.Context, shopID string, itemName string) (*ShopItem, error) {
	query := `SELECT id, shop_id, item_name, item_type, rarity, price, quantity, refresh_hours, min_realm FROM shop_inventory WHERE shop_id = $1 AND item_name = $2`
	row := r.db.Pool().QueryRow(ctx, query, shopID, itemName)
	var si ShopItem
	err := row.Scan(&si.ID, &si.ShopID, &si.ItemName, &si.ItemType, &si.Rarity, &si.Price, &si.Quantity, &si.RefreshHours, &si.MinRealm)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &si, nil
}

func (r *ShopRepository) DecrementShopStock(ctx context.Context, shopID string, itemName string, quantity int) error {
	query := `UPDATE shop_inventory SET quantity = quantity - $3 WHERE shop_id = $1 AND item_name = $2 AND quantity >= $3`
	_, err := r.db.Pool().Exec(ctx, query, shopID, itemName, quantity)
	return err
}

// ── Auctions ──

func (r *ShopRepository) CreateAuction(ctx context.Context, a *Auction) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	query := `INSERT INTO auctions (id, seller_id, item_id, item_name, quantity, price, deposit, status, created_at, expires_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Pool().Exec(ctx, query,
		a.ID, a.SellerID, a.ItemID, a.ItemName, a.Quantity, a.Price, a.Deposit, a.Status, a.CreatedAt, a.ExpiresAt)
	return err
}

func (r *ShopRepository) GetAuctionByID(ctx context.Context, auctionID string) (*Auction, error) {
	query := `SELECT id, seller_id, item_id, item_name, quantity, price, deposit, status, created_at, expires_at, buyer_id, sold_at FROM auctions WHERE id = $1`
	row := r.db.Pool().QueryRow(ctx, query, auctionID)
	var a Auction
	var expiresAt, soldAt *time.Time
	var buyerID *string
	err := row.Scan(&a.ID, &a.SellerID, &a.ItemID, &a.ItemName, &a.Quantity, &a.Price, &a.Deposit, &a.Status, &a.CreatedAt, &expiresAt, &buyerID, &soldAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.ExpiresAt = expiresAt
	if buyerID != nil {
		a.BuyerID = buyerID
	}
	if soldAt != nil {
		a.SoldAt = soldAt
	}
	return &a, nil
}

func (r *ShopRepository) ListActiveAuctions(ctx context.Context) ([]*Auction, error) {
	query := `SELECT id, seller_id, item_id, item_name, quantity, price, deposit, status, created_at, expires_at, buyer_id, sold_at FROM auctions WHERE status = 'active' AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC`
	rows, err := r.db.Pool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []*Auction
	for rows.Next() {
		var a Auction
		var expiresAt, soldAt *time.Time
		var buyerID *string
		if err := rows.Scan(&a.ID, &a.SellerID, &a.ItemID, &a.ItemName, &a.Quantity, &a.Price, &a.Deposit, &a.Status, &a.CreatedAt, &expiresAt, &buyerID, &soldAt); err != nil {
			return nil, err
		}
		a.ExpiresAt = expiresAt
		if buyerID != nil {
			a.BuyerID = buyerID
		}
		if soldAt != nil {
			a.SoldAt = soldAt
		}
		auctions = append(auctions, &a)
	}
	return auctions, rows.Err()
}

func (r *ShopRepository) BuyAuction(ctx context.Context, auctionID string, buyerID string) error {
	query := `UPDATE auctions SET status = 'sold', buyer_id = $2, sold_at = NOW() WHERE id = $1 AND status = 'active'`
	tag, err := r.db.Pool().Exec(ctx, query, auctionID, buyerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *ShopRepository) CancelAuction(ctx context.Context, auctionID string, sellerID string) error {
	query := `UPDATE auctions SET status = 'cancelled' WHERE id = $1 AND seller_id = $2 AND status = 'active'`
	tag, err := r.db.Pool().Exec(ctx, query, auctionID, sellerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
