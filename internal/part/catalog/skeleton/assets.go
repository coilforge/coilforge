package skeleton

// File overview:
// assets selects pre-generated vectors and icons used by the skeleton part.
// Subsystem: part catalog (skeleton) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

var asset = part.VectorAsset{Name: "template"}

// asset handles asset.
func (self *Template) asset() part.VectorAsset {
	return asset
}

// ToolbarIcon handles toolbar icon.
func ToolbarIcon() *ebiten.Image {
	return nil
}
