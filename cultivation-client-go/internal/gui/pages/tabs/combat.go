package tabs

import (
	"fmt"
	"image"
	"image/color"

	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/network"
	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CombatTab struct {
	attackBtn   widget.Clickable
	skillBtn    widget.Clickable
	fleeBtn     widget.Clickable
	exploreBtn  widget.Clickable
	list        widget.List
	listAdapter *BattleLogList
}

func NewCombatTab() *CombatTab {
	return &CombatTab{
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		listAdapter: &BattleLogList{},
	}
}

func (t *CombatTab) Layout(gtx layout.Context) layout.Dimensions {
	combat := store.GetGameStore().GetCombat()

	t.handleEvents(gtx)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("战斗")
			title.Color = theme.DefaultTheme.Primary
			title.Size = 24
			return layout.Inset{
				Top:    unit.Dp(16),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Bottom: unit.Dp(24),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if combat == nil || !combat.InCombat {
				return t.layoutNoCombat(gtx)
			}
			return t.layoutInCombat(gtx, combat)
		}),
	)
}

func (t *CombatTab) layoutNoCombat(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(32),
				Bottom: unit.Dp(24),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				noCombat := components.NewLabel("当前不在战斗中")
				noCombat.Color = theme.DefaultTheme.TextSecondary
				noCombat.Size = 18
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return noCombat.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutExploreButton(gtx)
				})
			})
		}),
	)
}

func (t *CombatTab) layoutInCombat(gtx layout.Context, combat *types.CombatState) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutEnemyInfo(gtx, combat)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutTurnInfo(gtx, combat)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutBattleLog(gtx, combat)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutActionButtons(gtx)
		}),
	)
}

func (t *CombatTab) layoutEnemyInfo(gtx layout.Context, combat *types.CombatState) layout.Dimensions {
	if combat.CurrentEnemy == nil {
		return layout.Dimensions{}
	}

	enemy := combat.CurrentEnemy
	healthPercent := float32(enemy.Health) / float32(enemy.MaxHealth)
	if healthPercent < 0 {
		healthPercent = 0
	}
	if healthPercent > 1 {
		healthPercent = 1
	}

	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						nameLabel := components.NewLabel(fmt.Sprintf("敌人: %s", enemy.Name))
						nameLabel.Color = theme.DefaultTheme.Text
						nameLabel.Size = 18
						return nameLabel.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						levelLabel := components.NewLabel(fmt.Sprintf(" (Lv.%d)", enemy.Level))
						levelLabel.Color = theme.DefaultTheme.Warning
						levelLabel.Size = 14
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return levelLabel.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutHealthBar(gtx, enemy.Health, enemy.MaxHealth, healthPercent)
				})
			}),
		)
	})
}

func (t *CombatTab) layoutHealthBar(gtx layout.Context, current, max int, percent float32) layout.Dimensions {
	barHeight := gtx.Dp(unit.Dp(20))

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			size.Y = barHeight
			return drawRectWithRadius(gtx, theme.DefaultTheme.Border, size, 4)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Max
					size.Y = barHeight
					size.X = int(float32(size.X) * percent)
					if size.X > 0 {
						hpColor := theme.DefaultTheme.Success
						if percent < 0.3 {
							hpColor = theme.DefaultTheme.Error
						} else if percent < 0.6 {
							hpColor = theme.DefaultTheme.Warning
						}
						return drawRectWithRadius(gtx, hpColor, size, 4)
					}
					return layout.Dimensions{Size: image.Point{X: 0, Y: barHeight}}
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: barHeight}}
				}),
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			text := fmt.Sprintf("%d / %d", current, max)
			lbl := material.Label(th, 12, text)
			lbl.Color = toNRGBA(theme.DefaultTheme.Text)
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: barHeight}}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return lbl.Layout(gtx)
				}),
			)
		}),
	)
}

