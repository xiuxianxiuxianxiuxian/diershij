package model

import (
	"fmt"
	"strings"
	"time"

	"cultivation-bubbletea/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
)

// Update processes all messages and user interactions.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// ── Window resize ──
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Input.Width = msg.Width - 4
		m.MainVP.Width = msg.Width - 60
		m.MainVP.Height = msg.Height - 6
		if m.MainVP.Width < 30 {
			m.MainVP.Width = 30
		}
		m.UpdateMainVPContent()
		return m, nil

	// ── WebSocket messages ──
	case wsMsg:
		cmds = append(cmds, waitForMsg(m.MsgCh))

		// Handle different message types
		switch msg.Type {
		case "state_sync":
			// State already cached by ws.go reader
			m.UpdateMainVPContent()
			m.SetNotification("状态已同步")

		case "entity_update":
			m.UpdateMainVPContent()

		case "op_result":
			m.handleOpResult(msg.Payload)

		case "error":
			errMsg := getStr(msg.Payload, "message")
			m.AddSystemMessage("[错误] " + errMsg)
			m.SetNotification("[错误] " + errMsg)
			m.UpdateMainVPContent()

		case "chat":
			sender := getStr(msg.Payload, "sender_name")
			if sender == "" {
				sender = getStr(msg.Payload, "sender_id")
			}
			content := getStr(msg.Payload, "content")
			if content != "" {
				m.AddChatMessage("world", sender, content)
				m.UpdateMainVPContent()
			}

		case "world_event":
			desc := getStr(msg.Payload, "description")
			m.AddChatMessage("event", "", desc)
			m.SetNotification("[世界事件] " + desc)
			m.UpdateMainVPContent()

		case "system", "announcement":
			content := getStr(msg.Payload, "message")
			if content == "" {
				content = getStr(msg.Payload, "content")
			}
			m.AddSystemMessage(content)
			m.SetNotification(content)
			m.UpdateMainVPContent()

		case "new_message":
			sender := getStr(msg.Payload, "sender_name")
			content := getStr(msg.Payload, "content")
			m.AddChatMessage("private", sender, content)
			m.SetNotification(fmt.Sprintf("[私信] %s: %s", sender, content))
			m.UpdateMainVPContent()

		case "friend_request":
			from := getStr(msg.Payload, "from_name")
			msg := fmt.Sprintf("%s 请求添加你为好友", from)
			m.AddSystemMessage(msg)
			m.SetNotification(msg)
			m.UpdateMainVPContent()
		}
		return m, tea.Batch(cmds...)

	// ── Tick (every second) ──
	case tickMsg:
		cmds = append(cmds, tickEvery(time.Second))
		// Decrement notification timer
		if m.NotifTimer > 0 {
			m.NotifTimer--
			if m.NotifTimer == 0 {
				m.Notification = ""
			}
		}
		return m, tea.Batch(cmds...)

	// ── Error ──
	case errMsg:
		m.Err = msg
		m.Connected = false
		return m, tea.Quit

	// ── Key press ──
	case tea.KeyMsg:
		// Handle global keybinds first
		switch msg.String() {
		case "ctrl+c", "ctrl+d":
			return m, tea.Quit

		case "tab":
			if m.InvMode {
				m.InvFilter = (m.InvFilter + 1) % 4
				return m, nil
			}
			if m.MapMode {
				m.MapMode = false
				return m, nil
			}
			m.ChatTab = (m.ChatTab + 1) % 3
			m.UpdateMainVPContent()
			return m, nil

		case "esc":
			if m.InvMode {
				m.InvMode = false
				return m, nil
			}
			if m.MapMode {
				m.MapMode = false
				return m, nil
			}
			m.Input.Reset()
			m.Input.Focus()
			m.Focus = 0
			return m, nil

		case "up":
			if m.InvMode {
				if m.InvCursor > 0 {
					m.InvCursor--
				}
				return m, nil
			}
			// Command history
			if m.HistIdx < len(m.CmdHistory)-1 && len(m.CmdHistory) > 0 {
				m.HistIdx++
				m.Input.SetValue(m.CmdHistory[len(m.CmdHistory)-1-m.HistIdx])
				m.Input.CursorEnd()
			}
			return m, nil

		case "down":
			if m.InvMode {
				items := m.getFilteredItems()
				if m.InvCursor < len(items)-1 {
					m.InvCursor++
				}
				return m, nil
			}
			if m.HistIdx > 0 {
				m.HistIdx--
				m.Input.SetValue(m.CmdHistory[len(m.CmdHistory)-1-m.HistIdx])
				m.Input.CursorEnd()
			} else if m.HistIdx == 0 {
				m.HistIdx = -1
				m.Input.SetValue("")
			}
			return m, nil

		case "pgup":
			m.MainVP.HalfPageUp()
			return m, nil
		case "pgdown":
			m.MainVP.HalfPageDown()
			return m, nil
		case "ctrl+up":
			m.MainVP.ScrollUp(3)
			return m, nil
		case "ctrl+down":
			m.MainVP.ScrollDown(3)
			return m, nil

		case "enter":
			line := strings.TrimSpace(m.Input.Value())
			m.Input.Reset()
			if line == "" {
				return m, nil
			}

			// Dispatch command
			m.CmdHistory = append(m.CmdHistory, line)
			m.HistIdx = -1

			// Parse and execute
			cmd := commands.Parse(line)
			if cmd == nil {
				m.AddSystemMessage("未知命令: " + strings.Fields(line)[0] + " (输入 help 查看帮助)")
				m.UpdateMainVPContent()
				return m, nil
			}

			switch cmd.Action {
			case "builtin":
				m.handleBuiltin(cmd)
			case "op":
				cmds = append(cmds, SendActionCmd(m.Conn, cmd.ActionType, cmd.Params))
			case "chat":
				cmds = append(cmds, SendChatCmd(m.Conn, cmd.Content, "world"))
			case "quit":
				return m, tea.Quit
			}
			return m, tea.Batch(cmds...)

		default:
			// Pass to input component
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)
			return m, cmd
		}
	}

	// Pass window size to input
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

