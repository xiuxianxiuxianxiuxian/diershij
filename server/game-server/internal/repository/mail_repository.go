package repository

import (
	"context"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/google/uuid"
)

// MailRepository defines the interface for mail data access.
type MailRepository interface {
	Create(ctx context.Context, mail *types.Mail) error
	GetByID(ctx context.Context, id string) (*types.Mail, error)
	GetByReceiver(ctx context.Context, receiverID types.EntityID, limit int) ([]*types.Mail, error)
	GetUnreadByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error)
	GetUnclaimedByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error)
	MarkAsRead(ctx context.Context, mailID string) error
	MarkAsClaimed(ctx context.Context, mailID string) error
	Delete(ctx context.Context, mailID string) error
}

type PostgresMailRepository struct {
	db *Database
}

func NewPostgresMailRepository(db *Database) *PostgresMailRepository {
	return &PostgresMailRepository{db: db}
}

func (r *PostgresMailRepository) Create(ctx context.Context, mail *types.Mail) error {
	if mail.ID == "" {
		mail.ID = uuid.New().String()
	}
	if mail.CreatedAt == 0 {
		mail.CreatedAt = time.Now().Unix()
	}

	query := `
		INSERT INTO mails (id, sender_id, receiver_id, sender_name, title, content,
		                   mail_type, is_read, has_attachment,
		                   attachment_item_id, attachment_item_name, attachment_quantity,
		                   attachment_spirit_stones, is_claimed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	senderID := interface{}(mail.SenderID)
	if mail.SenderID == "" {
		senderID = nil
	}
	attachmentItemID := interface{}(mail.AttachmentItemID)
	if mail.AttachmentItemID == "" {
		attachmentItemID = nil
	}

	_, err := r.db.Pool().Exec(ctx, query,
		mail.ID, senderID, mail.ReceiverID, mail.SenderName, mail.Title, mail.Content,
		mail.MailType, mail.IsRead, mail.HasAttachment,
		attachmentItemID, mail.AttachmentItemName, mail.AttachmentQuantity,
		mail.AttachmentSpiritStones, mail.IsClaimed, time.Unix(mail.CreatedAt, 0),
	)
	return err
}

func (r *PostgresMailRepository) GetByID(ctx context.Context, id string) (*types.Mail, error) {
	query := `
		SELECT id, sender_id, receiver_id, sender_name, title, content,
		       mail_type, is_read, has_attachment,
		       attachment_item_id, attachment_item_name, attachment_quantity,
		       attachment_spirit_stones, is_claimed, created_at
		FROM mails WHERE id = $1
	`
	row := r.db.Pool().QueryRow(ctx, query, id)
	return r.scanMail(row)
}

func (r *PostgresMailRepository) GetByReceiver(ctx context.Context, receiverID types.EntityID, limit int) ([]*types.Mail, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, sender_id, receiver_id, sender_name, title, content,
		       mail_type, is_read, has_attachment,
		       attachment_item_id, attachment_item_name, attachment_quantity,
		       attachment_spirit_stones, is_claimed, created_at
		FROM mails
		WHERE receiver_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Pool().Query(ctx, query, receiverID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []*types.Mail
	for rows.Next() {
		m, err := r.scanMail(rows)
		if err != nil {
			return nil, err
		}
		mails = append(mails, m)
	}
	return mails, rows.Err()
}

func (r *PostgresMailRepository) GetUnreadByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error) {
	query := `
		SELECT id, sender_id, receiver_id, sender_name, title, content,
		       mail_type, is_read, has_attachment,
		       attachment_item_id, attachment_item_name, attachment_quantity,
		       attachment_spirit_stones, is_claimed, created_at
		FROM mails
		WHERE receiver_id = $1 AND is_read = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool().Query(ctx, query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []*types.Mail
	for rows.Next() {
		m, err := r.scanMail(rows)
		if err != nil {
			return nil, err
		}
		mails = append(mails, m)
	}
	return mails, rows.Err()
}

func (r *PostgresMailRepository) GetUnclaimedByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error) {
	query := `
		SELECT id, sender_id, receiver_id, sender_name, title, content,
		       mail_type, is_read, has_attachment,
		       attachment_item_id, attachment_item_name, attachment_quantity,
		       attachment_spirit_stones, is_claimed, created_at
		FROM mails
		WHERE receiver_id = $1 AND has_attachment = true AND is_claimed = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool().Query(ctx, query, receiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mails []*types.Mail
	for rows.Next() {
		m, err := r.scanMail(rows)
		if err != nil {
			return nil, err
		}
		mails = append(mails, m)
	}
	return mails, rows.Err()
}

func (r *PostgresMailRepository) MarkAsRead(ctx context.Context, mailID string) error {
	_, err := r.db.Pool().Exec(ctx, "UPDATE mails SET is_read = true WHERE id = $1", mailID)
	return err
}

func (r *PostgresMailRepository) MarkAsClaimed(ctx context.Context, mailID string) error {
	_, err := r.db.Pool().Exec(ctx, "UPDATE mails SET is_claimed = true WHERE id = $1", mailID)
	return err
}

func (r *PostgresMailRepository) Delete(ctx context.Context, mailID string) error {
	_, err := r.db.Pool().Exec(ctx, "DELETE FROM mails WHERE id = $1", mailID)
	return err
}

func (r *PostgresMailRepository) scanMail(row interface{ Scan(dest ...interface{}) error }) (*types.Mail, error) {
	m := &types.Mail{}
	var createdAt time.Time
	var senderID, attachmentItemID *string

	err := row.Scan(
		&m.ID, &senderID, &m.ReceiverID, &m.SenderName, &m.Title, &m.Content,
		&m.MailType, &m.IsRead, &m.HasAttachment,
		&attachmentItemID, &m.AttachmentItemName, &m.AttachmentQuantity,
		&m.AttachmentSpiritStones, &m.IsClaimed, &createdAt,
	)
	if err != nil {
		return nil, err
	}

	if senderID != nil {
		m.SenderID = *senderID
	}
	if attachmentItemID != nil {
		m.AttachmentItemID = *attachmentItemID
	}
	m.CreatedAt = createdAt.Unix()

	return m, nil
}
