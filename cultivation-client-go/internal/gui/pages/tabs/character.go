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
	cultivateBtn       *components.Button
	meditateBtn        *components.Button
	sleepBtn           *components.Button
	breakthroughBtn    *components.Button
	exploreBtn         *components.Button
	gatherBtn          *components.Button
	craftBtn           *components.Button
	createMethodBtn    *components.Button
	feedbackMsg        string
	feedbackTime       time.Time
	feedbackColor      color.RGBA
	lastShownOpResult  *types.OperationResult
	lastOpResultString string
}

func NewCharacterTab() *CharacterTab {
	return &CharacterTab{
		cultivateBtn:    components.NewButton("修炼"),
		meditateBtn:     components.NewButton("打坐"),
		sleepBtn:        components.NewButton("休息"),
		breakthroughBtn: components.NewButton("突破"),
		exploreBtn:      components.NewButton("探索"),
		gatherBtn:       components.NewButton("采集"),
		craftBtn:        components.NewButton("炼制"),
		createMethodBtn: components.NewButton("自创功法"),
		feedbackMsg:     "",
		feedbackColor:   theme.DefaultTheme.Success,
	}
}

func (t *CharacterTab) Layout(gtx layout.Context) layout.Dimensions {
	char := store.GetGameStore().GetCharacter()

	t.handleActions(gtx)

	if t.feedbackMsg != "" && time.Since(t.feedbackTime) > 3*time.Second {
		t.feedbackMsg = ""
	}

	// Check for new operation result from server
	if opResult := store.GetGameStore().GetLastOperationResult(); opResult != nil && opResult != t.lastShownOpResult {
		t.lastShownOpResult = opResult
		detailStr := formatOpResult(opResult)
		if detailStr != "" {
			t.showFeedback(opResult.Message+" ("+detailStr+")", feedbackColor(opResult.Success))
		} else {
			t.showFeedback(opResult.Message, feedbackColor(opResult.Success))
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
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
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutActionButtons(gtx)
		}),
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

func (t *CharacterTab) layoutCharacterCard(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return components.DrawCardBackground(gtx, theme.DefaultTheme.Surface)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// 头部信息
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutHeader(gtx, char)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return components.DrawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 修炼进度
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutProgressSection(gtx, char)
						}),
						// Qi/Spiritual 进度条
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return t.layoutBar(gtx, "生命值", float32(char.Qi), float32(char.MaxQi),
											color.RGBA{R: 220, G: 60, B: 60, A: 255})
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return t.layoutBar(gtx, "灵力值", float32(char.SpiritualPower), float32(char.MaxSpiritualPower),
												color.RGBA{R: 60, G: 120, B: 220, A: 255})
										})
									}),
								)
							})
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return components.DrawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 属性区域（两列）
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutStatsGrid(gtx, char)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return components.DrawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 生活技能
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutLifeSkills(gtx, char)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return components.DrawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 灵性/状态
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutSpiritualAttrs(gtx, char)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(12),
								Bottom: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return components.DrawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 灵石
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutSpiritStones(gtx, char)
						}),
					)
				})
			}),
		)
	})
}

func (t *CharacterTab) layoutLifeSkills(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sectionLabel := components.NewLabel("生活技能")
			sectionLabel.Color = theme.DefaultTheme.Primary
			sectionLabel.Size = 14
			return sectionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return t.layoutStatItem(gtx, "煎药", char.AlchemyLevel, color.RGBA{R: 200, G: 150, B: 100, A: 255})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return t.layoutStatItem(gtx, "炼器", char.ArtificingLevel, color.RGBA{R: 150, G: 150, B: 200, A: 255})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return t.layoutStatItem(gtx, "阵法", char.FormationLevel, color.RGBA{R: 100, G: 180, B: 255, A: 255})
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return t.layoutStatItem(gtx, "控火", char.FireControl, color.RGBA{R: 255, G: 120, B: 50, A: 255})
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis: layout.Horizontal,
							}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return t.layoutStatItem(gtx, "草药", char.HerbKnowledge, color.RGBA{R: 100, G: 200, B: 100, A: 255})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return t.layoutStatItem(gtx, "采矿", char.MiningSkill, color.RGBA{R: 180, G: 180, B: 100, A: 255})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return t.layoutStatItem(gtx, "符篆", char.TalismanSkill, color.RGBA{R: 200, G: 100, B: 200, A: 255})
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return t.layoutStatItem(gtx, "御兽", char.BeastTaming, color.RGBA{R: 150, G: 200, B: 150, A: 255})
								}),
							)
						})
					}),
				)
			})
		}),
	)
}

