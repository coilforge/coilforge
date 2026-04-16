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

// asset handles asset.
func (t *Template) asset() part.VectorAsset {
	return templateAsset
}

// ToolbarIcon handles toolbar icon.
func ToolbarIcon() *ebiten.Image {
	return nil
}
