package tabs

import (
	"fmt"
	"image"
	"image/color"
	"sort"
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

type ChatChannel int

const (
	ChannelWorld ChatChannel = iota
	ChannelPrivate
)

type SocialTab struct {
	// 频道切换
	worldChannelBtn   widget.Clickable
	privateChannelBtn widget.Clickable
	currentChannel    ChatChannel

	// 好友列表
	friendList     widget.List
	selectedFriend string
	friendItems    map[string]*friendItem

	// 消息列表
	messageList widget.List

	// 聊天输入
	messageInput widget.Editor
	sendBtn      *components.Button

	// 好友操作按钮
	addFriendBtn    *components.Button
	deleteFriendBtn *components.Button

	// 添加好友对话框
	showAddFriendDialog bool
	addFriendInput      widget.Editor
	addFriendConfirmBtn *components.Button
	addFriendCancelBtn  *components.Button

	// 反馈消息
	feedbackMsg   string
	feedbackTime  time.Time
	feedbackColor color.RGBA
}

type friendItem struct {
	clickable widget.Clickable
	friend    types.Friend
}

func NewSocialTab() *SocialTab {
	addFriendBtn := components.NewButton("+")
	addFriendBtn.Color = theme.DefaultTheme.Success
	addFriendBtn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}

	deleteFriendBtn := components.NewButton("-")
	deleteFriendBtn.Color = theme.DefaultTheme.Error
	deleteFriendBtn.Inset = layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}

	return &SocialTab{
		currentChannel: ChannelWorld,
		friendItems:    make(map[string]*friendItem),
		friendList: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		messageList: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		messageInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		sendBtn:            btn("发送", theme.DefaultTheme.Primary),
		addFriendBtn:       addFriendBtn,
		deleteFriendBtn:    deleteFriendBtn,
		addFriendConfirmBtn: btn("添加", theme.DefaultTheme.Primary),
		addFriendCancelBtn:  btn("取消", theme.DefaultTheme.Border),
		feedbackColor:      theme.DefaultTheme.Success,
	}
}

func (t *SocialTab) Layout(gtx layout.Context) layout.Dimensions {
	social := store.GetGameStore().GetSocial()
	character := store.GetGameStore().GetCharacter()

	// 处理事件
	t.handleEvents(gtx)

	// 检查反馈消息是否过期
	if t.feedbackMsg != "" && time.Since(t.feedbackTime) > 3*time.Second {
		t.feedbackMsg = ""
	}

	content := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 标题
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("社交")
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
		// 主内容区域
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if social == nil {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					loading := components.NewLabel("加载中...")
					loading.Color = theme.DefaultTheme.TextSecondary
					return loading.Layout(gtx)
				})
			}
			return t.layoutMainContent(gtx, social, character)
		}),
		// 反馈消息
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

	// 添加好友对话框
	if t.showAddFriendDialog {
		return t.layoutAddFriendDialog(gtx, content)
	}

	return content
}

// 添加好友对话框布局
func (t *SocialTab) layoutAddFriendDialog(gtx layout.Context, content layout.Dimensions) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return content
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// 半透明背景
			return drawRect(gtx, color.RGBA{R: 0, G: 0, B: 0, A: 128}, gtx.Constraints.Max)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// 对话框
			dialogWidth := int(float32(gtx.Constraints.Max.X) * 0.6)
			if dialogWidth > 400 {
				dialogWidth = 400
			}
			dialogHeight := 200

			offsetX := (gtx.Constraints.Max.X - dialogWidth) / 2
			offsetY := (gtx.Constraints.Max.Y - dialogHeight) / 2

			return layout.Inset{
				Top:    unit.Dp(offsetY),
				Bottom: unit.Dp(offsetY),
				Left:   unit.Dp(offsetX),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = dialogWidth
				gtx.Constraints.Max.Y = dialogHeight

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
									title := components.NewLabel("添加好友")
									title.Color = theme.DefaultTheme.Primary
									title.Size = 18
									return title.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										border := widget.Border{
											Color:        toNRGBA(theme.DefaultTheme.Border),
											CornerRadius: unit.Dp(4),
											Width:        unit.Dp(1),
										}
										return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												editor := material.Editor(th, &t.addFriendInput, "输入好友名称")
												editor.Color = toNRGBA(theme.DefaultTheme.Text)
												return editor.Layout(gtx)
											})
										})
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis:      layout.Horizontal,
											Alignment: layout.Middle,
											Spacing:   layout.SpaceEvenly,
										}.Layout(gtx,
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return t.addFriendCancelBtn.Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return t.addFriendConfirmBtn.Layout(gtx)
												})
											}),
										)
									})
								}),
							)
						})
					}),
				)
			})
		}),
	)
}

