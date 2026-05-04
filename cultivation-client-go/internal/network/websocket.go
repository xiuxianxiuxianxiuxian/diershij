package network

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	mu        sync.RWMutex
	conn      *websocket.Conn
	url       string
	connected bool
	reconnect bool
	handlers  map[string][]func([]byte)
	sendChan  chan []byte
	done      chan struct{}
	reqID     atomic.Int64 // 自增请求ID
}

var wsInstance *WebSocketClient
var wsOnce sync.Once

func GetWebSocketClient() *WebSocketClient {
	wsOnce.Do(func() {
		wsInstance = &WebSocketClient{
			url:       "ws://localhost:8081/ws",
			connected: false,
			reconnect: true,
			handlers:  make(map[string][]func([]byte)),
			sendChan:  make(chan []byte, 100),
			done:      make(chan struct{}),
		}
	})
	return wsInstance
}

func (c *WebSocketClient) Connect() error {
	c.mu.Lock()
	if c.connected {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	token := store.GetAuthStore().GetToken()
	if token == "" {
		return fmt.Errorf("not authenticated")
	}

	wsURL := fmt.Sprintf("%s?token=%s", c.url, token)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.connected = true
	c.mu.Unlock()

	c.handleMessage(&types.WSMessage{
		Type:    "connected",
		Payload: map[string]interface{}{},
	})

	go c.readLoop()
	go c.writeLoop()

	return nil
}

func (c *WebSocketClient) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.reconnect = false
	if c.connected {
		close(c.done)
		if c.conn != nil {
			c.conn.Close()
		}
		c.connected = false
	}
}

// nextRequestID 生成唯一的请求ID
func (c *WebSocketClient) nextRequestID() string {
	return fmt.Sprintf("cli_%d_%d", time.Now().UnixMilli(), c.reqID.Add(1))
}

// Send 发送 WebSocket 消息
func (c *WebSocketClient) Send(messageType string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := types.WSMessage{
		Type:    types.WSMessageType(messageType),
		Payload: make(map[string]interface{}),
	}

	if err := json.Unmarshal(payload, &msg.Payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case c.sendChan <- msgBytes:
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout")
	}

	return nil
}

func (c *WebSocketClient) RegisterHandler(messageType string, handler func([]byte)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[messageType] = append(c.handlers[messageType], handler)
}

func (c *WebSocketClient) handleMessage(msg *types.WSMessage) {
	c.mu.RLock()
	handlers := c.handlers[string(msg.Type)]
	c.mu.RUnlock()

	if len(handlers) > 0 {
		payload, _ := json.Marshal(msg.Payload)
		for _, handler := range handlers {
			handler(payload)
		}
	}
}

func (c *WebSocketClient) readLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			var msg types.WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Printf("WebSocket read error: %v\n", err)

				c.handleMessage(&types.WSMessage{
					Type:    "disconnect",
					Payload: map[string]interface{}{},
				})

				c.reconnectLoop()
				return
			}

			c.handleMessage(&msg)
		}
	}
}

func (c *WebSocketClient) writeLoop() {
	for {
		select {
		case <-c.done:
			return
		case msg := <-c.sendChan:
			c.mu.RLock()
			conn := c.conn
			connected := c.connected
			c.mu.RUnlock()

			if connected && conn != nil {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					fmt.Printf("WebSocket write error: %v\n", err)
				}
			}
		}
	}
}

func (c *WebSocketClient) reconnectLoop() {
	for {
		c.mu.RLock()
		reconnect := c.reconnect
		c.mu.RUnlock()

		if !reconnect {
			return
		}

		time.Sleep(5 * time.Second)
		if err := c.Connect(); err != nil {
			fmt.Printf("Reconnection failed: %v\n", err)
		} else {
			return
		}
	}
}

// SendOperation 发送操作消息（自动附带 request_id）
func (c *WebSocketClient) SendOperation(actionType string, params map[string]interface{}) error {
	if params == nil {
		params = make(map[string]interface{})
	}
	data := map[string]interface{}{
		"action_type": actionType,
		"params":      params,
		"request_id":  c.nextRequestID(),
	}
	return c.Send("operation", data)
}

