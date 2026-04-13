package render

import (
	"coilforge/internal/core"
	"image/color"
)

var DarkMode = true

func WireColor(state int) color.RGBA {
	switch state {
	case core.NetHigh:
		return color.RGBA{R: 255, G: 180, B: 64, A: 255}
	case core.NetLow:
		return color.RGBA{R: 96, G: 180, B: 255, A: 255}
	case core.NetShort:
		return color.RGBA{R: 255, G: 64, B: 64, A: 255}
	default:
		return color.RGBA{R: 160, G: 160, B: 160, A: 255}
	}
}

func GridColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 42, G: 46, B: 54, A: 255}
	}
	return color.RGBA{R: 224, G: 228, B: 232, A: 255}
}

func SelectionColor() color.RGBA {
	return color.RGBA{R: 255, G: 208, B: 64, A: 255}
}

func GhostTint() color.RGBA {
	return color.RGBA{R: 255, G: 255, B: 255, A: 144}
}
