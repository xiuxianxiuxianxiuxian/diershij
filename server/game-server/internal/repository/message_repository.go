package repository

import (
	"context"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/google/uuid"
)

type MessageRepository interface {
	Create(ctx context.Context, message *types.DBMessage) error
	GetByID(ctx context.Context, id string) (*types.DBMessage, error)
	GetByReceiver(ctx context.Context, receiverID types.EntityID, limit int) ([]*types.DBMessage, error)
	GetUnreadByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.DBMessage, error)
	MarkAsRead(ctx context.Context, messageID string) error
	MarkAllAsRead(ctx context.Context, receiverID types.EntityID) error
	Delete(ctx context.Context, messageID string) error
	GetByType(ctx context.Context, receiverID types.EntityID, msgType string, limit int) ([]*types.DBMessage, error)
}

type PostgresMessageRepository struct {
	db *Database
}

func NewPostgresMessageRepository(db *Database) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, message *types.DBMessage) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	if message.CreatedAt == 0 {
		message.CreatedAt = time.Now().Unix()
	}

	query := `
		INSERT INTO messages (id, sender_id, receiver_id, type, content, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// Use nil for receiver_id when empty (world/broadcast messages) to avoid UUID cast error
	receiverID := interface{}(message.ReceiverID)
	if message.ReceiverID == "" {
		receiverID = nil
	}

	_, err := r.db.Pool().Exec(ctx, query,
		message.ID, message.SenderID, receiverID, message.Type, message.Content, message.IsRead, time.Unix(message.CreatedAt, 0),
	)
	return err
}

func (r *PostgresMessageRepository) GetByID(ctx context.Context, id string) (*types.DBMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, type, content, is_read, created_at
		FROM messages WHERE id = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, id)
	return r.scanMessage(row)
}

func (r *PostgresMessageRepository) GetByReceiver(ctx context.Context, receiverID types.EntityID, limit int) ([]*types.DBMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, sender_id, receiver_id, type, content, is_read, created_at
		FROM messages
		WHERE receiver_id = $1 OR receiver_id IS NULL
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Pool().Query(ctx, query, receiverID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*types.DBMessage
	for rows.Next() {
		msg, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *PostgresMessageRepository) GetUnreadByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.DBMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, type, content, is_read, created_at
		FROM messages
		WHERE (receiver_id = $1 OR receiver_id IS NULL) AND is_read = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool().Query(ctx, query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*types.DBMessage
	for rows.Next() {
		msg, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *PostgresMessageRepository) MarkAsRead(ctx context.Context, messageID string) error {
	query := `UPDATE messages SET is_read = true WHERE id = $1`
	_, err := r.db.Pool().Exec(ctx, query, messageID)
	return err
}

func (r *PostgresMessageRepository) MarkAllAsRead(ctx context.Context, receiverID types.EntityID) error {
	query := `UPDATE messages SET is_read = true WHERE receiver_id = $1 AND is_read = false`
	_, err := r.db.Pool().Exec(ctx, query, receiverID)
	return err
}

func (r *PostgresMessageRepository) Delete(ctx context.Context, messageID string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Pool().Exec(ctx, query, messageID)
	return err
}

func (r *PostgresMessageRepository) GetByType(ctx context.Context, receiverID types.EntityID, msgType string, limit int) ([]*types.DBMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, sender_id, receiver_id, type, content, is_read, created_at
		FROM messages
		WHERE type = $1 AND (receiver_id = $2 OR receiver_id IS NULL)
		ORDER BY created_at DESC
		LIMIT $3
	`
	rows, err := r.db.Pool().Query(ctx, query, msgType, receiverID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*types.DBMessage
	for rows.Next() {
		msg, err := r.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *PostgresMessageRepository) scanMessage(scanner interface {
	Scan(dest ...interface{}) error
}) (*types.DBMessage, error) {
	msg := &types.DBMessage{}
	var createdAt time.Time
	var senderID, receiverID *string

	err := scanner.Scan(
		&msg.ID, &senderID, &receiverID, &msg.Type, &msg.Content, &msg.IsRead, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	if senderID != nil {
		msg.SenderID = *senderID
	}
	if receiverID != nil {
		msg.ReceiverID = *receiverID
	}
	msg.CreatedAt = createdAt.Unix()

	return msg, nil
}
