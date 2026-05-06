package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the complete TUI layout.
func (m Model) View() string {
	if m.Err != nil {
		return errorStyle.Render(fmt.Sprintf("发生错误: %v\n按 Ctrl+C 退出", m.Err))
	}

	// Build panels
	leftPanel := m.renderLeftPanel()
	centerPanel := m.renderCenterPanel()
	rightPanel := m.renderRightPanel()

	// Build bottom bar
	bottomBar := m.renderBottomBar()

	// Build notification overlay
	notification := m.renderNotification()

	// Arrange main layout: three columns
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		centerPanel,
		rightPanel,
	)

	// Inventory or map overlay
	var overlay string
	if m.InvMode {
		overlay = m.renderInventoryOverlay()
	} else if m.MapMode {
		overlay = m.renderMapView()
	}

	// Stack: notification, main content, bottom bar
	var final string
	if notification != "" {
		final = notification + "\n"
	}
	final += mainContent + "\n"
	final += bottomBar

	// If there's an overlay, replace viewport area with overlay
	if overlay != "" {
		overlayBox := borderStyle.
			Width(m.Width-4).
			Height(m.Height-6).
			Render(overlay)
		final = overlayBox
	}

	return final
}

func (m Model) renderNotification() string {
	if m.Notification == "" {
		return ""
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#FFA500")).
		Foreground(lipgloss.Color("#000")).
		Bold(true).
		Width(m.Width).
		Padding(0, 1).
		Render(" " + m.Notification)
}

func (m Model) renderLeftPanel() string {
	entity := m.State.EntityCopy()
	width := 28
	if m.Width > 0 {
		w := m.Width / 5
		if w > 30 {
			w = 30
		}
		width = w
	}
	if width < 24 {
		width = 24
	}

	panelHeight := m.Height - 5
	if panelHeight < 10 {
		panelHeight = 10
	}

	var content strings.Builder

	// Character name & realm
	if entity != nil {
		name := getStr(entity, "name")
		realm := realmDisplay(getStr(entity, "realm"))
		status := statusDisplay(getStr(entity, "status"))

		content.WriteString(titleStyle.Render("⦿ " + name))
		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().Foreground(realmColor).Render("境界: " + realm))
		content.WriteString("\n")

		attrs, _ := entity["attributes"].(map[string]interface{})
		if attrs != nil {
			// HP bar
			hp := getFloatDef(attrs, "hp")
			maxHp := getFloatDef(attrs, "max_hp")
			if maxHp > 0 {
				content.WriteString(renderBar("气血", hp, maxHp, damageColor))
				content.WriteString("\n")
			}

			// Qi bar
			qi := getFloatDef(attrs, "qi")
			maxQi := getFloatDef(attrs, "max_qi")
			if maxQi > 0 {
				content.WriteString(renderBar("灵力", qi, maxQi, chatColor))
				content.WriteString("\n")
			}

			// Spiritual power bar
			sp := getFloatDef(attrs, "spiritual_power")
			maxSp := getFloatDef(attrs, "max_spiritual_power")
			if maxSp > 0 {
				content.WriteString(renderBar("神识", sp, maxSp, special))
				content.WriteString("\n")
			}

			// Cultivation progress bar
			prog := getFloatDef(attrs, "cultivation_progress")
			if prog > 0 {
				content.WriteString(renderBar("修为", prog, 100, highlight))
				content.WriteString("\n")
			}

			// Spirit stones
			if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
				if lg, _ := getInt64(ss, "low_grade"); lg > 0 {
					content.WriteString(infoStyle.Render(fmt.Sprintf("灵石: %d", lg)))
					content.WriteString("\n")
				}
			}

			// Status
			content.WriteString(fmt.Sprintf("状态: %s", statusDisplay(status)))
			content.WriteString("\n")
		}
	} else {
		content.WriteString(dimStyle.Render("等待角色数据..."))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Spiritual roots
	if entity != nil {
		if roots, ok := entity["spiritual_roots"].([]interface{}); ok && len(roots) > 0 {
			content.WriteString(titleStyle.Render("━ 灵根 ━"))
			content.WriteString("\n")
			for _, r := range roots {
				if root, ok := r.(map[string]interface{}); ok {
					elem := getStr(root, "element")
					purity := getFloatDef(root, "purity")
					colorHex := elementColors[elem]
					if colorHex == "" {
						colorHex = "#C9D1D9"
					}
					isMain, _ := root["is_main"].(bool)
					mark := ""
					if isMain {
						mark = " ★"
					}
					line := fmt.Sprintf("  %s %.0f%%%s", elem, purity, mark)
					content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorHex)).Render(line))
					content.WriteString("\n")
				}
			}
		}
	}

	content.WriteString("\n")

	// Equipped items
	items := m.State.ItemsCopy()
	var equipped []map[string]interface{}
	for _, it := range items {
		if item, ok := it.(map[string]interface{}); ok {
			if eq, ok := item["equipped"].(bool); ok && eq {
				equipped = append(equipped, item)
			}
		}
	}
	if len(equipped) > 0 {
		content.WriteString(titleStyle.Render("━ 装备 ━"))
		content.WriteString("\n")
		for _, eq := range equipped {
			name := getStr(eq, "name")
			slot := getStr(eq, "slot")
			content.WriteString(fmt.Sprintf("  [%s] %s\n", slot[:min(len(slot), 4)], name))
		}
		content.WriteString("\n")
	}

	// Main method
	if entity != nil {
		if method, ok := entity["main_method"].(map[string]interface{}); ok {
			content.WriteString(titleStyle.Render("━ 功法 ━"))
			content.WriteString("\n")
			mname := getStr(method, "name")
			mquality := getStr(method, "quality")
			if mname != "" {
				content.WriteString(fmt.Sprintf("  %s [%s]\n", mname, mquality))
			}
		}
	}

	// Pad to height
	rendered := lipgloss.NewStyle().
		Width(width).
		Height(panelHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(content.String())

	return rendered
}

