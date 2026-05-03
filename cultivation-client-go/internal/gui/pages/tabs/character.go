package tabs

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/network"
	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
	"gioui.org/layout"
	"gioui.org/unit"
)

type CharacterTab struct {
	cultivateBtn    *components.Button
	meditateBtn     *components.Button
	sleepBtn        *components.Button
	breakthroughBtn *components.Button
	feedbackMsg     string
	feedbackTime    time.Time
	feedbackColor   color.RGBA
}

func NewCharacterTab() *CharacterTab {
	return &CharacterTab{
		cultivateBtn:    components.NewButton("修炼"),
		meditateBtn:     components.NewButton("打坐"),
		sleepBtn:        components.NewButton("休息"),
		breakthroughBtn: components.NewButton("突破"),
		feedbackMsg:     "",
		feedbackColor:   theme.DefaultTheme.Success,
	}
}

func (t *CharacterTab) Layout(gtx layout.Context) layout.Dimensions {
	char := store.GetGameStore().GetCharacter()

	// 处理按钮点击
	t.handleActions(gtx)

	// 检查反馈消息是否过期（3秒后清除）
	if t.feedbackMsg != "" && time.Since(t.feedbackTime) > 3*time.Second {
		t.feedbackMsg = ""
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 标题
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("角色信息")
			title.Color = theme.DefaultTheme.Primary
			title.Size = 24
			return layout.Inset{
				Top:    unit.Dp(16),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Bottom: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return title.Layout(gtx)
			})
		}),
		// 角色卡片
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if char == nil {
				return layout.Inset{
					Left:  unit.Dp(16),
					Right: unit.Dp(16),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return components.NewLabel("加载中...").Layout(gtx)
				})
			}
			return t.layoutCharacterCard(gtx, char)
		}),
		// 操作按钮
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutActionButtons(gtx)
		}),
		// 反馈消息
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if t.feedbackMsg == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{
				Top:    unit.Dp(12),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Bottom: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := components.NewLabel(t.feedbackMsg)
				lbl.Color = t.feedbackColor
				lbl.Size = 14
				return lbl.Layout(gtx)
			})
		}),
	)
}

// 角色卡片布局
func (t *CharacterTab) layoutCharacterCard(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 卡片背景
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawCardBackground(gtx, theme.DefaultTheme.Surface)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// 头部：头像 + 名称/等级/境界
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutCharacterHeader(gtx, char)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(16),
								Bottom: unit.Dp(16),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 生命值进度条
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutProgressBar(gtx, "生命值", float32(char.Health), float32(char.MaxHealth), color.RGBA{R: 220, G: 60, B: 60, A: 255})
						}),
						// 能量值进度条
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return t.layoutProgressBar(gtx, "能量值", float32(char.Energy), float32(char.MaxEnergy), color.RGBA{R: 60, G: 120, B: 220, A: 255})
							})
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(16),
								Bottom: unit.Dp(16),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 属性网格
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutAttributesGrid(gtx, char)
						}),
					)
				})
			}),
		)
	})
}

// 角色头部布局（头像 + 信息）
func (t *CharacterTab) layoutCharacterHeader(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// 头像区域
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawAvatar(gtx, theme.DefaultTheme.Primary)
		}),
		// 角色信息
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// 角色名称
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						nameLabel := components.NewLabel(char.Name)
						nameLabel.Color = theme.DefaultTheme.Text
						nameLabel.Size = 20
						return nameLabel.Layout(gtx)
					}),
					// 等级
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							levelLabel := components.NewLabel(fmt.Sprintf("等级: %d", char.Level))
							levelLabel.Color = theme.DefaultTheme.TextSecondary
							levelLabel.Size = 14
							return levelLabel.Layout(gtx)
						})
					}),
					// 境界
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							realm := char.CultivationRealm
							if realm == "" {
								realm = "凡人"
							}
							realmLabel := components.NewLabel(fmt.Sprintf("境界: %s", realm))
							realmLabel.Color = theme.DefaultTheme.Secondary
							realmLabel.Size = 14
							return realmLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 进度条布局
