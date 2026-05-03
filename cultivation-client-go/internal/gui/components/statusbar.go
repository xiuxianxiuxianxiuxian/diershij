package components

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"cultivation-client/internal/gui/theme"
	"cultivation-client/internal/store"
	"gioui.org/layout"
	"gioui.org/unit"
)

type StatusBar struct {
	Height unit.Dp
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		Height: 40,
	}
}

func (s *StatusBar) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return drawRect(gtx, theme.DefaultTheme.Surface, gtx.Constraints.Max)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					layout.Rigid(s.layoutLeftSection),
					layout.Rigid(s.layoutRightSection),
				)
			})
		}),
	)
}

func (s *StatusBar) layoutLeftSection(gtx layout.Context) layout.Dimensions {
	char := store.GetGameStore().GetCharacter()
	if char == nil {
		return layout.Dimensions{Size: image.Point{X: 0, Y: int(s.Height)}}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			nameLabel := NewLabel(char.Name)
			nameLabel.Size = 16
			nameLabel.Color = theme.DefaultTheme.Text
			return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return nameLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hpLabel := NewLabel("HP:")
					hpLabel.Size = 12
					hpLabel.Color = theme.DefaultTheme.TextSecondary
					return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return hpLabel.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hpBar := NewProgressBar(float32(char.Health), float32(char.MaxHealth))
					hpBar.Color = color.RGBA{R: 255, G: 80, B: 80, A: 255}
					return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = 100
						gtx.Constraints.Min.X = 100
						return hpBar.Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					energyLabel := NewLabel("MP:")
					energyLabel.Size = 12
					energyLabel.Color = theme.DefaultTheme.TextSecondary
					return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return energyLabel.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					energyBar := NewProgressBar(float32(char.Energy), float32(char.MaxEnergy))
					energyBar.Color = color.RGBA{R: 80, G: 150, B: 255, A: 255}
					return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = 100
						gtx.Constraints.Min.X = 100
						return energyBar.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (s *StatusBar) layoutRightSection(gtx layout.Context) layout.Dimensions {
	world := store.GetGameStore().GetWorld()
	onlineCount := 0
	if world != nil {
		onlineCount = world.PlayersOnline
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			onlineLabel := NewLabel(fmt.Sprintf("在线: %d", onlineCount))
			onlineLabel.Size = 14
			onlineLabel.Color = theme.DefaultTheme.TextSecondary
			return layout.Inset{Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return onlineLabel.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			timeLabel := NewLabel(time.Now().Format("15:04:05"))
			timeLabel.Size = 14
			timeLabel.Color = theme.DefaultTheme.TextSecondary
			return timeLabel.Layout(gtx)
		}),
	)
}

func (s *StatusBar) GetCharacterHP() (int, int) {
	char := store.GetGameStore().GetCharacter()
	if char == nil {
		return 0, 0
	}
	return char.Health, char.MaxHealth
}

func (s *StatusBar) GetCharacterEnergy() (int, int) {
	char := store.GetGameStore().GetCharacter()
	if char == nil {
		return 0, 0
	}
	return char.Energy, char.MaxEnergy
}

func (s *StatusBar) GetOnlinePlayers() int {
	world := store.GetGameStore().GetWorld()
	if world == nil {
		return 0
	}
	return world.PlayersOnline
}

func (s *StatusBar) GetCurrentTime() string {
	return time.Now().Format("15:04:05")
}
