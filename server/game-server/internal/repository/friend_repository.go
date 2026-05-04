package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Friendship struct {
	EntityID  string    `json:"entity_id"`
	FriendID  string    `json:"friend_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendRequest struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendRepository struct {
	db *Database
}

func NewFriendRepository(db *Database) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) AddFriend(ctx context.Context, entityID, friendID string) error {
	query := `INSERT INTO friendships (entity_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Pool().Exec(ctx, query, entityID, friendID)
	return err
}

func (r *FriendRepository) RemoveFriend(ctx context.Context, entityID, friendID string) error {
	query := `DELETE FROM friendships WHERE entity_id = $1 AND friend_id = $2`
	_, err := r.db.Pool().Exec(ctx, query, entityID, friendID)
	return err
}

func (r *FriendRepository) GetFriends(ctx context.Context, entityID string) ([]*Friendship, error) {
	query := `SELECT entity_id, friend_id, created_at FROM friendships WHERE entity_id = $1`
	rows, err := r.db.Pool().Query(ctx, query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*Friendship
	for rows.Next() {
		var f Friendship
		if err := rows.Scan(&f.EntityID, &f.FriendID, &f.CreatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, &f)
	}
	return friends, nil
}

func (r *FriendRepository) AreFriends(ctx context.Context, entityID, friendID string) (bool, error) {
	query := `SELECT 1 FROM friendships WHERE entity_id = $1 AND friend_id = $2`
	var exists int
	err := r.db.Pool().QueryRow(ctx, query, entityID, friendID).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func (r *FriendRepository) CreateRequest(ctx context.Context, fromID, toID string) (string, error) {
	query := `INSERT INTO friend_requests (from_id, to_id, status) VALUES ($1, $2, 'pending') RETURNING id`
	var id string
	err := r.db.Pool().QueryRow(ctx, query, fromID, toID).Scan(&id)
	return id, err
}

func (r *FriendRepository) GetPendingRequest(ctx context.Context, fromID, toID string) (*FriendRequest, error) {
	query := `SELECT id, from_id, to_id, status, created_at FROM friend_requests WHERE from_id = $1 AND to_id = $2 AND status = 'pending'`
	var fr FriendRequest
	err := r.db.Pool().QueryRow(ctx, query, fromID, toID).Scan(&fr.ID, &fr.FromID, &fr.ToID, &fr.Status, &fr.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &fr, err
}

func (r *FriendRepository) GetRequestByID(ctx context.Context, requestID string) (*FriendRequest, error) {
	query := `SELECT id, from_id, to_id, status, created_at FROM friend_requests WHERE id = $1`
	var fr FriendRequest
	err := r.db.Pool().QueryRow(ctx, query, requestID).Scan(&fr.ID, &fr.FromID, &fr.ToID, &fr.Status, &fr.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &fr, err
}

func (r *FriendRepository) AcceptRequest(ctx context.Context, requestID string) error {
	query := `UPDATE friend_requests SET status = 'accepted', updated_at = NOW() WHERE id = $1 AND status = 'pending'`
	_, err := r.db.Pool().Exec(ctx, query, requestID)
	return err
}
