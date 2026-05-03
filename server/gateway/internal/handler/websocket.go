package handler

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/cultivation-world/gateway/internal/service"
    "github.com/cultivation-world/shared/types"
    "github.com/gorilla/websocket"
)

type WebSocketHub struct {
    clients    map[*WebSocketClient]bool
    entityMap  map[types.EntityID]*WebSocketClient
    broadcast  chan *BroadcastMessage
    register   chan *WebSocketClient
    unregister chan *WebSocketClient
    mu         sync.RWMutex
}

type BroadcastMessage struct {
    EntityID types.EntityID
    Message  types.Message
}

func NewWebSocketHub() *WebSocketHub {
    return &WebSocketHub{
        clients:    make(map[*WebSocketClient]bool),
        entityMap:  make(map[types.EntityID]*WebSocketClient),
        broadcast:  make(chan *BroadcastMessage, 256),
        register:   make(chan *WebSocketClient),
        unregister: make(chan *WebSocketClient),
    }
}

func (h *WebSocketHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.entityMap[client.entityID] = client
            h.mu.Unlock()
            log.Printf("Client connected: %s", client.entityID)

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                delete(h.entityMap, client.entityID)
                close(client.send)
            }
            h.mu.Unlock()
            log.Printf("Client disconnected: %s", client.entityID)

        case msg := <-h.broadcast:
            h.mu.RLock()
            if client, ok := h.entityMap[msg.EntityID]; ok {
                select {
                case client.send <- msg.Message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    delete(h.entityMap, client.entityID)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *WebSocketHub) Register(client *WebSocketClient) {
    h.register <- client
}

func (h *WebSocketHub) Unregister(client *WebSocketClient) {
    h.unregister <- client
}

func (h *WebSocketHub) BroadcastToEntity(entityID types.EntityID, msg types.Message) {
    h.broadcast <- &BroadcastMessage{EntityID: entityID, Message: msg}
}

func (h *WebSocketHub) BroadcastToAll(msg types.Message) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    for _, client := range h.entityMap {
        select {
        case client.send <- msg:
        default:
        }
    }
}

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

type WebSocketClient struct {
    entityID   types.EntityID
    conn       *websocket.Conn
    send       chan types.Message
    hub        *WebSocketHub
    gameClient *service.GameServiceClient
}

func NewWebSocketClient(entityID types.EntityID, conn *websocket.Conn, hub *WebSocketHub, gameClient *service.GameServiceClient) *WebSocketClient {
    return &WebSocketClient{
        entityID:   entityID,
        conn:       conn,
        send:       make(chan types.Message, 256),
        hub:        hub,
        gameClient: gameClient,
    }
}

func (c *WebSocketClient) ReadPump() {
    defer func() {
        c.hub.Unregister(c)
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        var msg types.Message
        if err := json.Unmarshal(message, &msg); err != nil {
            continue
        }

        c.handleMessage(&msg)
    }
}

func (c *WebSocketClient) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.conn.WriteJSON(message); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *WebSocketClient) handleMessage(msg *types.Message) {
    switch msg.Type {
    case types.MessageTypeOperation:
        c.handleOperation(msg)
    case types.MessageTypeChat:
        c.handleChat(msg)
    default:
        c.sendError("unknown message type")
    }
}

func (c *WebSocketClient) handleOperation(msg *types.Message) {
    actionType, ok := msg.Payload["action_type"].(string)
    if !ok {
        c.sendError("missing action_type")
        return
    }

    params := make(map[string]interface{})
    if p, ok := msg.Payload["params"].(map[string]interface{}); ok {
        params = p
    }

    op := types.NewOperation(c.entityID, types.ActionType(actionType), params)

    result, err := c.gameClient.ExecuteOperation(op)
    if err != nil {
        c.sendError(err.Error())
        return
    }

    payload := map[string]interface{}{
        "success":   result.Success,
        "message":   result.Message,
        "effects":   result.Effects,
        "timestamp": result.Timestamp,
    }
    c.send <- types.Message{
        Type:    types.MessageTypeOpResult,
        Payload: payload,
    }
}

func (c *WebSocketClient) handleChat(msg *types.Message) {
    content, ok := msg.Payload["content"].(string)
    if !ok {
        return
    }

    channel, _ := msg.Payload["channel"].(string)

    chatMsg := types.Message{
        Type: types.MessageTypeChat,
        Payload: map[string]interface{}{
            "sender_id":   string(c.entityID),
            "content":     content,
            "channel":     channel,
            "timestamp":   time.Now().UnixNano(),
        },
    }

    c.hub.BroadcastToAll(chatMsg)
}

func (c *WebSocketClient) sendError(message string) {
    c.send <- types.Message{
        Type: types.MessageTypeError,
        Payload: map[string]interface{}{
            "message": message,
        },
    }
}