func (m Model) renderCenterPanel() string {
	width := m.Width - 60
	height := m.Height - 5
	if width < 30 {
		width = 30
	}
	if height < 10 {
		height = 10
	}

	m.MainVP.Width = width
	m.MainVP.Height = height - 2 // account for tab bar

	// Tab bar
	tabTitles := []string{"聊天", "系统", "战斗"}
	var tabCells []string
	for i, t := range tabTitles {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.ChatTab {
			style = style.Foreground(titleColor).Bold(true).Underline(true)
		} else {
			style = style.Foreground(dimColor)
		}
		tabCells = append(tabCells, style.Render(t))
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabCells...)
	tabBar = lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(borderColor).
		Render(tabBar)

	// Viewport content
	vpContent := m.MainVP.View()

	content := tabBar + "\n" + vpContent

	// Pad to height
	leftWidth := 30
	if m.Width/5 < 30 {
		leftWidth = m.Width / 5
	}
	if leftWidth < 24 {
		leftWidth = 24
	}
	centerWidth := m.Width - leftWidth - 26

	return lipgloss.NewStyle().
		Width(centerWidth).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(content)
}

func (m Model) renderRightPanel() string {
	width := 24
	height := m.Height - 5
	if height < 10 {
		height = 10
	}

	var content strings.Builder

	// Position info
	entity := m.State.EntityCopy()
	if entity != nil {
		content.WriteString(titleStyle.Render("位置"))
		content.WriteString("\n")
		if pos, ok := entity["position"].(map[string]interface{}); ok {
			if rid, ok := pos["region_id"].(string); ok {
				rname := regionDisplay(rid)
				content.WriteString(fmt.Sprintf("  %s\n", rname))
			}
		}
		content.WriteString("\n")
	}

	// Spell/Combat info
	content.WriteString(titleStyle.Render("快捷信息"))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  Tab 切换面板"))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  ↑↓ 历史命令"))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  PgUp/PgDn 滚动"))
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  Esc 返回/取消"))
	content.WriteString("\n")
	content.WriteString("\n")

	// Friend count
	friends := m.State.Friends
	if friends != nil {
		content.WriteString(fmt.Sprintf("好友: %d 人\n", len(friends)))
	}

	// Item count
	items := m.State.ItemsCopy()
	if items != nil {
		content.WriteString(fmt.Sprintf("背包: %d 件\n", len(items)))
	}

	// Connected status
	if m.Connected {
		content.WriteString(successStyle.Render("● 已连接"))
	} else {
		content.WriteString(errorStyle.Render("○ 已断开"))
	}
	content.WriteString("\n")

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(content.String())
}

func (m Model) renderBottomBar() string {
	entity := m.State.EntityCopy()

	var statusParts []string
	if entity != nil {
		realm := realmDisplay(getStr(entity, "realm"))
		if realm != "" {
			statusParts = append(statusParts, realm)
		}

		attrs, _ := entity["attributes"].(map[string]interface{})
		if attrs != nil {
			qi := getFloatDef(attrs, "qi")
			maxQi := getFloatDef(attrs, "max_qi")
			statusParts = append(statusParts, fmt.Sprintf("灵%.0f/%.0f", qi, maxQi))

			sp := getFloatDef(attrs, "spiritual_power")
			maxSp := getFloatDef(attrs, "max_spiritual_power")
			statusParts = append(statusParts, fmt.Sprintf("神%.0f/%.0f", sp, maxSp))
		}

		if pos, ok := entity["position"].(map[string]interface{}); ok {
			if rid, ok := pos["region_id"].(string); ok {
				statusParts = append(statusParts, regionDisplay(rid))
			}
		}

		status := getStr(entity, "status")
		if status == "combat" {
			statusParts = append(statusParts, "⚔交战")
		}
	}

	statusStr := strings.Join(statusParts, " ")
	if statusStr != "" {
		statusStr = " " + statusStr
	}

	// Input prompt
	prompt := m.Input.View()

	// Combine: status bar on the left, input takes rest of line
	// Actually, for Bubble Tea, input is handled separately. We show the status bar and
	// the input line below it.
	statusBar := lipgloss.NewStyle().
		Width(m.Width).
		Background(lipgloss.Color("#1A1A2E")).
		Foreground(lipgloss.Color("#8B949E")).
		Render(statusStr)

	inputLine := lipgloss.NewStyle().
		Width(m.Width).
		Render(prompt)

	return statusBar + "\n" + inputLine
}

