package tabs

import (
	"image"

	"cultivation-client/internal/gui/components"
	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SettingsTab struct {
	// 音量滑块
	audioVolume *components.Slider
	musicVolume *components.Slider

	// 开关选项
	showDamageNumbers *components.Checkbox
	autoPlay          *components.Checkbox
	showFPS           *components.Checkbox

	// 语言选择
	languageSelector widget.Enum
	languages        []LanguageOption

	// 按钮
	saveButton  *components.Button
	resetButton *components.Button

	// 列表布局
	list widget.List
}

type LanguageOption struct {
	Code string
	Name string
}

func NewSettingsTab() *SettingsTab {
	settings := store.GetGameStore().GetSettings()

	t := &SettingsTab{
		audioVolume:       components.NewSlider(float32(settings.AudioVolume)),
		musicVolume:       components.NewSlider(float32(settings.MusicVolume)),
		showDamageNumbers: components.NewCheckbox("显示伤害数字", settings.ShowDamageNumbers),
		autoPlay:          components.NewCheckbox("自动播放", settings.AutoPlay),
		showFPS:           components.NewCheckbox("显示FPS", settings.ShowFPS),
		saveButton:        components.NewButton("保存设置"),
		resetButton:       components.NewButton("重置"),
		languages: []LanguageOption{
			{Code: "zh_CN", Name: "简体中文"},
			{Code: "zh_TW", Name: "繁體中文"},
			{Code: "en", Name: "English"},
		},
	}

	// 设置当前语言
	t.languageSelector.Value = settings.Language

	return t
}

func (t *SettingsTab) Layout(gtx layout.Context) layout.Dimensions {
	// 检查按钮点击
	if t.saveButton.Clicked(gtx) {
		t.Save()
	}
	if t.resetButton.Clicked(gtx) {
		t.Reset()
	}

	// 检查复选框变化
	if t.showDamageNumbers.Changed(gtx) {
		// 实时更新，但不保存
	}
	if t.autoPlay.Changed(gtx) {
		// 实时更新，但不保存
	}
	if t.showFPS.Changed(gtx) {
		// 实时更新，但不保存
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// 标题
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := components.NewLabel("设置")
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
		// 滚动内容区域
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			t.list.Axis = layout.Vertical
			return t.list.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// 音频设置组
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return t.layoutAudioGroup(gtx)
					}),
					// 游戏设置组
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return t.layoutGameGroup(gtx)
					}),
					// 界面设置组
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return t.layoutUIGroup(gtx)
					}),
				)
			})
		}),
		// 底部按钮
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(16),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Bottom: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceEnd,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							t.resetButton.Color = theme.DefaultTheme.Surface
							return t.resetButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						t.saveButton.Color = theme.DefaultTheme.Primary
						return t.saveButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}

// 音频设置组
func (t *SettingsTab) layoutAudioGroup(gtx layout.Context) layout.Dimensions {
	return t.layoutCard(gtx, "音频设置", func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// 音效音量
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutSliderRow(gtx, "音效音量", t.audioVolume)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
			}),
			// 音乐音量
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutSliderRow(gtx, "音乐音量", t.musicVolume)
			}),
		)
	})
}

// 游戏设置组
func (t *SettingsTab) layoutGameGroup(gtx layout.Context) layout.Dimensions {
	return t.layoutCard(gtx, "游戏设置", func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.showDamageNumbers.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.autoPlay.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(12)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.showFPS.Layout(gtx)
			}),
		)
	})
}

// 界面设置组
func (t *SettingsTab) layoutUIGroup(gtx layout.Context) layout.Dimensions {
	return t.layoutCard(gtx, "界面设置", func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return t.layoutLanguageSelector(gtx)
			}),
		)
	})
}

// 卡片式布局
func (t *SettingsTab) layoutCard(gtx layout.Context, title string, content func(gtx layout.Context) layout.Dimensions) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			// 背景
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Max
				return drawRectWithRadius(gtx, theme.DefaultTheme.Surface, size, 8)
			}),
			// 内容
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(16),
					Left:   unit.Dp(16),
					Right:  unit.Dp(16),
					Bottom: unit.Dp(16),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						// 标题
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							titleLabel := components.NewLabel(title)
							titleLabel.Color = theme.DefaultTheme.Primary
							titleLabel.Size = 18
							return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return titleLabel.Layout(gtx)
							})
						}),
						// 内容
						layout.Rigid(content),
					)
				})
			}),
		)
	})
}

// 滑块行布局
func (t *SettingsTab) layoutSliderRow(gtx layout.Context, label string, slider *components.Slider) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel(label)
			lbl.Color = theme.DefaultTheme.Text
			lbl.Size = 14
			return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return lbl.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return slider.Layout(gtx)
		}),
	)
}

// 语言选择器
func (t *SettingsTab) layoutLanguageSelector(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := components.NewLabel("选择语言")
			lbl.Color = theme.DefaultTheme.Text
			lbl.Size = 14
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return lbl.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.layoutLanguageOption(gtx, "zh_CN", "简体中文")
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.layoutLanguageOption(gtx, "zh_TW", "繁體中文")
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return t.layoutLanguageOption(gtx, "en", "English")
				}),
			)
		}),
	)
}

// 语言选项
func (t *SettingsTab) layoutLanguageOption(gtx layout.Context, code, name string) layout.Dimensions {
	isSelected := t.languageSelector.Value == code

	return layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(4),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.languageSelector.Layout(gtx, code, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// 单选按钮样式
					size := gtx.Dp(unit.Dp(18))
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							s := image.Point{X: size, Y: size}
							return drawCircle(gtx, theme.DefaultTheme.Border, s)
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							if isSelected {
								innerSize := gtx.Dp(unit.Dp(10))
								s := image.Point{X: innerSize, Y: innerSize}
								return drawCircle(gtx, theme.DefaultTheme.Primary, s)
							}
							return layout.Dimensions{Size: image.Point{X: size, Y: size}}
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, 14, name)
					if isSelected {
						lbl.Color = toNRGBA(theme.DefaultTheme.Primary)
					} else {
						lbl.Color = toNRGBA(theme.DefaultTheme.Text)
					}
					return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return lbl.Layout(gtx)
					})
				}),
			)
		})
	})
}

func (t *SettingsTab) Save() {
	store.GetGameStore().UpdateSettings(func(s *types.Settings) {
		s.AudioVolume = float64(t.audioVolume.Value())
		s.MusicVolume = float64(t.musicVolume.Value())
		s.ShowDamageNumbers = t.showDamageNumbers.Checked()
		s.AutoPlay = t.autoPlay.Checked()
		s.ShowFPS = t.showFPS.Checked()
		s.Language = t.languageSelector.Value
	})
}

func (t *SettingsTab) Reset() {
	// 恢复默认设置
	t.audioVolume.SetValue(0.8)
	t.musicVolume.SetValue(0.6)
	t.showDamageNumbers.SetChecked(true)
	t.autoPlay.SetChecked(false)
	t.showFPS.SetChecked(false)
	t.languageSelector.Value = "zh_CN"

	// 保存到 store
	t.Save()
}
