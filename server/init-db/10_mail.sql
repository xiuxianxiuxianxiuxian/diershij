-- Migration 10: Mail/Mailbox system
-- Supports system mail, player mail, attachments (items + spirit stones)

CREATE TABLE IF NOT EXISTS mails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES entities(id) ON DELETE SET NULL,
    receiver_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    sender_name VARCHAR(100) DEFAULT '',
    title VARCHAR(200) NOT NULL,
    content TEXT DEFAULT '',
    mail_type VARCHAR(20) NOT NULL DEFAULT 'system',   -- system, player, reward
    is_read BOOLEAN DEFAULT false,
    has_attachment BOOLEAN DEFAULT false,
    attachment_item_id UUID REFERENCES items(id) ON DELETE SET NULL,
    attachment_item_name VARCHAR(100) DEFAULT '',
    attachment_quantity INTEGER DEFAULT 0,
    attachment_spirit_stones BIGINT DEFAULT 0,
    is_claimed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '30 days')
);

CREATE INDEX IF NOT EXISTS idx_mails_receiver ON mails(receiver_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_mails_unread ON mails(receiver_id, is_read) WHERE is_read = false;
CREATE INDEX IF NOT EXISTS idx_mails_unclaimed ON mails(receiver_id, is_claimed, has_attachment) WHERE has_attachment = true AND is_claimed = false;