// 主内容区域布局
func (t *SocialTab) layoutMainContent(gtx layout.Context, social *types.SocialInfo, character *types.Character) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// 左侧好友列表（30%宽度）
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * 0.3)
			return t.layoutFriendList(gtx, social)
		}),
		// 右侧聊天区域（70%宽度）
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return t.layoutChatArea(gtx, social, character)
		}),
	)
}

// 好友列表布局
func (t *SocialTab) layoutFriendList(gtx layout.Context, social *types.SocialInfo) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(8),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			// 背景
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRectWithRadius(gtx, theme.DefaultTheme.Surface, gtx.Constraints.Max, 8)
			}),
			// 内容
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// 标题和操作按钮
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutFriendListHeader(gtx, len(social.Friends))
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(8),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 好友列表
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if len(social.Friends) == 0 {
								emptyLabel := components.NewLabel("暂无好友")
								emptyLabel.Color = theme.DefaultTheme.TextSecondary
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return emptyLabel.Layout(gtx)
								})
							}
							return t.layoutFriends(gtx, social.Friends)
						}),
					)
				})
			}),
		)
	})
}

// 好友列表头部
func (t *SocialTab) layoutFriendListHeader(gtx layout.Context, friendCount int) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel(fmt.Sprintf("好友 (%d)", friendCount))
			title.Color = theme.DefaultTheme.Text
			title.Size = 16
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.addFriendBtn.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return t.deleteFriendBtn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

// 图标按钮

// 好友列表项布局
func (t *SocialTab) layoutFriends(gtx layout.Context, friends []types.Friend) layout.Dimensions {
	// 按在线状态排序：在线在前
	sortedFriends := make([]types.Friend, len(friends))
	copy(sortedFriends, friends)
	sort.Slice(sortedFriends, func(i, j int) bool {
		if sortedFriends[i].Online != sortedFriends[j].Online {
			return sortedFriends[i].Online
		}
		return sortedFriends[i].Name < sortedFriends[j].Name
	})

	return material.List(th, &t.friendList).Layout(gtx, len(sortedFriends), func(gtx layout.Context, index int) layout.Dimensions {
		friend := sortedFriends[index]
		return t.layoutFriendItem(gtx, friend)
	})
}

// 单个好友项布局
func (t *SocialTab) layoutFriendItem(gtx layout.Context, friend types.Friend) layout.Dimensions {
	item, exists := t.friendItems[friend.ID]
	if !exists {
		item = &friendItem{friend: friend}
		t.friendItems[friend.ID] = item
	}
	item.friend = friend

	isSelected := t.selectedFriend == friend.ID

	return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 背景色
		bgColor := theme.DefaultTheme.Surface
		if isSelected {
			bgColor = theme.DefaultTheme.Hover
		}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRect(gtx, bgColor, gtx.Constraints.Max)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						// 在线状态指示器
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							statusColor := theme.DefaultTheme.Success
							if !friend.Online {
								statusColor = theme.DefaultTheme.Disabled
							}
							size := gtx.Dp(unit.Dp(8))
							return drawCircle(gtx, statusColor, image.Point{X: size, Y: size})
						}),
						// 好友名称和等级
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										nameLabel := components.NewLabel(friend.Name)
										nameLabel.Color = theme.DefaultTheme.Text
										nameLabel.Size = 14
										return nameLabel.Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										levelLabel := components.NewLabel(fmt.Sprintf("Lv.%d", friend.Level))
										levelLabel.Color = theme.DefaultTheme.TextSecondary
										levelLabel.Size = 12
										return levelLabel.Layout(gtx)
									}),
								)
							})
						}),
					)
				})
			}),
		)
	})
}

