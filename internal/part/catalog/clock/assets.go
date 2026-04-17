package clock

// File overview:
// assets selects pre-generated vectors and icons used by the clock part.
// Subsystem: part catalog (clock) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	_ "embed"

	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed toolbar_icon.png
var toolbarIconPNG []byte

var toolbarIconImage = part.LoadToolbarIconPNG(toolbarIconPNG)
var assetLow = part.VectorAsset{Name: "clock-low"}
var assetHigh = part.VectorAsset{Name: "clock-high"}

// asset handles asset.
func (self *Clock) asset() part.VectorAsset {
	if self.OutputHigh {
		return assetHigh
	}
	return assetLow
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
