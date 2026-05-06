package model

import (
	"fmt"
	"strings"
	"time"

	"cultivation-bubbletea/internal/client"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

// ── Message types for Bubble Tea ──

// wsMsg wraps a parsed WebSocket message for the tea program.
type wsMsg client.WsMessage

// errMsg wraps a fatal error.
type errMsg struct{ error }

// tickMsg is sent periodically for time-based updates (notifications, etc.).
type tickMsg time.Time

// ── Chat & combat log entries ──

// ChatEntry represents a single chat or system message.
type ChatEntry struct {
	Channel   string // "world", "private", "system", "event"
	Sender    string
	Content   string
	Timestamp time.Time
}

// CombatEntry represents a single combat log line.
type CombatEntry struct {
	Text      string
	Timestamp time.Time
	IsDamage  bool
	IsHeal    bool
	IsSystem  bool
}

// ── The main model ──

// Model is the top-level Bubble Tea model for the TUI client.
type Model struct {
	// Connection
	Conn      *websocket.Conn
	EntityID  string
	MsgCh     chan client.WsMessage
	State     *client.GameState
	Connected bool
	Err       error

	// Terminal
	Width  int
	Height int

	// UI components
	Input   textinput.Model
	MainVP  viewport.Model

	// Message logs
	ChatLog   []ChatEntry
	CombatLog []CombatEntry
	SysLog    []string

	// Notification
	Notification string
	NotifTimer   int

	// Command history
	CmdHistory []string
	HistIdx    int

	// Tabs & focus
	Focus   int // 0=input, 1=left panel, 2=right panel
	ChatTab int // 0=chat, 1=system, 2=combat

	// Inventory browsing mode
	InvMode   bool
	InvFilter int // 0=all, 1=equipment, 2=materials, 3=pills
	InvCursor int

	// Map view toggle
	MapMode bool

	// Loading indicator
	Loading bool

	// Input mode
	InputMode string // "", "chat"
}

// NewWithChannel creates a Model using external WS channel and state.
func NewWithChannel(conn *websocket.Conn, entityID string, msgCh chan client.WsMessage, state *client.GameState) Model {
	ti := textinput.New()
	ti.Placeholder = "输入命令 (输入 help 查看帮助)"
	ti.Focus()
	ti.Width = 80
	ti.Prompt = "> "

	vp := viewport.New(80, 20)
	vp.YPosition = 0

	return Model{
		Conn:      conn,
		EntityID:  entityID,
		MsgCh:     msgCh,
		State:     state,
		Connected: true,
		Input:     ti,
		MainVP:    vp,
		Focus:     0,
		ChatTab:   0,
		HistIdx:   -1,
	}
}


// Init initializes the Bubble Tea program.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForMsg(m.MsgCh),
		tickEvery(time.Second),
	)
}

// ── Commands ──