// 聊天区域布局
func (t *SocialTab) layoutChatArea(gtx layout.Context, social *types.SocialInfo, character *types.Character) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(8),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			// 背景
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRectWithRadius(gtx, theme.DefaultTheme.Surface, gtx.Constraints.Max, 8)
			}),
			// 内容
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// 频道切换按钮
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutChannelTabs(gtx)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(8),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 消息列表
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return t.layoutMessageList(gtx, social, character)
						}),
						// 分隔线
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top:    unit.Dp(8),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawDivider(gtx, theme.DefaultTheme.Border)
							})
						}),
						// 输入区域
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return t.layoutInputArea(gtx)
						}),
					)
				})
			}),
		)
	})
}

// 频道切换标签
func (t *SocialTab) layoutChannelTabs(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutChannelButton(gtx, &t.worldChannelBtn, "世界", ChannelWorld)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return t.layoutChannelButton(gtx, &t.privateChannelBtn, "私聊", ChannelPrivate)
			})
		}),
	)
}

// 频道按钮
func (t *SocialTab) layoutChannelButton(gtx layout.Context, clickable *widget.Clickable, text string, channel ChatChannel) layout.Dimensions {
	isActive := t.currentChannel == channel
	btnColor := theme.DefaultTheme.Primary
	if !isActive {
		btnColor = theme.DefaultTheme.Border
	}

	btn := material.Button(th, clickable, text)
	btn.Background = toNRGBA(btnColor)
	btn.Inset = layout.UniformInset(unit.Dp(8))
	return btn.Layout(gtx)
}

// 消息列表布局
func (t *SocialTab) layoutMessageList(gtx layout.Context, social *types.SocialInfo, character *types.Character) layout.Dimensions {
	// 过滤消息
	var messages []types.Message
	if t.currentChannel == ChannelWorld {
		// 世界频道显示所有消息
		messages = social.Messages
	} else {
		// 私聊频道只显示与选中好友的消息
		if t.selectedFriend != "" {
			for _, msg := range social.Messages {
				if msg.SenderID == t.selectedFriend {
					messages = append(messages, msg)
				}
			}
		}
	}

	if len(messages) == 0 {
		emptyText := "暂无消息"
		if t.currentChannel == ChannelPrivate && t.selectedFriend == "" {
			emptyText = "请选择一个好友"
		}
		emptyLabel := components.NewLabel(emptyText)
		emptyLabel.Color = theme.DefaultTheme.TextSecondary
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return emptyLabel.Layout(gtx)
		})
	}

	return material.List(th, &t.messageList).Layout(gtx, len(messages), func(gtx layout.Context, index int) layout.Dimensions {
		msg := messages[index]
		isSelf := character != nil && msg.SenderID == character.ID
		return t.layoutMessageItem(gtx, msg, isSelf)
	})
}

// 单条消息布局
func (t *SocialTab) layoutMessageItem(gtx layout.Context, msg types.Message, isSelf bool) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 根据发送者决定对齐方式
		if isSelf {
			// 自己发送的消息靠右显示
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.layoutMessageBubble(gtx, msg, isSelf)
				}),
			)
		}
		// 他人消息靠左显示
		return layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutMessageBubble(gtx, msg, isSelf)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
		)
	})
}

// 消息气泡布局
func (t *SocialTab) layoutMessageBubble(gtx layout.Context, msg types.Message, isSelf bool) layout.Dimensions {
	bgColor := theme.DefaultTheme.Hover
	if isSelf {
		bgColor = theme.DefaultTheme.Primary
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return drawRectWithRadius(gtx, bgColor, gtx.Constraints.Max, 8)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// 发送者名称和时间
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								senderLabel := components.NewLabel(msg.SenderName)
								senderLabel.Color = theme.DefaultTheme.Text
								if isSelf {
									senderLabel.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
								}
								senderLabel.Size = 12
								return senderLabel.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									timeStr := msg.Timestamp.Format("15:04")
									timeLabel := components.NewLabel(timeStr)
									timeLabel.Color = theme.DefaultTheme.TextSecondary
									if isSelf {
										timeLabel.Color = color.RGBA{R: 200, G: 200, B: 255, A: 255}
									}
									timeLabel.Size = 10
									return timeLabel.Layout(gtx)
								})
							}),
						)
					}),
					// 消息内容
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							contentLabel := components.NewLabel(msg.Content)
							contentLabel.Color = theme.DefaultTheme.Text
							if isSelf {
								contentLabel.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
							}
							contentLabel.Size = 14
							return contentLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// 输入区域布局
