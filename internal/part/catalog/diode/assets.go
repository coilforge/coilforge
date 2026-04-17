package diode

// File overview:
// assets selects pre-generated vectors and icons used by this part.
// Subsystem: part catalog assets.
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
var asset = part.VectorAsset{Name: "diode"}

// asset handles asset.
func (self *Diode) asset() part.VectorAsset {
	return asset
}

// toolbarIcon returns the pre-rasterized toolbar bitmap (e.g. 84×84 source, scaled in chrome).
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