func (t *CharacterTab) layoutSpiritualAttrs(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sectionLabel := components.NewLabel("灵性 / 状态")
			sectionLabel.Color = theme.DefaultTheme.Primary
			sectionLabel.Size = 14
			return sectionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "道心", char.DaoHeart, color.RGBA{R: 255, G: 200, B: 150, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "颂悟", char.Enlightenment, color.RGBA{R: 200, G: 180, B: 255, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "灵根纯度", char.RootPurity, color.RGBA{R: 100, G: 200, B: 255, A: 255})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "声望", char.Reputation, color.RGBA{R: 255, G: 200, B: 80, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "宗门贡献", char.SectContribution, color.RGBA{R: 150, G: 200, B: 100, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "毒性", char.PoisonLevel, color.RGBA{R: 100, G: 200, B: 80, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "咒症", char.CurseLevel, color.RGBA{R: 200, G: 80, B: 200, A: 255})
					}),
				)
			})
		}),
	)
}

func (t *CharacterTab) layoutSpiritStones(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sectionLabel := components.NewLabel("灵石")
			sectionLabel.Color = theme.DefaultTheme.Primary
			sectionLabel.Size = 14
			return sectionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatText(gtx, "下品灵石", fmt.Sprintf("%d", char.LowGradeStones), color.RGBA{R: 180, G: 180, B: 180, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatText(gtx, "中品灵石", fmt.Sprintf("%d", char.MediumGradeStones), color.RGBA{R: 100, G: 200, B: 100, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatText(gtx, "上品灵石", fmt.Sprintf("%d", char.HighGradeStones), color.RGBA{R: 100, G: 150, B: 255, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatText(gtx, "级品灵石", fmt.Sprintf("%d", char.PremiumGradeStones), color.RGBA{R: 255, G: 150, B: 100, A: 255})
					}),
				)
			})
		}),
	)
}


// 头部：头像 + 名称/境界/状态/位置
func (t *CharacterTab) layoutHeader(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return components.DrawAvatar(gtx, theme.DefaultTheme.Primary)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						nameLabel := components.NewLabel(char.Name)
						nameLabel.Color = theme.DefaultTheme.Text
						nameLabel.Size = 20
						return nameLabel.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							realmLabel := components.NewLabel(fmt.Sprintf("境界: %s", realmDisplayName(char.CultivationRealm)))
							realmLabel.Color = theme.DefaultTheme.Secondary
							realmLabel.Size = 14
							return realmLabel.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							statusStr := statusDisplayName(char.Status)
							statusLabel := components.NewLabel(fmt.Sprintf("状态: %s", statusStr))
							statusLabel.Color = theme.DefaultTheme.TextSecondary
							statusLabel.Size = 14
							return statusLabel.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							posLabel := components.NewLabel(fmt.Sprintf("位置: %s (%.1f, %.1f)", char.RegionID, char.PositionX, char.PositionY))
							posLabel.Color = theme.DefaultTheme.TextSecondary
							posLabel.Size = 14
							return posLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 修炼进度
func (t *CharacterTab) layoutProgressSection(gtx layout.Context, char *types.Character) layout.Dimensions {
	return t.layoutBar(gtx, "修炼进度", float32(char.CultivationProgress), 100,
		color.RGBA{R: 255, G: 200, B: 80, A: 255})
}

// 属性网格（两列：左列基础属性，右列战斗/心境/业力）
func (t *CharacterTab) layoutStatsGrid(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// 左列：修炼资质
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						sectionLabel := components.NewLabel("修炼资质")
						sectionLabel.Color = theme.DefaultTheme.Primary
						sectionLabel.Size = 14
						return sectionLabel.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return t.layoutAttrRow(gtx, char)
						})
					}),
				)
			})
		}),
		// 右列：战斗属性 + 心境 + 业力
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						sectionLabel := components.NewLabel("战斗属性")
						sectionLabel.Color = theme.DefaultTheme.Primary
						sectionLabel.Size = 14
						return sectionLabel.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return t.layoutCombatAttrs(gtx, char)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return t.layoutMentalAndLifespan(gtx, char)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return t.layoutKarma(gtx, char)
						})
					}),
				)
			})
		}),
	)
}

