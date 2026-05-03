package theme

import "image/color"

type Theme struct {
	Background       color.RGBA
	Surface          color.RGBA
	Primary          color.RGBA
	PrimaryVariant   color.RGBA
	Secondary        color.RGBA
	Text             color.RGBA
	TextSecondary    color.RGBA
	Error            color.RGBA
	Success          color.RGBA
	Warning          color.RGBA
	Border           color.RGBA
	Hover            color.RGBA
	Active           color.RGBA
	Disabled         color.RGBA
}

var DefaultTheme = Theme{
	Background:       color.RGBA{R: 20, G: 20, B: 30, A: 255},
	Surface:          color.RGBA{R: 30, G: 30, B: 45, A: 255},
	Primary:          color.RGBA{R: 100, G: 150, B: 255, A: 255},
	PrimaryVariant:   color.RGBA{R: 70, G: 120, B: 220, A: 255},
	Secondary:        color.RGBA{R: 180, G: 100, B: 255, A: 255},
	Text:             color.RGBA{R: 255, G: 255, B: 255, A: 255},
	TextSecondary:    color.RGBA{R: 180, G: 180, B: 200, A: 255},
	Error:            color.RGBA{R: 255, G: 80, B: 80, A: 255},
	Success:          color.RGBA{R: 80, G: 200, B: 120, A: 255},
	Warning:          color.RGBA{R: 255, G: 200, B: 80, A: 255},
	Border:           color.RGBA{R: 60, G: 60, B: 80, A: 255},
	Hover:            color.RGBA{R: 40, G: 40, B: 60, A: 255},
	Active:           color.RGBA{R: 50, G: 50, B: 70, A: 255},
	Disabled:         color.RGBA{R: 80, G: 80, B: 100, A: 255},
}

var DarkTheme = Theme{
	Background:       color.RGBA{R: 10, G: 10, B: 15, A: 255},
	Surface:          color.RGBA{R: 20, G: 20, B: 30, A: 255},
	Primary:          color.RGBA{R: 80, G: 130, B: 255, A: 255},
	PrimaryVariant:   color.RGBA{R: 60, G: 100, B: 200, A: 255},
	Secondary:        color.RGBA{R: 160, G: 80, B: 255, A: 255},
	Text:             color.RGBA{R: 240, G: 240, B: 255, A: 255},
	TextSecondary:    color.RGBA{R: 160, G: 160, B: 180, A: 255},
	Error:            color.RGBA{R: 255, G: 60, B: 60, A: 255},
	Success:          color.RGBA{R: 60, G: 180, B: 100, A: 255},
	Warning:          color.RGBA{R: 255, G: 180, B: 60, A: 255},
	Border:           color.RGBA{R: 40, G: 40, B: 60, A: 255},
	Hover:            color.RGBA{R: 30, G: 30, B: 45, A: 255},
	Active:           color.RGBA{R: 40, G: 40, B: 55, A: 255},
	Disabled:         color.RGBA{R: 60, G: 60, B: 80, A: 255},
}

func (t *Theme) WithPrimary(primary color.RGBA) *Theme {
	t.Primary = primary
	return t
}

func (t *Theme) WithBackground(background color.RGBA) *Theme {
	t.Background = background
	return t
}
