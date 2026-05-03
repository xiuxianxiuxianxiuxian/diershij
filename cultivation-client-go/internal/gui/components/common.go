package components

import (
	"fmt"
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

var (
	th = material.NewTheme()
)

func toNRGBA(c color.RGBA) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

// ToNRGBA converts RGBA to NRGBA.
func ToNRGBA(c color.RGBA) color.NRGBA {
	return toNRGBA(c)
}

// DrawRect draws a filled rectangle.
func DrawRect(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// drawRect is an alias for DrawRect for internal use.
func drawRect(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	return DrawRect(gtx, c, size)
}

// DrawRectWithRadius draws a filled rounded rectangle.
func DrawRectWithRadius(gtx layout.Context, c color.RGBA, size image.Point, radius int) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{Max: size},
		SE:   radius,
		SW:   radius,
		NE:   radius,
		NW:   radius,
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func drawRectWithRadius(gtx layout.Context, c color.RGBA, size image.Point, radius int) layout.Dimensions {
	return DrawRectWithRadius(gtx, c, size, radius)
}

// DrawCircle draws a filled circle.
func DrawCircle(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Ellipse{
		Max: size,
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func drawCircle(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	return DrawCircle(gtx, c, size)
}

// DrawDivider draws a horizontal 1px divider line.
func DrawDivider(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Constraints.Max
	size.Y = gtx.Dp(unit.Dp(1))
	return DrawRect(gtx, c, size)
}

// DrawCardBackground draws a card-style rounded background.
func DrawCardBackground(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Constraints.Max
	return DrawRectWithRadius(gtx, c, size, 12)
}

// DrawAvatar draws a circular avatar placeholder.
func DrawAvatar(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Dp(unit.Dp(80))
	outerSize := size + gtx.Dp(unit.Dp(4))
	defer clip.Ellipse{
		Max: image.Point{X: outerSize, Y: outerSize},
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(theme.DefaultTheme.PrimaryVariant)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	innerOffset := (outerSize - size) / 2
	return layout.Inset{
		Top:    unit.Dp(innerOffset),
		Left:   unit.Dp(innerOffset),
		Right:  unit.Dp(0),
		Bottom: unit.Dp(0),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		defer clip.Ellipse{
			Max: image.Point{X: size, Y: size},
		}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Point{X: size, Y: size}}
	})
}

// DrawCheckmark draws a filled checkmark square.
func DrawCheckmark(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func drawCheckmark(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	return DrawCheckmark(gtx, c, size)
}

type Button struct {
	Text   string
	widget widget.Clickable
	Color  color.RGBA
	Inset  layout.Inset
}

func NewButton(text string) *Button {
	return &Button{
		Text:  text,
		Color: theme.DefaultTheme.Primary,
		Inset: layout.UniformInset(unit.Dp(12)),
	}
}

func (b *Button) Layout(gtx layout.Context) layout.Dimensions {
	btn := material.Button(th, &b.widget, b.Text)
	btn.Background = toNRGBA(b.Color)
	btn.Inset = b.Inset
	return btn.Layout(gtx)
}

func (b *Button) Clicked(gtx layout.Context) bool {
	return b.widget.Clicked(gtx)
}

type InputField struct {
	Placeholder string
	Password    bool
	widget      widget.Editor
}

func NewInputField(placeholder string) *InputField {
	return &InputField{
		Placeholder: placeholder,
		Password:    false,
	}
}

func (i *InputField) Text() string {
	return i.widget.Text()
}

func (i *InputField) SetText(text string) {
	i.widget.SetText(text)
}

func (i *InputField) Layout(gtx layout.Context) layout.Dimensions {
	editor := material.Editor(th, &i.widget, i.Placeholder)
	border := widget.Border{
		Color:        toNRGBA(theme.DefaultTheme.Border),
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(1),
	}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return editor.Layout(gtx)
		})
	})
}

type Label struct {
	Text  string
	Color color.RGBA
	Size  unit.Sp
}

func NewLabel(text string) *Label {
	return &Label{
		Text:  text,
		Color: theme.DefaultTheme.Text,
		Size:  16,
	}
}

func (l *Label) Layout(gtx layout.Context) layout.Dimensions {
	lbl := material.Label(th, l.Size, l.Text)
	lbl.Color = toNRGBA(l.Color)
	return lbl.Layout(gtx)
}

type ProgressBar struct {
	Value     float32
	Max       float32
	Color     color.RGBA
	BgColor   color.RGBA
	ShowText  bool
	Height    unit.Dp
	TextColor color.RGBA
}

func NewProgressBar(value, max float32) *ProgressBar {
	return &ProgressBar{
		Value:     value,
		Max:       max,
		Color:     theme.DefaultTheme.Success,
		BgColor:   theme.DefaultTheme.Border,
		ShowText:  true,
		Height:    unit.Dp(16),
		TextColor: theme.DefaultTheme.Text,
	}
}

func (p *ProgressBar) SetColor(color color.RGBA) {
	p.Color = color
}

func (p *ProgressBar) SetBgColor(color color.RGBA) {
	p.BgColor = color
}

func (p *ProgressBar) SetShowText(show bool) {
	p.ShowText = show
}

func (p *ProgressBar) Layout(gtx layout.Context) layout.Dimensions {
	var progress float32
	if p.Max > 0 {
		progress = float32(p.Value) / float32(p.Max)
		if progress < 0 {
			progress = 0
		}
		if progress > 1 {
			progress = 1
		}
	}

	barHeight := gtx.Dp(p.Height)

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Max
			size.Y = barHeight
			return drawRectWithRadius(gtx, p.BgColor, size, 4)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Constraints.Max
					size.Y = barHeight
					size.X = int(float32(size.X) * progress)
					if size.X > 0 {
						return drawRectWithRadius(gtx, p.Color, size, 4)
					}
					return layout.Dimensions{Size: image.Point{X: 0, Y: barHeight}}
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: barHeight}}
				}),
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !p.ShowText {
				return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: barHeight}}
			}
			percentage := int(progress * 100)
			text := fmt.Sprintf("%d%%", percentage)
			lbl := material.Label(th, 12, text)
			lbl.Color = toNRGBA(p.TextColor)
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

type ListItem struct {
	Title    string
	Subtitle string
	Icon     *image.RGBA
	Selected bool
	widget   widget.Clickable
}

func NewListItem(title, subtitle string) *ListItem {
	return &ListItem{
		Title:    title,
		Subtitle: subtitle,
		Selected: false,
	}
}

func (l *ListItem) Layout(gtx layout.Context) layout.Dimensions {
	return l.widget.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		bg := theme.DefaultTheme.Surface
		if l.Selected {
			bg = theme.DefaultTheme.Hover
		}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return drawRect(gtx, bg, gtx.Constraints.Max)
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, 16, l.Title)
							lbl.Color = toNRGBA(theme.DefaultTheme.Text)
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, 12, l.Subtitle)
							lbl.Color = toNRGBA(theme.DefaultTheme.TextSecondary)
							return lbl.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}