// handleOpResult processes an operation result from the game server.
func (m *Model) handleOpResult(payload map[string]interface{}) {
	success, _ := payload["success"].(bool)
	message := getStr(payload, "message")

	tag := "OK"
	if !success {
		tag = "失败"
	}

	effects, _ := payload["effects"].(map[string]interface{})

	// Special rendering for large data sets (leaderboard, mail, nearby players)
	if effects != nil {
		if entries, ok := effects["entries"].([]interface{}); ok && len(entries) > 0 {
			m.renderLeaderboard(message, tag, entries)
			return
		}
		if mails, ok := effects["mails"].([]interface{}); ok && len(mails) > 0 {
			m.renderMailList(message, tag, mails)
			return
		}
		if players, ok := effects["players"].([]interface{}); ok && len(players) > 0 {
			m.renderNearbyPlayers(message, tag, players)
			return
		}
		if shopItems, ok := effects["items"].([]interface{}); ok && len(shopItems) > 0 {
			m.renderShopItems(message, tag, effects)
			return
		}
			if events, ok := effects["events"].([]interface{}); ok && len(events) > 0 {
				m.renderWorldEvents(message, tag, events)
				return
			}
			if auctions, ok := effects["auctions"].([]interface{}); ok && len(auctions) > 0 {
				m.renderAuctionList(message, tag, auctions)
				return
			}

		// Cache friends and sect
		if _, ok := effects["friends"]; ok {
			m.State.UpdateFriends(effects)
		}
		if sectID, ok := effects["sect_id"].(string); ok && sectID != "" {
			m.State.UpdateSect(effects)
		}

		// Combat log entries
		if hits, ok := effects["hits"].([]interface{}); ok {
			for _, h := range hits {
				if hm, ok := h.(map[string]interface{}); ok {
					dmg := getFloatDef(hm, "damage")
					target := getStr(hm, "target_name")
					var text string
					if getBoolDef(hm, "is_crit") {
						text = fmt.Sprintf("暴击! 对 %s 造成 %.0f 伤害", target, dmg)
						m.AddCombatEntry(text, true, false, false)
					} else if getBoolDef(hm, "is_dodge") {
						text = fmt.Sprintf("%s 闪避了攻击", target)
						m.AddCombatEntry(text, false, false, true)
					} else {
						text = fmt.Sprintf("对 %s 造成 %.0f 伤害", target, dmg)
						m.AddCombatEntry(text, true, false, false)
					}
				}
			}
		}

		// Parse effects for display
		effectsStr := formatEffects(effects)
		if effectsStr != "" {
			message = message + " | " + effectsStr
		}
	}

	line := fmt.Sprintf("[%s] %s", tag, message)
	if success {
		m.AddChatMessage("system", "", line)
	} else {
		m.AddSystemMessage(line)
	}
	m.SetNotification(line)
	m.UpdateMainVPContent()
}

