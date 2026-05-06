package types

// Mail represents a mail message in the mailbox system.
type Mail struct {
	ID                    string `json:"id"`
	SenderID              string `json:"sender_id"`
	ReceiverID            string `json:"receiver_id"`
	SenderName            string `json:"sender_name"`
	Title                 string `json:"title"`
	Content               string `json:"content"`
	MailType              string `json:"mail_type"` // system, player, reward
	IsRead                bool   `json:"is_read"`
	HasAttachment         bool   `json:"has_attachment"`
	AttachmentItemID      string `json:"attachment_item_id"`
	AttachmentItemName    string `json:"attachment_item_name"`
	AttachmentQuantity    int    `json:"attachment_quantity"`
	AttachmentSpiritStones int64 `json:"attachment_spirit_stones"`
	IsClaimed             bool   `json:"is_claimed"`
	CreatedAt             int64  `json:"created_at"`
}
