package indicator

// File overview:
// assets selects pre-generated vectors and icons used by the indicator part.
// Subsystem: part catalog (indicator) assets.
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
var assetOff = part.VectorAsset{Name: "indicator-off"}
var assetOn = part.VectorAsset{Name: "indicator-on"}

// asset handles asset.
func (self *Indicator) asset() part.VectorAsset {
	if self.Lit {
		return assetOn
	}
	return assetOff
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