func (t *CombatTab) layoutTurnInfo(gtx layout.Context, combat *types.CombatState) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		turnLabel := components.NewLabel(fmt.Sprintf("第 %d 回合", combat.TurnNumber))
		turnLabel.Color = theme.DefaultTheme.Secondary
		turnLabel.Size = 16
		return turnLabel.Layout(gtx)
	})
}

func (t *CombatTab) layoutBattleLog(gtx layout.Context, combat *types.CombatState) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRectWithRadius(gtx, theme.DefaultTheme.Surface, gtx.Constraints.Max, 8)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							logTitle := components.NewLabel("战斗日志")
							logTitle.Color = theme.DefaultTheme.TextSecondary
							logTitle.Size = 14
							return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return logTitle.Layout(gtx)
							})
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if len(combat.BattleLog) == 0 {
								emptyLabel := components.NewLabel("暂无战斗记录")
								emptyLabel.Color = theme.DefaultTheme.TextSecondary
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return emptyLabel.Layout(gtx)
								})
							}
							t.listAdapter.SetLogs(combat.BattleLog)
							return material.List(th, &t.list).Layout(gtx, len(combat.BattleLog), func(gtx layout.Context, index int) layout.Dimensions {
								return t.listAdapter.Layout(gtx, index)
							})
						}),
					)
				})
			}),
		)
	})
}

func (t *CombatTab) layoutActionButtons(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(24),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceEvenly,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutButton(gtx, &t.attackBtn, "攻击", theme.DefaultTheme.Error)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutButton(gtx, &t.skillBtn, "技能", theme.DefaultTheme.Primary)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutButton(gtx, &t.fleeBtn, "逃跑", theme.DefaultTheme.Warning)
			}),
		)
	})
}

func (t *CombatTab) layoutExploreButton(gtx layout.Context) layout.Dimensions {
	return t.layoutButton(gtx, &t.exploreBtn, "寻找敌人", theme.DefaultTheme.Primary)
}

func (t *CombatTab) layoutButton(gtx layout.Context, clickable *widget.Clickable, text string, btnColor color.RGBA) layout.Dimensions {
	btn := material.Button(th, clickable, text)
	btn.Background = toNRGBA(btnColor)
	btn.Inset = layout.UniformInset(unit.Dp(12))
	return btn.Layout(gtx)
}

func (t *CombatTab) handleEvents(gtx layout.Context) {
	ws := network.GetWebSocketClient()

	if t.attackBtn.Clicked(gtx) {
		if err := ws.SendOperation("combat", map[string]interface{}{}); err != nil {
			fmt.Printf("攻击失败: %v\n", err)
		}
	}

	if t.skillBtn.Clicked(gtx) {
		if err := ws.SendOperation("use_skill", map[string]interface{}{}); err != nil {
			fmt.Printf("使用技能失败: %v\n", err)
		}
	}

	if t.fleeBtn.Clicked(gtx) {
		if err := ws.SendOperation("flee", map[string]interface{}{}); err != nil {
			fmt.Printf("逃跑失败: %v\n", err)
		}
	}

	if t.exploreBtn.Clicked(gtx) {
		if err := ws.Explore(); err != nil {
			fmt.Printf("探索失败: %v\n", err)
		}
	}
}

type BattleLogList struct {
	logs []types.CombatLog
}

func (l *BattleLogList) SetLogs(logs []types.CombatLog) {
	l.logs = logs
}

func (l *BattleLogList) Layout(gtx layout.Context, index int) layout.Dimensions {
	if index >= len(l.logs) {
		return layout.Dimensions{}
	}

	log := l.logs[index]

	var logColor color.RGBA
	switch log.Type {
	case "attack":
		logColor = theme.DefaultTheme.Error
	case "heal":
		logColor = theme.DefaultTheme.Success
	default:
		logColor = theme.DefaultTheme.Text
	}

	timeStr := log.Timestamp.Format("15:04:05")
	logText := fmt.Sprintf("[%s] %s", timeStr, log.Message)

	lbl := material.Label(th, 12, logText)
	lbl.Color = toNRGBA(logColor)

	return layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return lbl.Layout(gtx)
	})
}