func waitForMsg(msgCh <-chan client.WsMessage) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-msgCh
		if !ok {
			return errMsg{fmt.Errorf("消息通道已关闭")}
		}
		return wsMsg(msg)
	}
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// SendActionCmd returns a tea.Cmd that sends a game action.
func SendActionCmd(conn *websocket.Conn, actionType string, params map[string]interface{}) tea.Cmd {
	return func() tea.Msg {
		err := client.SendAction(conn, actionType, params)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

// SendChatCmd returns a tea.Cmd that sends a chat message.
func SendChatCmd(conn *websocket.Conn, content, channel string) tea.Cmd {
	return func() tea.Msg {
		err := client.SendChat(conn, content, channel)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

// ── Helpers ──

// SetNotification sets a notification that auto-dismisses after ~5 ticks.
func (m *Model) SetNotification(msg string) {
	m.Notification = msg
	m.NotifTimer = 5
}

// AddChatMessage appends a chat entry to the log, keeping the last 500.
func (m *Model) AddChatMessage(channel, sender, content string) {
	m.ChatLog = append(m.ChatLog, ChatEntry{
		Channel:   channel,
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
	})
	if len(m.ChatLog) > 500 {
		m.ChatLog = m.ChatLog[len(m.ChatLog)-500:]
	}
}

// AddCombatEntry appends a combat log entry.
func (m *Model) AddCombatEntry(text string, isDamage, isHeal, isSystem bool) {
	m.CombatLog = append(m.CombatLog, CombatEntry{
		Text:      text,
		Timestamp: time.Now(),
		IsDamage:  isDamage,
		IsHeal:    isHeal,
		IsSystem:  isSystem,
	})
	if len(m.CombatLog) > 500 {
		m.CombatLog = m.CombatLog[len(m.CombatLog)-500:]
	}
}

// AddSystemMessage appends a system message.
func (m *Model) AddSystemMessage(text string) {
	m.SysLog = append(m.SysLog, text)
	if len(m.SysLog) > 200 {
		m.SysLog = m.SysLog[len(m.SysLog)-200:]
	}
}

// UpdateMainVPContent refreshes the viewport content based on the active tab.
func (m *Model) UpdateMainVPContent() {
	vpWidth := m.MainVP.Width
	if vpWidth > m.Width-60 {
		vpWidth = m.Width - 60
	}
	if vpWidth < 30 {
		vpWidth = 30
	}

	var content string
	switch m.ChatTab {
	case 0: // chat
		content = renderChatLog(m.ChatLog, vpWidth)
	case 1: // system
		content = renderSystemLog(m.SysLog, vpWidth)
	case 2: // combat
		content = renderCombatLog(m.CombatLog, vpWidth)
	}
	// Only auto-scroll if user was already at the bottom
	atBottom := m.MainVP.ScrollPercent() >= 0.99
	m.MainVP.SetContent(content)
	if atBottom {
		m.MainVP.GotoBottom()
	}
}

// ── Styling ──

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7C56DC"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	chatColor     = lipgloss.Color("#58A6FF")
	systemColor   = lipgloss.Color("#FFA657")
	combatColor   = lipgloss.Color("#FF6B6B")
	healColor     = lipgloss.Color("#56D364")
	damageColor   = lipgloss.Color("#FF4444")
	eventColor    = lipgloss.Color("#FFD700")
	realmColor    = lipgloss.Color("#DA8BFF")
	infoColor     = lipgloss.AdaptiveColor{Light: "#1A1A2E", Dark: "#C9D1D9"}
	dimColor      = lipgloss.AdaptiveColor{Light: "#888", Dark: "#666"}
	successColor  = lipgloss.AdaptiveColor{Light: "#28A745", Dark: "#3FB950"}
	errorColor    = lipgloss.AdaptiveColor{Light: "#D73A49", Dark: "#F85149"}
	borderColor   = lipgloss.AdaptiveColor{Light: "#CCC", Dark: "#444"}
	titleColor    = lipgloss.AdaptiveColor{Light: "#0366D6", Dark: "#58A6FF"}
	bgColor       = lipgloss.AdaptiveColor{Light: "#FFF", Dark: "#0D1117"}
)

// Element colors for spiritual roots
var elementColors = map[string]string{
	"金": "#FFD700",
	"木": "#56D364",
	"水": "#58A6FF",
	"火": "#FF4444",
	"土": "#BB8F4A",
	"雷": "#DA8BFF",
	"风": "#79C0FF",
	"冰": "#8EC8E8",
	"阴": "#BC8CFF",
	"阳": "#FFA657",
	"变异": "#FF6BC1",
}

// ── Render helpers ──

func renderChatLog(log []ChatEntry, width int) string {
	if len(log) == 0 {
		return dimStyle.Render("暂无消息")
	}
	var b strings.Builder
	start := len(log) - 100
	if start < 0 {
		start = 0
	}
	for _, e := range log[start:] {
		var line string
		switch e.Channel {
		case "world":
			line = fmt.Sprintf("[聊天] %s: %s", e.Sender, e.Content)
			line = lipgloss.NewStyle().Foreground(chatColor).Render(line)
		case "private":
			line = fmt.Sprintf("[私信] %s: %s", e.Sender, e.Content)
			line = lipgloss.NewStyle().Foreground(special).Render(line)
		case "system":
			line = lipgloss.NewStyle().Foreground(systemColor).Render("[系统] " + e.Content)
		case "event":
			line = lipgloss.NewStyle().Foreground(eventColor).Render("[事件] " + e.Content)
		default:
			line = e.Content
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func renderSystemLog(log []string, width int) string {
	if len(log) == 0 {
		return dimStyle.Render("暂无系统消息")
	}
	var b strings.Builder
	start := len(log) - 100
	if start < 0 {
		start = 0
	}
	for _, s := range log[start:] {
		b.WriteString(lipgloss.NewStyle().Foreground(systemColor).Render(s))
		b.WriteByte('\n')
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func renderCombatLog(log []CombatEntry, width int) string {
	if len(log) == 0 {
		return dimStyle.Render("暂无战斗记录")
	}
	var b strings.Builder
	start := len(log) - 100
	if start < 0 {
		start = 0
	}
	for _, e := range log[start:] {
		var styled string
		switch {
		case e.IsDamage:
			styled = lipgloss.NewStyle().Foreground(damageColor).Render(e.Text)
		case e.IsHeal:
			styled = lipgloss.NewStyle().Foreground(healColor).Render(e.Text)
		case e.IsSystem:
			styled = lipgloss.NewStyle().Foreground(systemColor).Render(e.Text)
		default:
			styled = e.Text
		}
		b.WriteString(styled)
		b.WriteByte('\n')
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

// ── Common styles ──

var (
	dimStyle    = lipgloss.NewStyle().Foreground(dimColor)
	titleStyle  = lipgloss.NewStyle().Foreground(titleColor).Bold(true)
	infoStyle   = lipgloss.NewStyle().Foreground(infoColor)
	successStyle = lipgloss.NewStyle().Foreground(successColor)
	errorStyle  = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(borderColor)
)