func (m Model) renderInventoryOverlay() string {
	filtered := m.getFilteredItems()
	width := m.Width - 10
	height := m.Height - 8
	if width < 40 {
		width = 40
	}
	if height < 10 {
		height = 10
	}

	var content strings.Builder

	// Tabs
	filters := []string{"全部", "装备", "材料", "丹药"}
	var tabs []string
	for i, f := range filters {
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == m.InvFilter {
			style = style.Foreground(titleColor).Bold(true).Underline(true)
		} else {
			style = style.Foreground(dimColor)
		}
		tabs = append(tabs, style.Render(f))
	}
	content.WriteString(strings.Join(tabs, " "))
	content.WriteString("\n")
	content.WriteString(strings.Repeat("─", width-4) + "\n")

	if len(filtered) == 0 {
		content.WriteString(dimStyle.Render("  背包为空"))
	} else {
		for i, item := range filtered {
			name := getStr(item, "name")
			qty := getFloatDef(item, "quantity")
			slot := getStr(item, "slot")

			cursor := "  "
			if i == m.InvCursor {
				cursor = "→ "
			}

			equipped := ""
			if getBoolDef(item, "equipped") {
				equipped = " [已装备]"
			}

			slotStr := ""
			if slot != "" {
				slotStr = fmt.Sprintf(" [%s]", slot)
			}

			qtyStr := ""
			if qty > 1 {
				qtyStr = fmt.Sprintf(" x%.0f", qty)
			}

			line := fmt.Sprintf("%s%s%s%s%s", cursor, name, slotStr, equipped, qtyStr)

			if i == m.InvCursor {
				content.WriteString(lipgloss.NewStyle().Foreground(highlight).Render(line))
			} else {
				content.WriteString(line)
			}
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  ↑↓ 选择  Tab 切换分类  Esc 返回"))

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(highlight).
		Padding(1, 2).
		Align(lipgloss.Left).
		Render(content.String())
}

func (m Model) renderMapView() string {
	entity := m.State.EntityCopy()
	width := m.Width - 10
	height := m.Height - 8

	var content strings.Builder
	content.WriteString(titleStyle.Render("世界地图"))
	content.WriteString("\n\n")

	currentRegion := ""
	if entity != nil {
		if pos, ok := entity["position"].(map[string]interface{}); ok {
			currentRegion = getStr(pos, "region_id")
		}
	}

	// ASCII map
	mapStr := `
        ┌───────────── 天渊禁地 ─────────────┐
        │                                     │
  玄冰谷 ─── 灵雾山脉 ─── 天机城 ─── 落日沙漠
        │        │              │
        │   青云镇          幽冥谷
        │        │
        └──── 中州城 ──── 万妖山脉 ──── 东海
`

	content.WriteString(mapStr)
	content.WriteString("\n")

	// Show connected regions
	content.WriteString(titleStyle.Render("相邻区域"))
	content.WriteString("\n")
	switch currentRegion {
	case "qingyun_town":
		content.WriteString("  灵雾山脉 (north)\n")
		content.WriteString("  中州城 (south)\n")
	case "mist_mountains":
		content.WriteString("  青云镇 (south)\n  玄冰谷 (north)\n  天机城 (east)\n")
	case "zhongzhou":
		content.WriteString("  青云镇 (north)\n  万妖山脉 (east)\n")
	default:
		content.WriteString(dimStyle.Render("  未知区域\n"))
	}

	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  当前位置: "))
	if currentRegion != "" {
		content.WriteString(regionDisplay(currentRegion))
	}
	content.WriteString("\n")
	content.WriteString(dimStyle.Render("  Esc 关闭地图"))

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(special).
		Padding(1, 2).
		Align(lipgloss.Left).
		Render(content.String())
}

// ── Render helpers ──

func renderBar(label string, current, max float64, color lipgloss.TerminalColor) string {
	width := 16
	filled := int(current / max * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := "["
	bar += strings.Repeat("█", filled)
	bar += strings.Repeat("░", width-filled)
	bar += "]"

	styled := lipgloss.NewStyle().Foreground(color).Render(bar)
	return fmt.Sprintf("%s %s %.0f/%.0f", label, styled, current, max)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func regionDisplay(rid string) string {
	switch rid {
	case "qingyun_town":
		return "青云镇"
	case "mist_mountains":
		return "灵雾山脉"
	case "zhongzhou":
		return "中州城"
	case "wan_yao_mountains":
		return "万妖山脉"
	case "xuan_bing_valley":
		return "玄冰谷"
	case "tianji_city":
		return "天机城"
	case "luori_desert":
		return "落日沙漠"
	case "mingyou_valley":
		return "幽冥谷"
	case "tianyuan_forbidden":
		return "天渊禁地"
	case "donghai":
		return "东海"
	}
	return rid
}