func (t *SocialTab) layoutInputArea(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		// 输入框
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			border := widget.Border{
				Color:        toNRGBA(theme.DefaultTheme.Border),
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(1),
			}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, &t.messageInput, t.getInputPlaceholder())
					editor.Color = toNRGBA(theme.DefaultTheme.Text)
					return editor.Layout(gtx)
				})
			})
		}),
		// 发送按钮
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.sendBtn.Layout(gtx)
			})
		}),
	)
}

// 获取输入框占位符文本
func (t *SocialTab) getInputPlaceholder() string {
	if t.currentChannel == ChannelPrivate {
		if t.selectedFriend == "" {
			return "请先选择一个好友"
		}
		return "输入私聊消息..."
	}
	return "输入世界消息..."
}

// 处理事件
func (t *SocialTab) handleEvents(gtx layout.Context) {
	// 频道切换
	if t.worldChannelBtn.Clicked(gtx) {
		t.currentChannel = ChannelWorld
		t.selectedFriend = ""
	}

	if t.privateChannelBtn.Clicked(gtx) {
		t.currentChannel = ChannelPrivate
	}

	// 好友选择
	social := store.GetGameStore().GetSocial()
	if social != nil {
		for _, friend := range social.Friends {
			if item, exists := t.friendItems[friend.ID]; exists {
				if item.clickable.Clicked(gtx) {
					t.selectedFriend = friend.ID
					t.currentChannel = ChannelPrivate
				}
			}
		}
	}

	// 发送消息
	if t.sendBtn.Clicked(gtx) || (t.messageInput.Submit && t.messageInput.Text() != "") {
		t.sendMessage()
	}

	// 添加好友按钮
	if t.addFriendBtn.Clicked(gtx) {
		t.showAddFriendDialog = true
		t.addFriendInput.SetText("")
	}

	// 对话框取消按钮
	if t.addFriendCancelBtn.Clicked(gtx) {
		t.showAddFriendDialog = false
		t.addFriendInput.SetText("")
	}

	// 对话框确认按钮
	if t.addFriendConfirmBtn.Clicked(gtx) {
		friendName := t.addFriendInput.Text()
		if friendName != "" {
			ws := network.GetWebSocketClient()
			if err := ws.SendOperation("add_friend", map[string]interface{}{
				"name": friendName,
			}); err != nil {
				t.showFeedback("添加好友失败: "+err.Error(), theme.DefaultTheme.Error)
			} else {
				t.showFeedback("已发送好友请求: "+friendName, theme.DefaultTheme.Success)
			}
		}
		t.showAddFriendDialog = false
		t.addFriendInput.SetText("")
	}

	// 删除好友
	if t.deleteFriendBtn.Clicked(gtx) {
		if t.selectedFriend != "" {
			ws := network.GetWebSocketClient()
			if err := ws.SendOperation("remove_friend", map[string]interface{}{
				"friend_id": t.selectedFriend,
			}); err != nil {
				t.showFeedback("删除好友失败: "+err.Error(), theme.DefaultTheme.Error)
			} else {
				t.showFeedback("已删除好友", theme.DefaultTheme.Success)
				t.selectedFriend = ""
			}
		} else {
			t.showFeedback("请先选择一个好友", theme.DefaultTheme.Warning)
		}
	}
}

// 发送消息
func (t *SocialTab) sendMessage() {
	content := t.messageInput.Text()
	if content == "" {
		return
	}

	ws := network.GetWebSocketClient()

	var msgType string
	var receiverID string

	if t.currentChannel == ChannelWorld {
		msgType = "world"
		receiverID = ""
	} else {
		if t.selectedFriend == "" {
			t.showFeedback("请先选择一个好友", theme.DefaultTheme.Warning)
			return
		}
		msgType = "private"
		receiverID = t.selectedFriend
	}

	if err := ws.SendMessage(content, msgType, receiverID); err != nil {
		t.showFeedback("发送失败: "+err.Error(), theme.DefaultTheme.Error)
	} else {
		t.messageInput.SetText("")
	}
}

// 显示反馈消息
func (t *SocialTab) showFeedback(msg string, color color.RGBA) {
	t.feedbackMsg = msg
	t.feedbackTime = time.Now()
	t.feedbackColor = color
}