// 修炼资质行：悟性、根骨、机缘、神识
func (t *CharacterTab) layoutAttrRow(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutStatItem(gtx, "悟性", char.Comprehension, color.RGBA{R: 180, G: 140, B: 255, A: 255})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutStatItem(gtx, "根骨", char.Constitution, color.RGBA{R: 100, G: 200, B: 100, A: 255})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutStatItem(gtx, "机缘", char.Luck, color.RGBA{R: 255, G: 180, B: 100, A: 255})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutStatItemFloat(gtx, "神识", char.DivineSense, color.RGBA{R: 100, G: 200, B: 255, A: 255})
		}),
	)
}

// 战斗属性：攻击、防御、速度
func (t *CharacterTab) layoutCombatAttrs(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return t.layoutStatItem(gtx, "攻击", char.Attack, color.RGBA{R: 255, G: 100, B: 100, A: 255})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return t.layoutStatItem(gtx, "防御", char.Defense, color.RGBA{R: 100, G: 200, B: 100, A: 255})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return t.layoutStatItem(gtx, "速度", char.Speed, color.RGBA{R: 100, G: 150, B: 255, A: 255})
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "暴击率", char.CritRate, color.RGBA{R: 255, G: 80, B: 80, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "暴击伤害", char.CritDamage, color.RGBA{R: 200, G: 60, B: 60, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "闪避率", char.DodgeRate, color.RGBA{R: 80, G: 200, B: 255, A: 255})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "命中率", char.HitRate, color.RGBA{R: 100, G: 255, B: 100, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "穿透", char.Penetration, color.RGBA{R: 200, G: 100, B: 255, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItemFloat(gtx, "减伤", char.DamageReduction, color.RGBA{R: 80, G: 200, B: 80, A: 255})
					}),
				)
			})
		}),
	)
}

// 心境 + 寿元
func (t *CharacterTab) layoutMentalAndLifespan(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sectionLabel := components.NewLabel("心境 / 寿元")
			sectionLabel.Color = theme.DefaultTheme.Primary
			sectionLabel.Size = 14
			return sectionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "心神", char.MentalStability, color.RGBA{R: 180, G: 100, B: 255, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lifespanLabel := fmt.Sprintf("%d/%d", char.RemainingLifespan, char.MaxLifespan)
						return t.layoutStatText(gtx, "寿元", lifespanLabel, color.RGBA{R: 255, G: 200, B: 100, A: 255})
					}),
				)
			})
		}),
	)
}

// 业力
func (t *CharacterTab) layoutKarma(gtx layout.Context, char *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sectionLabel := components.NewLabel("业力")
			sectionLabel.Color = theme.DefaultTheme.Primary
			sectionLabel.Size = 14
			return sectionLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "业力值", char.KarmaValue, color.RGBA{R: 255, G: 100, B: 100, A: 255})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return t.layoutStatItem(gtx, "功德", char.Merit, color.RGBA{R: 255, G: 200, B: 80, A: 255})
					}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return t.layoutStatItem(gtx, "业债", char.KarmicDebt, color.RGBA{R: 150, G: 80, B: 80, A: 255})
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							hmLabel := char.HeavenlyMark
							if hmLabel == "" {
								hmLabel = "无"
							}
							return t.layoutStatText(gtx, "天印", hmLabel, color.RGBA{R: 100, G: 200, B: 255, A: 255})
						}),
				)
			})
		}),
	)
}

