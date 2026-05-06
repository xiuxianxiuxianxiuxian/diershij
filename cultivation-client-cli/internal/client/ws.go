package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp int64                  `json:"timestamp,omitempty"`
}

var currentEntity map[string]interface{}
var currentSpells []interface{}
var currentItems []interface{}
var currentFriends []interface{}
var currentSect map[string]interface{}

// ConnectWebSocket dials the gateway and returns the connection.
// It starts a background goroutine that reads messages and sends
// formatted display strings to msgCh. doneCh is closed when the
// connection is lost.
func ConnectWebSocket(token string, msgCh chan<- string, doneCh chan<- struct{}) (*websocket.Conn, error) {
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
		defer func() {
			close(doneCh)
		}()

		for {
			_, raw, err := conn.ReadMessage()
			if err != nil {
				msgCh <- fmt.Sprintf("[断开] %v", err)
				return
			}

			var msg wsMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}

			conn.SetReadDeadline(time.Now().Add(120 * time.Second))

			switch msg.Type {
			case "state_sync":
				if ent, ok := msg.Payload["entity"].(map[string]interface{}); ok {
					currentEntity = ent
				}
				if s, ok := msg.Payload["spells"].([]interface{}); ok {
					currentSpells = s
				}
				if it, ok := msg.Payload["items"].([]interface{}); ok {
					currentItems = it
				}
				msgCh <- formatStateSync(msg.Payload)

			case "entity_update":
				if ent, ok := msg.Payload["entity"].(map[string]interface{}); ok {
					currentEntity = ent
				}
				if s, ok := msg.Payload["spells"].([]interface{}); ok {
					currentSpells = s
				}
				if it, ok := msg.Payload["items"].([]interface{}); ok {
					currentItems = it
				}

			case "op_result":
				msgCh <- formatOpResult(msg.Payload)

			case "error":
				msgCh <- fmt.Sprintf("[错误] %s", getStr(msg.Payload, "message"))

			case "chat":
				sender := getStr(msg.Payload, "sender_name")
				if sender == "" {
					sender = getStr(msg.Payload, "sender_id")
				}
				content := getStr(msg.Payload, "content")
				if content != "" {
					msgCh <- fmt.Sprintf("[聊天] %s: %s", sender, content)
				}

			case "world_event":
				msgCh <- fmt.Sprintf("[世界事件] %s", getStr(msg.Payload, "description"))

			case "system":
				msgCh <- fmt.Sprintf("[系统] %s", getStr(msg.Payload, "message"))

			case "announcement":
				msgCh <- fmt.Sprintf("[公告] %s", getStr(msg.Payload, "content"))

			case "new_message":
				msgCh <- fmt.Sprintf("[消息] 来自 %s: %s",
					getStr(msg.Payload, "sender_name"),
					getStr(msg.Payload, "content"))

			case "friend_request":
				msgCh <- fmt.Sprintf("[好友] %s 请求添加你为好友",
					getStr(msg.Payload, "from_name"))
			}
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

// GetCharacter returns the latest cached entity data from state_sync/entity_update.
func GetCharacter() map[string]interface{} {
	return currentEntity
}

// GetSpells returns the spells list from the cached state_sync data.
func GetSpells() []interface{} {
	return currentSpells
}

// GetItems returns the inventory items list from the cached state_sync data.
func GetItems() []interface{} {
	return currentItems
}

// GetFriends returns the friends list cached from op_result.
func GetFriends() []interface{} {
	return currentFriends
}

// GetSect returns the cached sect info.
func GetSect() map[string]interface{} {
	return currentSect
}

// CacheFriends updates the friends cache from an op_result payload.
func CacheFriends(payload map[string]interface{}) {
	if f, ok := payload["friends"].([]interface{}); ok {
		currentFriends = f
	}
}

// CacheSect updates the sect cache from an op_result payload.
func CacheSect(payload map[string]interface{}) {
	if s, ok := payload["sect_id"].(string); ok && s != "" {
		currentSect = payload
	}
}

// GetEquippedItems returns items where equipped=true.
func GetEquippedItems() []map[string]interface{} {
	items := GetItems()
	var equipped []map[string]interface{}
	for _, it := range items {
		if item, ok := it.(map[string]interface{}); ok {
			if eq, ok := item["equipped"].(bool); ok && eq {
				equipped = append(equipped, item)
			}
		}
	}
	return equipped
}

// ── formatting helpers ──

func formatOpResult(payload map[string]interface{}) string {
	success, _ := payload["success"].(bool)
	message := getStr(payload, "message")
	tag := "OK"
	if !success {
		tag = "失败"
	}

	effects, _ := payload["effects"].(map[string]interface{})

	// 排行榜 entries 多行显示
	if effects != nil {
		if entries, ok := effects["entries"].([]interface{}); ok && len(entries) > 0 {
			title := message
			s := fmt.Sprintf("[%s] %s\n", tag, title)
			s += strings.Repeat("─", 60) + "\n"
			s += fmt.Sprintf("  %-4s %-20s %-10s %s\n", "排名", "名称", "数值", "境界")
			s += strings.Repeat("─", 60) + "\n"
			for _, e := range entries {
				em, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				rank := getIntDef(em, "rank")
				name := getStr(em, "name")
				val := getFloatDef(em, "value")
				realm := getStr(em, "realm")
				s += fmt.Sprintf("  #%-3d %-20s %-10.1f %s\n", rank, name, val, realmDisplay(realm))
			}
			s += strings.Repeat("─", 60)
			return s
		}

		// 邮件列表多行显示
		if mails, ok := effects["mails"].([]interface{}); ok && len(mails) > 0 {
			s := fmt.Sprintf("[%s] %s\n", tag, message)
			s += strings.Repeat("─", 60) + "\n"
			for _, m := range mails {
				mm, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				mid := getStr(mm, "id")
				title := getStr(mm, "title")
				sender := getStr(mm, "sender_name")
				isRead := false
				if r, ok := mm["is_read"].(bool); ok {
					isRead = r
				}
				hasAtt := false
				if a, ok := mm["has_attachment"].(bool); ok {
					hasAtt = a
				}
				status := ""
				if !isRead {
					status += "[未读]"
				}
				if hasAtt {
					status += "[附件]"
				}
				shortID := mid
				if len(shortID) > 8 {
					shortID = shortID[:8]
				}
				s += fmt.Sprintf("  %s %-8s %s (%s)%s\n", status, shortID, title, sender,
					map[bool]string{true: " ✔已领", false: ""}[getBoolDef(mm, "is_claimed")])
			}
			s += strings.Repeat("─", 60) + fmt.Sprintf("\n 共 %d 封", len(mails))
			return s
		}

		// 附近玩家 list
		if players, ok := effects["players"].([]interface{}); ok && len(players) > 0 {
			s := fmt.Sprintf("[%s] %s\n", tag, message)
			s += strings.Repeat("─", 50) + "\n"
			s += fmt.Sprintf("  %-20s %-15s %s\n", "名称", "境界", "神识")
			s += strings.Repeat("─", 50) + "\n"
			for _, p := range players {
				pm, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				name := getStr(pm, "name")
				realm := getStr(pm, "realm")
				spirit := getFloatDef(pm, "spirit")
				maxSp := getFloatDef(pm, "max_spirit")
				s += fmt.Sprintf("  %-20s %-15s %.0f/%.0f\n", name, realmDisplay(realm), spirit, maxSp)
			}
			s += strings.Repeat("─", 50) + fmt.Sprintf("\n 共 %d 人在此区域", len(players))
			return s
		}
	}

	s := fmt.Sprintf("[%s] %s", tag, message)

	if effects != nil && len(effects) > 0 {
		details := formatEffects(effects)
		if details != "" {
			s += " | " + details
		}

		// Cache friends list if present in effects
		if _, ok := effects["friends"]; ok {
			CacheFriends(effects)
		}
		// Cache sect info if present in effects
		if sectID, ok := effects["sect_id"].(string); ok && sectID != "" {
			CacheSect(effects)
		}
	}
	return s
}

func getBoolDef(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func formatEffects(effects map[string]interface{}) string {
	parts := []string{}

	if v, ok := effects["cultivation_gain"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("修为+%.2f", v))
	}
	if v, ok := effects["progress"].(float64); ok {
		parts = append(parts, fmt.Sprintf("进度%.1f%%", v))
	}
	if v, ok := effects["qi_cost"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("灵力-%.0f", v))
	}
	if v, ok := effects["qi_recovery"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("灵力+%.0f", v))
	}
	if v, ok := effects["spiritual_recovery"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("神识+%.0f", v))
	}
	if v, ok := effects["damage_dealt"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("伤害%.0f", v))
	}
	if v, ok := effects["damage"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("伤害%.0f", v))
	}
	if v, ok := effects["resource"].(string); ok {
		qty := 1.0
		if q, ok := effects["quantity"].(float64); ok {
			qty = q
		}
		parts = append(parts, fmt.Sprintf("%s x%.0f", v, qty))
	}
	if v, ok := effects["new_realm"].(string); ok {
		parts = append(parts, fmt.Sprintf("晋升%s", realmDisplay(v)))
	}
	if v, ok := effects["success_rate"].(float64); ok {
		parts = append(parts, fmt.Sprintf("成功率%.0f%%", v*100))
	}
	if v, ok := effects["sect_name"].(string); ok {
		parts = append(parts, fmt.Sprintf("宗门%s", v))
	}
	if v, ok := effects["method_quality"].(string); ok {
		parts = append(parts, fmt.Sprintf("品质%s", v))
	}
	if v, ok := effects["price"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("灵石%.0f", v))
	}
	if v, ok := effects["cost"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("消耗%.0f", v))
	}
	if v, ok := effects["target_name"].(string); ok {
		parts = append(parts, v)
	}
	if v, ok := effects["skill_exp"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("技能+%.0f", v))
	}
	if v, ok := effects["is_crit"].(bool); ok && v {
		parts = append(parts, "暴击!")
	}
	if v, ok := effects["is_dodge"].(bool); ok && v {
		parts = append(parts, "闪避!")
	}
	if v, ok := effects["discoveries"].([]interface{}); ok {
		for _, d := range v {
			if s, ok := d.(string); ok {
				parts = append(parts, s)
			}
		}
	}
	if v, ok := effects["item_name"].(string); ok {
		parts = append(parts, v)
	}
	if v, ok := effects["slot"].(string); ok {
		parts = append(parts, fmt.Sprintf("装备位:%s", v))
	}
	if v, ok := effects["spell_name"].(string); ok {
		parts = append(parts, fmt.Sprintf("法术:%s", v))
	}
	if v, ok := effects["count"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("共%d位", int(v)))
	}
	if v, ok := effects["member_count"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("成员%d人", int(v)))
	}
	if v, ok := effects["founder_name"].(string); ok {
		parts = append(parts, fmt.Sprintf("宗主:%s", v))
	}
	if v, ok := effects["quality"].(string); ok {
		parts = append(parts, fmt.Sprintf("品质:%s", v))
	}
	if v, ok := effects["element_affinity"].(string); ok && v != "" {
		parts = append(parts, fmt.Sprintf("亲和:%s", v))
	}
	if v, ok := effects["cultivation_multiplier"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("修炼倍率:%.1f", v))
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}

