package components

import (
	"image"
	"image/color"

	"cultivation-client/internal/gui/theme"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					return drawRect(gtx, theme.DefaultTheme.Background, gtx.Constraints.Max)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx, s.layoutItems(gtx, selectedID, onSelect)...)
				}),
			)
		}),
	)
}

func (s *Sidebar) layoutItems(gtx layout.Context, selectedID string, onSelect func(string)) []layout.FlexChild {
	children := make([]layout.FlexChild, 0, len(s.Items))

	for i := range s.Items {
		index := i
		item := &s.Items[i]

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			minHeight := gtx.Dp(unit.Dp(48))
			if gtx.Constraints.Min.Y < minHeight {
				gtx.Constraints.Min.Y = minHeight
			}
			return s.layoutMenuItem(gtx, index, item, selectedID, onSelect)
		}))
	}

	return children
}

func (s *Sidebar) layoutMenuItem(gtx layout.Context, index int, item *SidebarItem, selectedID string, onSelect func(string)) layout.Dimensions {
	isSelected := item.ID == selectedID
	isHovered := s.clickables[index].Hovered()

	bgColor := s.getBackgroundColor(isSelected, isHovered)

	return s.clickables[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if s.clickables[index].Clicked(gtx) {
			s.Selected = index
			if onSelect != nil {
				onSelect(item.ID)
			}
		}

		minHeight := gtx.Dp(unit.Dp(48))
		if gtx.Constraints.Min.Y < minHeight {
			gtx.Constraints.Min.Y = minHeight
		}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRect(gtx, bgColor, gtx.Constraints.Max)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(12),
					Bottom: unit.Dp(12),
					Left:   unit.Dp(16),
					Right:  unit.Dp(16),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), item.Label)
					if isSelected {
						lbl.Color = toNRGBA(theme.DefaultTheme.Text)
					} else {
						lbl.Color = toNRGBA(theme.DefaultTheme.TextSecondary)
					}
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}

func (s *Sidebar) getBackgroundColor(isSelected, isHovered bool) color.RGBA {
	if isSelected {
		return theme.DefaultTheme.Active
	}
	if isHovered {
		return theme.DefaultTheme.Hover
	}
	return theme.DefaultTheme.Surface
}

func drawRect(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