// 通用单个属性项（int 值）
func (t *CharacterTab) layoutStatItem(gtx layout.Context, name string, value int, dotColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel(name)
			lbl.Color = theme.DefaultTheme.TextSecondary
			lbl.Size = 12
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(unit.Dp(8))
						return components.DrawRect(gtx, dotColor, image.Point{X: size, Y: size})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							valLabel := components.NewLabel(fmt.Sprintf("%d", value))
							valLabel.Color = theme.DefaultTheme.Text
							valLabel.Size = 16
							return valLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 通用单个属性项（float 值）
func (t *CharacterTab) layoutStatItemFloat(gtx layout.Context, name string, value float64, dotColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel(name)
			lbl.Color = theme.DefaultTheme.TextSecondary
			lbl.Size = 12
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(unit.Dp(8))
						return components.DrawRect(gtx, dotColor, image.Point{X: size, Y: size})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							valLabel := components.NewLabel(fmt.Sprintf("%.0f", value))
							valLabel.Color = theme.DefaultTheme.Text
							valLabel.Size = 16
							return valLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 通用单个属性文本项（字符串值）
func (t *CharacterTab) layoutStatText(gtx layout.Context, name string, text string, dotColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel(name)
			lbl.Color = theme.DefaultTheme.TextSecondary
			lbl.Size = 12
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(unit.Dp(8))
						return components.DrawRect(gtx, dotColor, image.Point{X: size, Y: size})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							valLabel := components.NewLabel(text)
							valLabel.Color = theme.DefaultTheme.Text
							valLabel.Size = 16
							return valLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 进度条
func (t *CharacterTab) layoutBar(gtx layout.Context, label string, value, max float32, barColor color.RGBA) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
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

// 操作按钮
func (t *CharacterTab) layoutActionButtons(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 第一行：修炼/打坐/休息/突破
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Top:    unit.Dp(8),
				Bottom: unit.Dp(4),
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
		}),
		// 第二行：探索/采集/炼制/自创功法
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Top:    unit.Dp(4),
				Bottom: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Spacing:   layout.SpaceEvenly,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							t.exploreBtn.Color = theme.DefaultTheme.Primary
							return t.exploreBtn.Layout(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							t.gatherBtn.Color = theme.DefaultTheme.Success
							return t.gatherBtn.Layout(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							t.craftBtn.Color = theme.DefaultTheme.Secondary
							return t.craftBtn.Layout(gtx)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							t.createMethodBtn.Color = theme.DefaultTheme.Warning
							return t.createMethodBtn.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

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

	if t.exploreBtn.Clicked(gtx) {
		if err := ws.Explore(); err != nil {
			t.showFeedback("探索失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始探索...", theme.DefaultTheme.Success)
		}
	}

	if t.gatherBtn.Clicked(gtx) {
		if err := ws.Gather("herb", 1); err != nil {
			t.showFeedback("采集失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始采集草药...", theme.DefaultTheme.Success)
		}
	}

	if t.craftBtn.Clicked(gtx) {
		if err := ws.Craft(""); err != nil {
			t.showFeedback("炼制失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("开始炼制...", theme.DefaultTheme.Success)
		}
	}

	if t.createMethodBtn.Clicked(gtx) {
		if err := ws.CreateMethod(); err != nil {
			t.showFeedback("自创功法失败: "+err.Error(), theme.DefaultTheme.Error)
		} else {
			t.showFeedback("尝试自创功法...", theme.DefaultTheme.Success)
		}
	}
}

func (t *CharacterTab) showFeedback(msg string, c color.RGBA) {
	t.feedbackMsg = msg
	t.feedbackTime = time.Now()
	t.feedbackColor = c
}

// formatOpResult 格式化操作结果的 Effects 为可读字符串
func formatOpResult(result *types.OperationResult) string {
	if result == nil || len(result.Effects) == 0 {
		return ""
	}
	parts := []string{}
	if qi, ok := result.Effects["qi_cost"].(float64); ok && qi > 0 {
		parts = append(parts, fmt.Sprintf("消耗灵力 %.0f", qi))
	}
	if gain, ok := result.Effects["cultivation_gain"].(float64); ok && gain > 0 {
		parts = append(parts, fmt.Sprintf("修为 +%.2f", gain))
	}
	if qr, ok := result.Effects["qi_recovery"].(float64); ok && qr > 0 {
		parts = append(parts, fmt.Sprintf("灵力 +%.0f", qr))
	}
	if sr, ok := result.Effects["spiritual_recovery"].(float64); ok && sr > 0 {
		parts = append(parts, fmt.Sprintf("神识 +%.0f", sr))
	}
	if dmg, ok := result.Effects["damage"].(float64); ok && dmg > 0 {
		parts = append(parts, fmt.Sprintf("造成伤害 %.0f", dmg))
	}
	if dmg, ok := result.Effects["damage_dealt"].(float64); ok && dmg > 0 {
		parts = append(parts, fmt.Sprintf("造成伤害 %.0f", dmg))
	}
	if resName, ok := result.Effects["resource"].(string); ok {
		if qty, ok := result.Effects["quantity"].(float64); ok {
			parts = append(parts, fmt.Sprintf("%s x%.0f", resName, qty))
		} else {
			parts = append(parts, resName)
		}
	}
	if discoveries, ok := result.Effects["discoveries"].([]interface{}); ok {
		for _, d := range discoveries {
			if s, ok := d.(string); ok {
				parts = append(parts, s)
			}
		}
	}
	if price, ok := result.Effects["price"].(float64); ok && price > 0 {
		parts = append(parts, fmt.Sprintf("灵石 %d", int(price)))
	}
	if sectName, ok := result.Effects["sect_name"].(string); ok {
		parts = append(parts, fmt.Sprintf("宗门: %s", sectName))
	}
	if methodQ, ok := result.Effects["method_quality"].(string); ok {
		parts = append(parts, fmt.Sprintf("功法品质: %s", methodQ))
	}
	if newRealm, ok := result.Effects["new_realm"].(string); ok {
		parts = append(parts, fmt.Sprintf("晋升: %s", realmDisplayName(newRealm)))
	}
	if cost, ok := result.Effects["cost"].(float64); ok && cost > 0 {
		parts = append(parts, fmt.Sprintf("消耗 %d", int(cost)))
	}
	joined := ""
	for i, p := range parts {
		if i > 0 {
			joined += " | "
		}
		joined += p
	}
	return joined
}

// feedbackColor 根据操作成功/失败返回对应颜色
func feedbackColor(success bool) color.RGBA {
	if success {
		return theme.DefaultTheme.Success
	}
	return theme.DefaultTheme.Error
}

// realmDisplayName 境界代码转中文名
func realmDisplayName(realm string) string {
	switch realm {
	case "mortal":
		return "凡人"
	case "qi_condensation":
		return "练气期"
	case "foundation":
		return "筑基期"
	case "golden_core":
		return "金丹期"
	case "nascent_soul":
		return "元婴期"
	case "soul_transformation":
		return "化神期"
	case "void_refinement":
		return "炼虚期"
	case "integration":
		return "合体期"
	case "mahayana":
		return "大乘期"
	case "tribulation":
		return "渡劫期"
	default:
		if realm == "" {
			return "凡人"
		}
		return realm
	}
}

// statusDisplayName 状态代码转中文名
func statusDisplayName(status string) string {
	switch status {
	case "normal":
		return "正常"
	case "cultivating":
		return "修炼中"
	case "combat":
		return "战斗中"
	case "resting":
		return "休息中"
	case "dead":
		return "已死亡"
	case "exploring":
		return "探索中"
	case "crafting":
		return "炼制中"
	case "meditating":
		return "打坐中"
	default:
		return status
	}
}