func formatStateSync(payload map[string]interface{}) string {
	entity, ok := payload["entity"].(map[string]interface{})
	if !ok {
		return "状态同步已接收"
	}

	name := getStr(entity, "name")
	realm := getStr(entity, "realm")
	status := getStr(entity, "status")

	attrs, _ := entity["attributes"].(map[string]interface{})

	s := fmt.Sprintf("%s | 境界:%s | 状态:%s", name, realmDisplay(realm), statusDisplay(status))

	if attrs != nil {
		if qi, ok := getFloat(attrs, "qi"); ok {
			maxQi, _ := getFloat(attrs, "max_qi")
			s += fmt.Sprintf(" | 灵力:%.0f/%.0f", qi, maxQi)
		}
		if sp, ok := getFloat(attrs, "spiritual_power"); ok {
			maxSp, _ := getFloat(attrs, "max_spiritual_power")
			s += fmt.Sprintf(" 神识:%.0f/%.0f", sp, maxSp)
		}
		if prog, ok := getFloat(attrs, "cultivation_progress"); ok {
			s += fmt.Sprintf(" | 修为:%.1f%%", prog)
		}
		if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
			if lg, _ := getInt64(ss, "low_grade"); lg > 0 {
				s += fmt.Sprintf(" | 灵石:%d", lg)
			}
		}
	}

	if pos, ok := entity["position"].(map[string]interface{}); ok {
		if rid, ok := pos["region_id"].(string); ok {
			s += fmt.Sprintf(" | 位置:%s", rid)
		}
	}

	return "≡ " + s
}