func (m *Model) renderLeaderboard(title, tag string, entries []interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 50) + "\n")
	b.WriteString(fmt.Sprintf("  %-4s %-20s %-10s %s\n", "排名", "名称", "数值", "境界"))
	b.WriteString(strings.Repeat("─", 50) + "\n")
	for _, e := range entries {
		em, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		rank := getIntDef(em, "rank")
		name := getStr(em, "name")
		val := getFloatDef(em, "value")
		realm := realmDisplay(getStr(em, "realm"))
		b.WriteString(fmt.Sprintf("  #%-3d %-20s %-10.1f %s\n", rank, name, val, realm))
	}
	b.WriteString(strings.Repeat("─", 50))
	m.AddChatMessage("system", "", b.String())
}

func (m *Model) renderMailList(title, tag string, mails []interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 50) + "\n")
	for _, m1 := range mails {
		mm, ok := m1.(map[string]interface{})
		if !ok {
			continue
		}
		mid := getStr(mm, "id")
		t := getStr(mm, "title")
		sender := getStr(mm, "sender_name")
		shortID := mid
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		status := ""
		if !getBoolDef(mm, "is_read") {
			status += "[未读]"
		}
		if getBoolDef(mm, "has_attachment") {
			status += "[附件]"
		}
		b.WriteString(fmt.Sprintf("  %s %-8s %s (%s)%s\n", status, shortID, t, sender,
			map[bool]string{true: " ✔已领", false: ""}[getBoolDef(mm, "is_claimed")]))
	}
	b.WriteString(strings.Repeat("─", 50))
	m.AddChatMessage("system", "", b.String())
}

func (m *Model) renderNearbyPlayers(title, tag string, players []interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 40) + "\n")
	b.WriteString(fmt.Sprintf("  %-20s %-15s %s\n", "名称", "境界", "神识"))
	b.WriteString(strings.Repeat("─", 40) + "\n")
	for _, p := range players {
		pm, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		name := getStr(pm, "name")
		realm := realmDisplay(getStr(pm, "realm"))
		spirit := getFloatDef(pm, "spirit")
		maxSp := getFloatDef(pm, "max_spirit")
		b.WriteString(fmt.Sprintf("  %-20s %-15s %.0f/%.0f\n", name, realm, spirit, maxSp))
	}
	b.WriteString(strings.Repeat("─", 40))
	m.AddChatMessage("system", "", b.String())
}

func (m *Model) renderShopItems(title, tag string, effects map[string]interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 50) + "\n")
	if items, ok := effects["items"].([]interface{}); ok {
		for _, item := range items {
			if im, ok := item.(map[string]interface{}); ok {
				name := getStr(im, "name")
				price := getFloatDef(im, "price")
				qty := getFloatDef(im, "quantity")
				b.WriteString(fmt.Sprintf("  %-20s 灵石:%.0f 库存:%.0f\n", name, price, qty))
			}
		}
	}
	b.WriteString(strings.Repeat("─", 50))
	m.AddChatMessage("system", "", b.String())
}

func (m *Model) renderWorldEvents(title, tag string, events []interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 50) + "\n")
	for _, e := range events {
		if em, ok := e.(map[string]interface{}); ok {
			name := getStr(em, "name")
			desc := getStr(em, "description")
			region := getStr(em, "region_id")
			rname := regionDisplay(region)
			b.WriteString(fmt.Sprintf("  %s — %s [%s]\n", name, desc, rname))
		}
	}
	b.WriteString(strings.Repeat("─", 50))
	m.AddChatMessage("system", "", b.String())
}

func (m *Model) renderAuctionList(title, tag string, auctions []interface{}) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("[%s] %s\n", tag, title))
	b.WriteString(strings.Repeat("─", 55) + "\n")
	b.WriteString(fmt.Sprintf("  %-10s %-20s %-8s %s\n", "ID", "物品", "价格", "数量"))
	b.WriteString(strings.Repeat("─", 55) + "\n")
	for _, a := range auctions {
		if am, ok := a.(map[string]interface{}); ok {
			id := getStr(am, "id")
			itemName := getStr(am, "item_name")
			price := getFloatDef(am, "price")
			qty := getFloatDef(am, "quantity")
			shortID := id
			if len(shortID) > 8 {
				shortID = shortID[:8]
			}
			b.WriteString(fmt.Sprintf("  %-10s %-20s %-8.0f %.0f\n", shortID, itemName, price, qty))
		}
	}
	b.WriteString(strings.Repeat("─", 55))
	m.AddChatMessage("system", "", b.String())
}

