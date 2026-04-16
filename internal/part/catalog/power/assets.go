package power

// File overview:
// assets selects pre-generated vectors and icons used by the power part.
// Subsystem: part catalog (power) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

// asset handles asset.
func (p *Power) asset() part.VectorAsset {
	if p.Kind == "gnd" {
		return gndAsset
	}
	return vccAsset
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
