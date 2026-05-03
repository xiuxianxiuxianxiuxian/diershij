package tabs

import (
	"image"
	"image/color"

	"cultivation-client/internal/gui/components"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

var th = material.NewTheme()

func toNRGBA(c color.RGBA) color.NRGBA {
	return components.ToNRGBA(c)
}

func drawRect(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	return components.DrawRect(gtx, c, size)
}

func drawRectWithRadius(gtx layout.Context, c color.RGBA, size image.Point, radius int) layout.Dimensions {
	return components.DrawRectWithRadius(gtx, c, size, radius)
}

func drawCircle(gtx layout.Context, c color.RGBA, size image.Point) layout.Dimensions {
	return components.DrawCircle(gtx, c, size)
}

func drawDivider(gtx layout.Context, c color.RGBA) layout.Dimensions {
	return components.DrawDivider(gtx, c)
}

func drawCardBackground(gtx layout.Context, c color.RGBA) layout.Dimensions {
	return components.DrawCardBackground(gtx, c)
}

func drawAvatar(gtx layout.Context, c color.RGBA) layout.Dimensions {
	return components.DrawAvatar(gtx, c)
}

func btn(text string, c color.RGBA) *components.Button {
	b := components.NewButton(text)
	b.Color = c
	return b
}
