package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chzyer/readline"
)

// completionFunc is set from main.go to avoid circular import with commands package.
var completionFunc func(prefix string) []string

// SetCompletionFunc allows the commands package to provide tab-completion without circular imports.
func SetCompletionFunc(fn func(prefix string) []string) {
	completionFunc = fn
}

// GetStatus returns the entity's status string for context awareness.
func GetStatus() string {
	entity := GetCharacter()
	if entity == nil {
		return ""
	}
	return getStr(entity, "status")
}

var globalRL *readline.Instance

// SetReadline sets the readline instance to reuse (called from main after auth).
func SetReadline(rl *readline.Instance) {
	globalRL = rl
}

// Screen manages the terminal display with split-screen layout.
type Screen struct {
	rl      *readline.Instance
	msgCh   chan string
	doneCh  <-chan struct{}
	once    sync.Once
	ownRL   bool
}

// NewScreen creates a new Screen with tab-completion and context-aware prompt.
func NewScreen(msgCh chan string, doneCh <-chan struct{}) (*Screen, error) {
	rl := globalRL
	if rl == nil {
		var err error
		rl, err = readline.NewEx(&readline.Config{
			Prompt:          "> ",
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			HistoryFile:     ".cli_history",
		})
		if err != nil {
			return nil, fmt.Errorf("readline: %w", err)
		}
	}
	globalRL = nil

	rl.Config.AutoComplete = readline.NewPrefixCompleter(
		readline.PcItemDynamic(func(prefix string) []string {
			if completionFunc != nil {
				return completionFunc(prefix)
			}
			return nil
		}),
	)

	rl.SetPrompt("> ")
	return &Screen{rl: rl, msgCh: msgCh, doneCh: doneCh, ownRL: globalRL == nil}, nil
}

// Close closes the screen and restores the terminal.
func (s *Screen) Close() {
	if s.ownRL {
		s.rl.Close()
	}
}

// WriteLine prints a message above the current prompt.
func (s *Screen) WriteLine(msg string) {
	// Truncate very long lines
	if len(msg) > 200 {
		msg = msg[:200] + "…"
	}
	s.rl.Stdout().Write([]byte(msg + "\n"))
}

// Start enters the main readline loop, calling handler for each input line.
func (s *Screen) Start(handler func(string)) {
	// Background: consume msgCh → display above prompt
	go func() {
		for msg := range s.msgCh {
			s.WriteLine(msg)
		}
	}()

	// Send initial status after entity loads
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.updatePrompt()
	}()

	for {
		s.updatePrompt()
		line, err := s.rl.Readline()
		if err != nil {
			break // EOF or Ctrl+C
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		handler(line)
	}
}

func (s *Screen) updatePrompt() {
	entity := GetCharacter()
	if entity == nil {
		s.rl.SetPrompt("> ")
		return
	}

	realm := realmDisplay(getStr(entity, "realm"))
	attrs, _ := entity["attributes"].(map[string]interface{})

	var parts []string
	if realm != "" {
		parts = append(parts, realm)
	}
	if attrs != nil {
		if qi, ok := getFloat(attrs, "qi"); ok {
			maxQi, _ := getFloat(attrs, "max_qi")
			parts = append(parts, fmt.Sprintf("灵%.0f/%.0f", qi, maxQi))
		}
		if sp, ok := getFloat(attrs, "spiritual_power"); ok {
			maxSp, _ := getFloat(attrs, "max_spiritual_power")
			parts = append(parts, fmt.Sprintf("神%.0f/%.0f", sp, maxSp))
		}
	}
	if pos, ok := entity["position"].(map[string]interface{}); ok {
		if rid, ok := pos["region_id"].(string); ok {
			parts = append(parts, rid)
		}
	}

	status := getStr(entity, "status")
	if status == "combat" {
		parts = append(parts, "⚔交战")
	}

	if len(parts) == 0 {
		s.rl.SetPrompt("> ")
	} else {
		s.rl.SetPrompt(fmt.Sprintf("[%s] > ", strings.Join(parts, " ")))
	}
	s.rl.Refresh()
}

