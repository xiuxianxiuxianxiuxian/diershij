package components

import (
	"cultivation-client/internal/gui/theme"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Sidebar struct {
	Width      unit.Dp
	Items      []SidebarItem
	Selected   int
	clickables []widget.Clickable
}

type SidebarItem struct {
	ID    string
	Label string
	Icon  []byte
}

func NewSidebar() *Sidebar {
	s := &Sidebar{
		Width: 200,
		Items: []SidebarItem{
			{ID: "character", Label: "角色"},
			{ID: "combat", Label: "战斗"},
			{ID: "world", Label: "世界"},
			{ID: "social", Label: "社交"},
			{ID: "settings", Label: "设置"},
		},
		Selected: 0,
	}
	s.clickables = make([]widget.Clickable, len(s.Items))
	return s
}

func (s *Sidebar) Select(index int) {
	if index >= 0 && index < len(s.Items) {
		s.Selected = index
	}
}

func (s *Sidebar) SelectedID() string {
	if s.Selected >= 0 && s.Selected < len(s.Items) {
		return s.Items[s.Selected].ID
	}
	return ""
}

func (s *Sidebar) Layout(gtx layout.Context, selectedID string, onSelect func(string)) layout.Dimensions {
	// 先绘制整个侧边栏背景
	drawRect(gtx, theme.DefaultTheme.Background, gtx.Constraints.Max)

	// 构建所有菜单项
	rigidChildren := make([]layout.FlexChild, 0, len(s.Items)*2)
	for i := range s.Items {
		idx := i
		rigidChildren = append(rigidChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.layoutItem(gtx, idx, s.Items[idx].ID == selectedID, selectedID, onSelect)
		}))
		if i < len(s.Items)-1 {
			rigidChildren = append(rigidChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Max
				sz.Y = gtx.Dp(unit.Dp(1))
				return drawRect(gtx, theme.DefaultTheme.Border, sz)
			}))
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, rigidChildren...)
}

func (s *Sidebar) layoutItem(gtx layout.Context, index int, isSelected bool, selectedID string, onSelect func(string)) layout.Dimensions {
	item := &s.Items[index]
	click := &s.clickables[index]

	if click.Clicked(gtx) {
		s.Selected = index
		if onSelect != nil {
			onSelect(item.ID)
		}
	}

	bg := theme.DefaultTheme.Surface
	if isSelected {
		bg = theme.DefaultTheme.Active
	} else if click.Hovered() {
		bg = theme.DefaultTheme.Hover
	}

	minH := gtx.Dp(unit.Dp(48))
	gtx.Constraints.Min.Y = minH

	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 固定高度，确保每个项目高度一致
		gtx.Constraints.Min.Y = minH
		// 背景
		drawRect(gtx, bg, gtx.Constraints.Max)
		// 文字
		return layout.Inset{
			Top:    unit.Dp(12),
			Left:   unit.Dp(16),
			Bottom: unit.Dp(12),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), item.Label)
			if isSelected {
				lbl.Color = toNRGBA(theme.DefaultTheme.Text)
			} else {
				lbl.Color = toNRGBA(theme.DefaultTheme.TextSecondary)
			}
			return lbl.Layout(gtx)
		})
	})
}