// handleBuiltin processes client-side commands.
func (m *Model) handleBuiltin(cmd *commands.Command) {
	switch cmd.Name {
	case "help":
		m.renderHelp()
	case "clear", "cls":
		m.ChatLog = nil
		m.CombatLog = nil
		m.SysLog = nil
		m.MainVP.SetContent("")
	case "status", "st":
		m.renderStatus()
	case "attributes", "attrs":
		m.renderAttributes()
	case "skills":
		m.renderSkills()
	case "equip":
		m.renderEquippedItems()
	case "inventory", "bag":
		m.InvMode = true
		m.InvCursor = 0
		m.InvFilter = 0
	case "map":
		m.MapMode = true
	case "nearby", "near":
		// Will be handled by the op dispatch
	default:
		m.AddSystemMessage("未知命令: " + cmd.Raw)
		m.UpdateMainVPContent()
	}
}

func (m *Model) renderHelp() {
	helpText := `可用命令:
  help/清屏/退出 — 系统命令
  角色 (状态/属性/技能/装备) — 角色信息
  修炼 (修炼/打坐/休息/突破/探索/采集/移动/炼制/自创) — 修炼相关
  战斗 (攻击/自动/逃跑/技能/法术) — 战斗
  背包 (列表/使用/丢弃) — 物品管理
  社交 (聊天/私信/好友/宗门) — 社交
  商店 (列表/物品/购买/出售) — NPC商店
  拍卖 (列表/上架/购买/取消/查看) — 拍卖行
  装备 (列表/装备/卸下) — 装备管理
  功法 (列表/学习/主修) — 功法系统
  排行榜 (修为榜/战力榜/财富榜/功德榜) — 排行榜
  邮件 (列表/读取/领取/删除) — 信箱
  附近 — 查看附近玩家
  法术/技能 (学习/施展) — 法术系统
  世界事件/events — 查看活跃世界事件
  交易 — 玩家交易
  map — 区域地图
  Tab — 切换面板  ↑↓ — 历史  PgUp/PgDn — 滚动  Esc — 取消`
	m.AddChatMessage("system", "", helpText)
	m.UpdateMainVPContent()
}

func (m *Model) renderStatus() {
	entity := m.State.EntityCopy()
	if entity == nil {
		m.AddSystemMessage("暂无角色信息")
		m.UpdateMainVPContent()
		return
	}

	name := getStr(entity, "name")
	realm := realmDisplay(getStr(entity, "realm"))
	status := statusDisplay(getStr(entity, "status"))
	attrs, _ := entity["attributes"].(map[string]interface{})

	var b strings.Builder
	b.WriteString(fmt.Sprintf("角色: %s | 境界: %s | 状态: %s\n", name, realm, status))
	if attrs != nil {
		qi := getFloatDef(attrs, "qi")
		maxQi := getFloatDef(attrs, "max_qi")
		sp := getFloatDef(attrs, "spiritual_power")
		maxSp := getFloatDef(attrs, "max_spiritual_power")
		prog := getFloatDef(attrs, "cultivation_progress")
		hp := getFloatDef(attrs, "hp")
		maxHp := getFloatDef(attrs, "max_hp")
		age := getFloatDef(attrs, "age")
		lifespan := getFloatDef(attrs, "lifespan")
		b.WriteString(fmt.Sprintf("灵力: %.0f/%.0f  神识: %.0f/%.0f\n", qi, maxQi, sp, maxSp))
		b.WriteString(fmt.Sprintf("气血: %.0f/%.0f  修为进度: %.1f%%\n", hp, maxHp, prog))
		b.WriteString(fmt.Sprintf("寿元: %.0f/%.0f\n", age, lifespan))

		if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
			if lg, _ := getInt64(ss, "low_grade"); lg > 0 {
				b.WriteString(fmt.Sprintf("灵石: %d\n", lg))
			}
		}
	}
	m.AddChatMessage("system", "", b.String())
	m.UpdateMainVPContent()
}