func (t *CharacterTab) layoutProgressBar(gtx layout.Context, label string, value, max float32, barColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 标签
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := components.NewLabel(label)
					lbl.Color = theme.DefaultTheme.TextSecondary
					lbl.Size = 12
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					valueText := fmt.Sprintf("%.0f/%.0f", value, max)
					lbl := components.NewLabel(valueText)
					lbl.Color = theme.DefaultTheme.Text
					lbl.Size = 12
					return lbl.Layout(gtx)
				}),
			)
		}),
		// 进度条
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				pb := components.NewProgressBar(value, max)
				pb.SetColor(barColor)
				pb.SetBgColor(theme.DefaultTheme.Border)
				pb.ShowText = false
				pb.Height = unit.Dp(10)
				return pb.Layout(gtx)
			})
		}),
	)
}

// 属性网格布局
func (t *CharacterTab) layoutAttributesGrid(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutAttributeItem(gtx, "攻击", char.Attack, color.RGBA{R: 255, G: 100, B: 100, A: 255})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutAttributeItem(gtx, "防御", char.Defense, color.RGBA{R: 100, G: 200, B: 100, A: 255})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutAttributeItem(gtx, "速度", char.Speed, color.RGBA{R: 100, G: 150, B: 255, A: 255})
		}),
	)
}

// 单个属性项布局
func (t *CharacterTab) layoutAttributeItem(gtx layout.Context, name string, value int, iconColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// 属性名称
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel(name)
			lbl.Color = theme.DefaultTheme.TextSecondary
			lbl.Size = 12
			return lbl.Layout(gtx)
		}),
		// 属性值
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					// 小色块图标
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(unit.Dp(8))
						return drawRect(gtx, iconColor, image.Point{X: size, Y: size})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							valueLabel := components.NewLabel(fmt.Sprintf("%d", value))
							valueLabel.Color = theme.DefaultTheme.Text
							valueLabel.Size = 16
							return valueLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 操作按钮布局
func (t *CharacterTab) layoutActionButtons(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceEvenly,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					t.cultivateBtn.Color = theme.DefaultTheme.Primary
					return t.cultivateBtn.Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					t.meditateBtn.Color = theme.DefaultTheme.Secondary
					return t.meditateBtn.Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					t.sleepBtn.Color = theme.DefaultTheme.Success
					return t.sleepBtn.Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					t.breakthroughBtn.Color = theme.DefaultTheme.Warning
					return t.breakthroughBtn.Layout(gtx)
				})
			}),
		)
	})
}

// 处理操作按钮点击
func (t *CharacterTab) handleActions(gtx layout.Context) {
	ws := network.GetWebSocketClient()

	if t.cultivateBtn.Clicked(gtx) {
		if err := ws.Cultivate(); err != nil {
			t.showFeedback("修炼失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始修炼...", theme.DefaultTheme.Success)
		}
	}

	if t.meditateBtn.Clicked(gtx) {
		if err := ws.Meditate(); err != nil {
			t.showFeedback("打坐失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始打坐...", theme.DefaultTheme.Success)
		}
	}

	if t.sleepBtn.Clicked(gtx) {
		if err := ws.Sleep(); err != nil {
			t.showFeedback("休息失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始休息...", theme.DefaultTheme.Success)
		}
	}

	if t.breakthroughBtn.Clicked(gtx) {
		if err := ws.Breakthrough(); err != nil {
			t.showFeedback("突破失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("尝试突破境界...", theme.DefaultTheme.Success)
		}
	}
}

// 显示反馈消息
func (t *CharacterTab) showFeedback(msg string, color color.RGBA) {
	t.feedbackMsg = msg
	t.feedbackTime = time.Now()
	t.feedbackColor = color
}