func realmDisplay(r string) string {
	switch r {
	case "mortal":
		return "凡人"
	case "qi_condensation":
		return "练气期"
	case "foundation":
		return "筑基期"
	case "golden_core":
		return "金丹期"
	case "nascent_soul":
		return "元婴期"
	case "soul_transformation":
		return "化神期"
	case "void_refinement":
		return "炼虚期"
	case "integration":
		return "合体期"
	case "mahayana":
		return "大乘期"
	case "tribulation":
		return "渡劫期"
	}
	return r
}

func statusDisplay(s string) string {
	switch s {
	case "normal":
		return "正常"
	case "cultivating":
		return "修炼中"
	case "combat":
		return "战斗中"
	case "resting":
		return "休息中"
	case "dead":
		return "已死亡"
	case "exploring":
		return "探索中"
	case "crafting":
		return "炼制中"
	case "meditating":
		return "打坐中"
	}
	return s
}

func getStr(m map[string]interface{}, key string) string {
	if s, ok := m[key].(string); ok {
		return s
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func getInt(m map[string]interface{}, key string) (int, bool) {
		v, ok := m[key]
		if !ok {
			return 0, false
		}
		switch n := v.(type) {
		case float64:
			return int(n), true
		case int:
			return n, true
		case int64:
			return int(n), true
		}
		return 0, false
	}

	func getIntDef(m map[string]interface{}, key string) int {
		v, ok := getInt(m, key)
		if !ok {
			return 0
		}
		return v
	}

	func getFloatDef(m map[string]interface{}, key string) float64 {
		v, ok := getFloat(m, key)
		if !ok {
			return 0
		}
		return v
	}

	func getInt64(m map[string]interface{}, key string) (int64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}
