package commands

import (
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		input string
		want  string // first name of resolved node, or "" if nil
	}{
		{"角色 状态", "状态"},
		{"角色 属性", "属性"},
		{"角色 技能", "技能"},
		{"status", "status"},
		{"help", "帮助"},
		{"修炼 修炼", "修炼"},
		{"cultivate", "cultivate"},
		{"战斗 目标", "目标"},
		{"combat", "combat"},
		{"战斗 逃跑", "逃跑"},
		{"背包 列表", "列表"},
		{"社交 好友 添加", "添加"},
		{"社交 聊天", "聊天"},
		{"买卖", ""},
		{"unknown", ""},
	}
	for _, tc := range tests {
		tokens := split(tc.input)
		node, remaining := resolve(tokens)
		if tc.want == "" {
			if node != nil {
				t.Errorf("resolve(%q) = %v, want nil", tc.input, node.Names[0])
			}
			continue
		}
		if node == nil {
			t.Errorf("resolve(%q) = nil, want %s", tc.input, tc.want)
			continue
		}
		if node.Names[0] != tc.want {
			t.Errorf("resolve(%q) = %s, want %s (remaining=%v)", tc.input, node.Names[0], tc.want, remaining)
		}
	}
}

func split(s string) []string {
	var r []string
	start := -1
	for i, c := range s {
		if c == ' ' {
			if start >= 0 {
				r = append(r, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		r = append(r, s[start:])
	}
	return r
}
