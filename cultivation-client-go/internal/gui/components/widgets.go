package components

import (
	"fmt"
	"image"
	"image/color"

	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/types"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CharacterInfo struct {
	Name      string
	Level     int
	Health    int
	MaxHealth int
	Energy    int
	MaxEnergy int
	Attack    int
	Defense   int
	Speed     int
}

func NewCharacterInfo() *CharacterInfo {
	return &CharacterInfo{
		Name:      "",
		Level:     1,
		Health:    100,
		MaxHealth: 100,
		Energy:    50,
		MaxEnergy: 50,
		Attack:    10,
		Defense:   5,
		Speed:     8,
	}
}

func (c *CharacterInfo) SetCharacter(char *types.Character) {
	if char != nil {
		c.Name = char.Name
		c.Level = char.Level
		c.Health = char.Health
		c.MaxHealth = char.MaxHealth
		c.Energy = char.Energy
		c.MaxEnergy = char.MaxEnergy
		c.Attack = char.Attack
		c.Defense = char.Defense
		c.Speed = char.Speed
	}
}

func (c *CharacterInfo) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			nameLabel := NewLabel("名称: " + c.Name)
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return nameLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			levelLabel := NewLabel(fmt.Sprintf("等级: %d", c.Level))
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return levelLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hpLabel := NewLabel(fmt.Sprintf("生命值: %d/%d", c.Health, c.MaxHealth))
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return hpLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			energyLabel := NewLabel(fmt.Sprintf("能量: %d/%d", c.Energy, c.MaxEnergy))
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return energyLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			attackLabel := NewLabel(fmt.Sprintf("攻击: %d", c.Attack))
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return attackLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			defenseLabel := NewLabel(fmt.Sprintf("防御: %d", c.Defense))
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return defenseLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			speedLabel := NewLabel(fmt.Sprintf("速度: %d", c.Speed))
			return speedLabel.Layout(gtx)
		}),
	)
}

type Slider struct {
	widget  widget.Float
	min     float32
	max     float32
	changed bool
}

func NewSlider(value float32) *Slider {
	s := &Slider{
		min: 0,
		max: 1,
	}
	s.widget.Value = value
	return s
}

func (s *Slider) Value() float32 {
	return s.widget.Value
}

func (s *Slider) SetValue(v float32) {
	if v >= s.min && v <= s.max {
		s.widget.Value = v
	}
}

func (s *Slider) Changed() bool {
	return s.changed
}

func (s *Slider) Layout(gtx layout.Context) layout.Dimensions {
	if s.widget.Update(gtx) {
		s.changed = true
	}

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Max
					size.Y = gtx.Dp(unit.Dp(8))
					return drawRect(gtx, theme.DefaultTheme.Border, size)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Max
					size.Y = gtx.Dp(unit.Dp(8))
					size.X = int(float32(size.X) * s.widget.Value)
					return drawRect(gtx, theme.DefaultTheme.Primary, size)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					thumbSize := gtx.Dp(unit.Dp(20))
					thumbPos := int(float32(gtx.Constraints.Max.X-thumbSize) * s.widget.Value)
					return layout.Inset{
						Left: unit.Dp(thumbPos),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := image.Point{X: thumbSize, Y: thumbSize}
						return drawCircle(gtx, theme.DefaultTheme.PrimaryVariant, size)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			valueText := fmt.Sprintf("%.0f%%", s.widget.Value*100)
			lbl := material.Label(th, 14, valueText)
			lbl.Color = toNRGBA(theme.DefaultTheme.Text)
			return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return lbl.Layout(gtx)
			})
		}),
	)
}

type Checkbox struct {
	widget widget.Bool
	label  string
}

func NewCheckbox(label string, checked bool) *Checkbox {
	c := &Checkbox{
		label: label,
	}
	c.widget.Value = checked
	return c
}

func (c *Checkbox) Checked() bool {
	return c.widget.Value
}

func (c *Checkbox) SetChecked(checked bool) {
	c.widget.Value = checked
}

func (c *Checkbox) Changed(gtx layout.Context) bool {
	return c.widget.Update(gtx)
}

func (c *Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	return c.widget.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				boxSize := gtx.Dp(unit.Dp(20))
				return layout.Stack{Alignment: layout.Center}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						size := image.Point{X: boxSize, Y: boxSize}
						return drawRect(gtx, theme.DefaultTheme.Surface, size)
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						border := widget.Border{
							Color:        toNRGBA(theme.DefaultTheme.Border),
							CornerRadius: unit.Dp(4),
							Width:        unit.Dp(2),
						}
						return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							if c.widget.Value {
								checkSize := gtx.Dp(unit.Dp(12))
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									size := image.Point{X: checkSize, Y: checkSize}
									return drawCheckmark(gtx, theme.DefaultTheme.Primary, size)
								})
							}
							return layout.Dimensions{Size: image.Point{X: boxSize, Y: boxSize}}
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, 14, c.label)
				lbl.Color = toNRGBA(theme.DefaultTheme.Text)
				return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}

type Icon struct {
	Data  []byte
	Color color.RGBA
}

func NewIcon(data []byte) *Icon {
	return &Icon{
		Data:  data,
		Color: theme.DefaultTheme.Text,
	}
}

func (i *Icon) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{Size: gtx.Constraints.Constrain(image.Point{X: 24, Y: 24})}
}