func (m *Model) renderAttributes() {
	entity := m.State.EntityCopy()
	if entity == nil {
		m.AddSystemMessage("暂无角色数据")
		m.UpdateMainVPContent()
		return
	}
	attrs, _ := entity["attributes"].(map[string]interface{})
	if attrs == nil {
		m.AddSystemMessage("无属性数据")
		m.UpdateMainVPContent()
		return
	}

	var b strings.Builder
	b.WriteString("=== 详细属性 ===\n")
	for k, v := range attrs {
		if k == "spirit_stones" {
			continue
		}
		switch val := v.(type) {
		case float64:
			if val != 0 || k == "cultivation_progress" {
				b.WriteString(fmt.Sprintf("  %s: %.2f\n", k, val))
			}
		case string:
			if val != "" {
				b.WriteString(fmt.Sprintf("  %s: %s\n", k, val))
			}
		case bool:
			b.WriteString(fmt.Sprintf("  %s: %v\n", k, val))
		}
	}
	m.AddChatMessage("system", "", b.String())
	m.UpdateMainVPContent()
}

func (m *Model) renderSkills() {
	spells := m.State.SpellsCopy()
	if len(spells) == 0 {
		m.AddSystemMessage("无已学法术")
		m.UpdateMainVPContent()
		return
	}
	var b strings.Builder
	b.WriteString("=== 法术/技能 ===\n")
	for _, s := range spells {
		if sm, ok := s.(map[string]interface{}); ok {
			name := getStr(sm, "name")
			desc := getStr(sm, "description")
			if desc == "" {
				desc = getStr(sm, "spell_name")
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", name, desc))
		}
	}
	m.AddChatMessage("system", "", b.String())
	m.UpdateMainVPContent()
}

func (m *Model) renderEquippedItems() {
	items := m.State.ItemsCopy()
	var equipped []map[string]interface{}
	for _, it := range items {
		if item, ok := it.(map[string]interface{}); ok {
			if eq, ok := item["equipped"].(bool); ok && eq {
				equipped = append(equipped, item)
			}
		}
	}
	if len(equipped) == 0 {
		m.AddSystemMessage("未装备任何物品")
		m.UpdateMainVPContent()
		return
	}
	var b strings.Builder
	b.WriteString("=== 已装备 ===\n")
	for _, eq := range equipped {
		name := getStr(eq, "name")
		slot := getStr(eq, "slot")
		b.WriteString(fmt.Sprintf("  [%s] %s\n", slot, name))
	}
	m.AddChatMessage("system", "", b.String())
	m.UpdateMainVPContent()
}

// getFilteredItems returns items filtered by InvFilter.
func (m *Model) getFilteredItems() []map[string]interface{} {
	items := m.State.ItemsCopy()
	var filtered []map[string]interface{}
	for _, it := range items {
		if item, ok := it.(map[string]interface{}); ok {
			if m.InvFilter == 0 {
				filtered = append(filtered, item)
			} else if m.InvFilter == 1 && getStr(item, "slot") != "" {
				filtered = append(filtered, item)
			} else if m.InvFilter == 2 && getStr(item, "type") == "material" {
				filtered = append(filtered, item)
			} else if m.InvFilter == 3 && getStr(item, "type") == "potion" {
				filtered = append(filtered, item)
			}
		}
	}
	return filtered
}

// formatEffects converts common effect fields to a compact string.
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
	if v, ok := effects["damage_dealt"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("伤害%.0f", v))
	}
	if v, ok := effects["new_realm"].(string); ok {
		parts = append(parts, fmt.Sprintf("晋升%s", realmDisplay(v)))
	}
	if v, ok := effects["success_rate"].(float64); ok {
		parts = append(parts, fmt.Sprintf("成功率%.0f%%", v*100))
	}
	if v, ok := effects["resource"].(string); ok {
		qty := 1.0
		if q, ok := effects["quantity"].(float64); ok {
			qty = q
		}
		parts = append(parts, fmt.Sprintf("%s x%.0f", v, qty))
	}
	if v, ok := effects["price"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("灵石%.0f", v))
	}
	if v, ok := effects["cost"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("消耗%.0f", v))
	}
	if v, ok := effects["skill_exp"].(float64); ok && v > 0 {
		parts = append(parts, fmt.Sprintf("技能+%.0f", v))
	}
	if v, ok := effects["is_crit"].(bool); ok && v {
		parts = append(parts, "暴击!")
	}
	if v, ok := effects["spell_name"].(string); ok {
		parts = append(parts, fmt.Sprintf("法术:%s", v))
	}
	return strings.Join(parts, " ")
}

// ── Helpers ──

func getStr(m map[string]interface{}, key string) string {
	if s, ok := m[key].(string); ok {
		return s
	}
	return ""
}

func getFloatDef(m map[string]interface{}, key string) float64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return 0
}

func getIntDef(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
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

func getBoolDef(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
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
