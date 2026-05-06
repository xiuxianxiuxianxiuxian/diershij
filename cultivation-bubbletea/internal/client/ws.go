package client

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WsMessage represents a parsed WebSocket message from the gateway.
type WsMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// GameState holds cached entity data from state_sync / entity_update messages.
type GameState struct {
	mu       sync.RWMutex
	Entity   map[string]interface{} `json:"entity"`
	Spells   []interface{}          `json:"spells"`
	Items    []interface{}          `json:"items"`
	Friends  []interface{}          `json:"friends"`
	Sect     map[string]interface{} `json:"sect"`
}

// Entity returns a copy of the current entity data.
func (gs *GameState) EntityCopy() map[string]interface{} {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return copyMap(gs.Entity)
}

// Spells returns a copy of the spells list.
func (gs *GameState) SpellsCopy() []interface{} {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	out := make([]interface{}, len(gs.Spells))
	copy(out, gs.Spells)
	return out
}

// ItemsCopy returns a copy of the items list.
func (gs *GameState) ItemsCopy() []interface{} {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	out := make([]interface{}, len(gs.Items))
	copy(out, gs.Items)
	return out
}

// GetStatus returns the entity's status string.
func (gs *GameState) GetStatus() string {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	if gs.Entity == nil {
		return ""
	}
	if s, ok := gs.Entity["status"].(string); ok {
		return s
	}
	return ""
}

// UpdateFromStateSync updates the game state from a state_sync or entity_update payload.
func (gs *GameState) UpdateFromStateSync(payload map[string]interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if ent, ok := payload["entity"].(map[string]interface{}); ok {
		gs.Entity = ent
	}
	if s, ok := payload["spells"].([]interface{}); ok {
		gs.Spells = s
	}
	if it, ok := payload["items"].([]interface{}); ok {
		gs.Items = it
	}
}

// UpdateFriends caches friends data.
func (gs *GameState) UpdateFriends(payload map[string]interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if f, ok := payload["friends"].([]interface{}); ok {
		gs.Friends = f
	}
}

// UpdateSect caches sect data.
func (gs *GameState) UpdateSect(payload map[string]interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if s, ok := payload["sect_id"].(string); ok && s != "" {
		gs.Sect = payload
	}
}

// ConnectWebSocket dials the gateway WebSocket, starts a background reader
// goroutine, and pipes parsed messages to msgCh. doneCh is closed when the
// connection drops. GameState is updated in real-time for synchronous access.
func ConnectWebSocket(token string, msgCh chan<- WsMessage, state *GameState) (*websocket.Conn, error) {
	u := fmt.Sprintf("ws://localhost:8081/ws?token=%s", token)
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	go func() {
		defer conn.Close()

		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				msgCh <- WsMessage{Type: "system", Payload: map[string]interface{}{
					"message": fmt.Sprintf("[断开] %v", err),
				}}
				return
			}

			var msg struct {
				Type    string                 `json:"type"`
				Payload map[string]interface{} `json:"payload"`
			}
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}

			conn.SetReadDeadline(time.Now().Add(120 * time.Second))

			// Update game state cache for state_sync and entity_update
			switch msg.Type {
			case "state_sync", "entity_update":
				if state != nil {
					state.UpdateFromStateSync(msg.Payload)
				}
			}

			msgCh <- WsMessage{Type: msg.Type, Payload: msg.Payload}
		}
	}()

	return conn, nil
}

// SendAction sends a game operation over WebSocket.
func SendAction(conn *websocket.Conn, actionType string, params map[string]interface{}) error {
	msg := map[string]interface{}{
		"type": "operation",
		"payload": map[string]interface{}{
			"action_type": actionType,
			"params":      params,
		},
	}
	return conn.WriteJSON(msg)
}

// SendChat sends a chat message over WebSocket.
func SendChat(conn *websocket.Conn, content, channel string) error {
	msg := map[string]interface{}{
		"type": "chat",
		"payload": map[string]interface{}{
			"content": content,
			"channel": channel,
		},
	}
	return conn.WriteJSON(msg)
}

func copyMap(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}
	out := make(map[string]interface{}, len(original))
	for k, v := range original {
		out[k] = v
	}
	return out
}
