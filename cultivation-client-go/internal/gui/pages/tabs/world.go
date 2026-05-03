package tabs

import (
	"fmt"
	"image/color"
	"time"

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

type WorldTab struct {
	exploreBtn    widget.Clickable
	moveBtn       widget.Clickable
	gatherBtn     widget.Clickable
	eventList     widget.List
	feedbackMsg   string
	feedbackTime  time.Time
	feedbackColor color.RGBA
	showMoveDialog bool
	selectedRegion int
}

func NewWorldTab() *WorldTab {
	return &WorldTab{
		eventList: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		feedbackColor: theme.DefaultTheme.Success,
		selectedRegion: -1,
	}
}

func (t *WorldTab) Layout(gtx layout.Context) layout.Dimensions {
	world := store.GetGameStore().GetWorld()

	t.handleActions(gtx)

	if t.feedbackMsg != "" && time.Since(t.feedbackTime) > 3*time.Second {
		t.feedbackMsg = ""
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("世界")
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
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if world == nil {
				return layout.Inset{
					Left:  unit.Dp(16),
					Right: unit.Dp(16),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					loading := components.NewLabel("加载中...")
					loading.Color = theme.DefaultTheme.TextSecondary
					return loading.Layout(gtx)
				})
			}
			return t.layoutAnnouncements(gtx, world)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if world == nil {
				return layout.Dimensions{}
			}
			return t.layoutMapCard(gtx, world)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if world == nil {
				return layout.Dimensions{}
			}
			return t.layoutEventList(gtx, world)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutActionButtons(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if t.feedbackMsg == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{
				Top:    unit.Dp(8),
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

func (t *WorldTab) layoutAnnouncements(gtx layout.Context, world *types.WorldState) layout.Dimensions {
	if len(world.Announcements) == 0 {
		return layout.Dimensions{}
	}

	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				announcement := world.Announcements[0]
				bgColor := t.getAnnouncementColor(announcement.Priority)
				
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Constraints.Max
						size.Y = gtx.Dp(unit.Dp(40))
						return drawRectWithRadius(gtx, bgColor, size, 8)
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Left:   unit.Dp(12),
							Right:  unit.Dp(12),
							Top:    unit.Dp(8),
							Bottom: unit.Dp(8),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									icon := components.NewLabel("📢 ")
									icon.Size = 16
									return icon.Layout(gtx)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									content := components.NewLabel(announcement.Content)
									content.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
									content.Size = 14
									return content.Layout(gtx)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

func (t *WorldTab) getAnnouncementColor(priority int) color.RGBA {
	switch priority {
	case 3:
		return color.RGBA{R: 255, G: 100, B: 100, A: 255}
	case 2:
		return color.RGBA{R: 255, G: 180, B: 80, A: 255}
	case 1:
		return color.RGBA{R: 255, G: 220, B: 100, A: 255}
	default:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255}
	}
}

func (t *WorldTab) layoutMapCard(gtx layout.Context, world *types.WorldState) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawCardBackground(gtx, theme.DefaultTheme.Surface)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutMapHeader(gtx, world)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutMapInfo(gtx, world)
						}),
					)
				})
			}),
		)
	})
}

func (t *WorldTab) layoutMapHeader(gtx layout.Context, world *types.WorldState) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			mapIcon := components.NewLabel("🗺️ ")
			mapIcon.Size = 24
			return mapIcon.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				mapName := components.NewLabel(world.CurrentMap)
				mapName.Color = theme.DefaultTheme.Text
				mapName.Size = 20
				return mapName.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					onlineIcon := components.NewLabel("👥 ")
					onlineIcon.Size = 14
					onlineIcon.Color = theme.DefaultTheme.Success
					return onlineIcon.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					onlineLabel := components.NewLabel(fmt.Sprintf("%d 人在线", world.PlayersOnline))
					onlineLabel.Color = theme.DefaultTheme.Success
					onlineLabel.Size = 14
					return onlineLabel.Layout(gtx)
				}),
			)
		}),
	)
}

func (t *WorldTab) layoutMapInfo(gtx layout.Context, world *types.WorldState) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutInfoItem(gtx, "坐标", "X: 128, Y: 256", "📍")
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutInfoItem(gtx, "区域", "安全区", "🛡️")
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutInfoItem(gtx, "天气", "晴朗", "☀️")
		}),
	)
}

func (t *WorldTab) layoutInfoItem(gtx layout.Context, label, value, icon string) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			iconLabel := components.NewLabel(icon + " " + label)
			iconLabel.Color = theme.DefaultTheme.TextSecondary
			iconLabel.Size = 12
			return iconLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				valueLabel := components.NewLabel(value)
				valueLabel.Color = theme.DefaultTheme.Text
				valueLabel.Size = 14
				return valueLabel.Layout(gtx)
			})
		}),
	)
}

