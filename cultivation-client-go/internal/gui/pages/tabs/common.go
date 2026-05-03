package tabs

import (
	"image"
	"image/color"

	"cultivation-client/internal/gui/theme"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var th = material.NewTheme()

func toNRGBA(c color.RGBA) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

// 绘制卡片背景
func drawCardBackground(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Constraints.Max
	defer clip.RRect{
		Rect: image.Rectangle{Max: size},
		SE:   12,
		SW:   12,
		NE:   12,
		NW:   12,
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// 绘制分隔线
func drawDivider(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Constraints.Max
	size.Y = gtx.Dp(unit.Dp(1))
	return drawRect(gtx, c, size)
}

// 绘制矩形
func drawRect(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// 绘制圆角矩形
func drawRectWithRadius(gtx layout.Context, c color.RGBA, size image.Point, radius int) layout.Dimensions {
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

// 绘制头像
func drawAvatar(gtx layout.Context, c color.RGBA) layout.Dimensions {
	size := gtx.Dp(unit.Dp(80))
	// 外圈装饰
	outerSize := size + gtx.Dp(unit.Dp(4))
	defer clip.Ellipse{
		Max: image.Point{X: outerSize, Y: outerSize},
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(theme.DefaultTheme.PrimaryVariant)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// 内圈头像区域
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

// 绘制圆形
func drawCircle(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	defer clip.Ellipse{
		Max: size,
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: toNRGBA(c)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
