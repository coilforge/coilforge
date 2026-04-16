package indicator

// File overview:
// assets selects pre-generated vectors and icons used by the indicator part.
// Subsystem: part catalog (indicator) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

// asset handles asset.
func (ind *Indicator) asset() part.VectorAsset {
	if ind.Lit {
		return indicatorOnAsset
	}
	return indicatorOffAsset
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