func (t *WorldTab) layoutEventList(gtx layout.Context, world *types.WorldState) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRectWithRadius(gtx, theme.DefaultTheme.Surface, gtx.Constraints.Max, 12)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									eventTitle := components.NewLabel("⚡ 世界事件")
									eventTitle.Color = theme.DefaultTheme.Secondary
									eventTitle.Size = 16
									return eventTitle.Layout(gtx)
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									countLabel := components.NewLabel(fmt.Sprintf("(%d)", len(world.Events)))
									countLabel.Color = theme.DefaultTheme.TextSecondary
									countLabel.Size = 12
									return countLabel.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if len(world.Events) == 0 {
								emptyLabel := components.NewLabel("暂无世界事件")
								emptyLabel.Color = theme.DefaultTheme.TextSecondary
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return emptyLabel.Layout(gtx)
								})
							}
							return material.List(th, &t.eventList).Layout(gtx, len(world.Events), func(gtx layout.Context, index int) layout.Dimensions {
								return t.layoutEventItem(gtx, world.Events[index])
							})
						}),
					)
				})
			}),
		)
	})
}

func (t *WorldTab) layoutEventItem(gtx layout.Context, event types.WorldEvent) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(6),
		Bottom: unit.Dp(6),
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
						eventIcon := components.NewLabel(t.getEventIcon(event.Type))
						eventIcon.Size = 18
						return eventIcon.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							eventType := components.NewLabel(t.getEventTypeName(event.Type))
							eventType.Color = t.getEventColor(event.Type)
							eventType.Size = 14
							return eventType.Layout(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						timeStr := event.StartTime.Format("15:04")
						timeLabel := components.NewLabel(timeStr)
						timeLabel.Color = theme.DefaultTheme.TextSecondary
						timeLabel.Size = 12
						return timeLabel.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:  unit.Dp(4),
					Left: unit.Dp(26),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					descLabel := components.NewLabel(event.Description)
					descLabel.Color = theme.DefaultTheme.Text
					descLabel.Size = 13
					return descLabel.Layout(gtx)
				})
			}),
		)
	})
}

func (t *WorldTab) getEventIcon(eventType string) string {
	switch eventType {
	case "boss_spawn":
		return "👹"
	case "resource_spawn":
		return "💎"
	case "activity":
		return "🎉"
	case "weather":
		return "🌪️"
	case "pvp":
		return "⚔️"
	default:
		return "📌"
	}
}

func (t *WorldTab) getEventTypeName(eventType string) string {
	switch eventType {
	case "boss_spawn":
		return "BOSS刷新"
	case "resource_spawn":
		return "资源刷新"
	case "activity":
		return "活动"
	case "weather":
		return "天气变化"
	case "pvp":
		return "PVP事件"
	default:
		return "其他"
	}
}

func (t *WorldTab) getEventColor(eventType string) color.RGBA {
	switch eventType {
	case "boss_spawn":
		return theme.DefaultTheme.Error
	case "resource_spawn":
		return theme.DefaultTheme.Success
	case "activity":
		return theme.DefaultTheme.Warning
	case "weather":
		return theme.DefaultTheme.Primary
	case "pvp":
		return color.RGBA{R: 255, G: 100, B: 100, A: 255}
	default:
		return theme.DefaultTheme.TextSecondary
	}
}

func (t *WorldTab) layoutActionButtons(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Spacing:   layout.SpaceEvenly,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutButton(gtx, &t.exploreBtn, "🔍 探索", theme.DefaultTheme.Primary)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutButton(gtx, &t.moveBtn, "🚶 移动", theme.DefaultTheme.Secondary)
				})
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.layoutButton(gtx, &t.gatherBtn, "⛏️ 采集", theme.DefaultTheme.Success)
				})
			}),
		)
	})
}

func (t *WorldTab) layoutButton(gtx layout.Context, clickable *widget.Clickable, text string, btnColor color.RGBA) layout.Dimensions {
	btn := material.Button(th, clickable, text)
	btn.Background = toNRGBA(btnColor)
	btn.Inset = layout.UniformInset(unit.Dp(12))
	return btn.Layout(gtx)
}

func (t *WorldTab) handleActions(gtx layout.Context) {
	ws := network.GetWebSocketClient()

	if t.exploreBtn.Clicked(gtx) {
		if err := ws.Explore(); err != nil {
			t.showFeedback("探索失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始探索周围区域...", theme.DefaultTheme.Success)
		}
	}

	if t.moveBtn.Clicked(gtx) {
		t.showMoveDialog = true
		regions := []string{"青云镇", "东荒", "北域", "南疆", "西域"}
		if t.selectedRegion >= 0 && t.selectedRegion < len(regions) {
			region := regions[t.selectedRegion]
			if err := ws.Move(region, 0, 0); err != nil {
				t.showFeedback("移动失败: "+err.Error(), theme.DefaultTheme.Error)
			} else {
				t.showFeedback("正在前往 "+region+"...", theme.DefaultTheme.Success)
			}
			t.selectedRegion = -1
			t.showMoveDialog = false
		} else {
			t.showFeedback("请选择目标区域", theme.DefaultTheme.Warning)
			t.selectedRegion = 0
		}
	}

	if t.gatherBtn.Clicked(gtx) {
		if err := ws.Gather("", 1); err != nil {
			t.showFeedback("采集失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始采集资源...", theme.DefaultTheme.Success)
		}
	}
}

func (t *WorldTab) showFeedback(msg string, c color.RGBA) {
	t.feedbackMsg = msg
	t.feedbackTime = time.Now()
	t.feedbackColor = c
}
