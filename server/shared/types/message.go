package types

type MessageType string

const (
    MessageTypeAuth         MessageType = "auth"
    MessageTypeAuthResult   MessageType = "auth_result"
    MessageTypeOperation    MessageType = "operation"
    MessageTypeOpResult     MessageType = "op_result"
    MessageTypeStateSync    MessageType = "state_sync"
    MessageTypeEntityUpdate MessageType = "entity_update"
    MessageTypeWorldEvent   MessageType = "world_event"
    MessageTypeChat         MessageType = "chat"
    MessageTypeSystem       MessageType = "system"
    MessageTypeError        MessageType = "error"
)

type Message struct {
    Type      MessageType             `json:"type"`
    Payload   map[string]interface{}  `json:"payload"`
    Timestamp int64                   `json:"timestamp"`
    RequestID string                  `json:"request_id,omitempty"`
}

type AuthPayload struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Token    string `json:"token,omitempty"`
}

type AuthResultPayload struct {
    Success bool     `json:"success"`
    Token   string   `json:"token,omitempty"`
    Entity  *Entity  `json:"entity,omitempty"`
    Message string   `json:"message"`
}

type StateSyncPayload struct {
    Entity       *Entity      `json:"entity"`
    Region       *Region      `json:"region"`
    NearbyEntities []EntityID `json:"nearby_entities"`
    WorldTime    int64        `json:"world_time"`
}

type ChatPayload struct {
    SenderID   EntityID `json:"sender_id"`
    SenderName string   `json:"sender_name"`
    Channel    string   `json:"channel"`
    Content    string   `json:"content"`
}

type ErrorPayload struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// DBMessage represents a message stored in the database (chat messages)
type DBMessage struct {
    ID         string    `json:"id"`
    SenderID   string    `json:"sender_id"`
    ReceiverID string    `json:"receiver_id"`
    Type       string    `json:"type"`
    Content    string    `json:"content"`
    IsRead     bool      `json:"is_read"`
    CreatedAt  int64     `json:"created_at"`
}