// Cultivate 修炼
func (c *WebSocketClient) Cultivate() error {
	return c.SendOperation("cultivate", nil)
}

// Meditate 打坐
func (c *WebSocketClient) Meditate() error {
	return c.SendOperation("meditate", nil)
}

// Sleep 休息
func (c *WebSocketClient) Sleep() error {
	return c.SendOperation("sleep", nil)
}

// Breakthrough 突破
func (c *WebSocketClient) Breakthrough() error {
	return c.SendOperation("breakthrough", nil)
}

// Move 移动
func (c *WebSocketClient) Move(regionID string, x, y float64) error {
	params := map[string]interface{}{
		"region_id": regionID,
		"x":         x,
		"y":         y,
	}
	return c.SendOperation("move", params)
}

// Combat 战斗
func (c *WebSocketClient) Combat(targetID string) error {
	params := map[string]interface{}{
		"target_id": targetID,
	}
	return c.SendOperation("combat", params)
}

// Gather 采集
func (c *WebSocketClient) Gather(resourceType string, quantity int) error {
	params := map[string]interface{}{
		"resource_type": resourceType,
		"quantity":      quantity,
	}
	return c.SendOperation("gather", params)
}

// Explore 探索
func (c *WebSocketClient) Explore() error {
	return c.SendOperation("explore", nil)
}

// SendMessage 发消息
func (c *WebSocketClient) SendMessage(content string, msgType string, receiverID string) error {
	params := map[string]interface{}{
		"content":      content,
		"message_type": msgType,
		"receiver_id":  receiverID,
	}
	return c.SendOperation("send_message", params)
}

// CastSpell 施法
func (c *WebSocketClient) CastSpell(spellID string, targetID string) error {
	params := map[string]interface{}{
		"spell_id":  spellID,
		"target_id": targetID,
	}
	return c.SendOperation("cast_spell", params)
}

// AddFriend 加好友
func (c *WebSocketClient) AddFriend(name string) error {
	params := map[string]interface{}{
		"name": name,
	}
	return c.SendOperation("add_friend", params)
}

// RemoveFriend 删除好友
func (c *WebSocketClient) RemoveFriend(friendID string) error {
	params := map[string]interface{}{
		"friend_id": friendID,
	}
	return c.SendOperation("remove_friend", params)
}

// AcceptFriendRequest 接受好友请求
func (c *WebSocketClient) AcceptFriendRequest(requestID string) error {
	params := map[string]interface{}{
		"request_id": requestID,
	}
	return c.SendOperation("accept_friend", params)
}

// CreateSect 创建宗门
func (c *WebSocketClient) CreateSect(sectName string) error {
	params := map[string]interface{}{
		"sect_name": sectName,
	}
	return c.SendOperation("form_sect", params)
}

// JoinSect 加入宗门
func (c *WebSocketClient) JoinSect(sectID string) error {
	params := map[string]interface{}{
		"sect_id": sectID,
	}
	return c.SendOperation("join_sect", params)
}

// LeaveSect 离开宗门
func (c *WebSocketClient) LeaveSect() error {
	return c.SendOperation("leave_sect", nil)
}

// Flee 逃跑
func (c *WebSocketClient) Flee() error {
	return c.SendOperation("flee", nil)
}

// UseSkill 使用技能
func (c *WebSocketClient) UseSkill() error {
	return c.SendOperation("use_skill", nil)
}

// Craft 炼制（炼器/炼丹）
func (c *WebSocketClient) Craft(recipeID string) error {
	params := map[string]interface{}{
		"recipe_id": recipeID,
	}
	return c.SendOperation("craft", params)
}

// Trade 交易
func (c *WebSocketClient) Trade(targetID string, itemID string, price float64) error {
	params := map[string]interface{}{
		"target_id": targetID,
		"item_id":   itemID,
		"price":     price,
	}
	return c.SendOperation("trade", params)
}

// CreateMethod 自创功法
func (c *WebSocketClient) CreateMethod() error {
	return c.SendOperation("create_method", nil)
}
