package switches

// File overview:
// assets selects pre-generated vectors and icons used by the switches part.
// Subsystem: part catalog (switches) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

// asset handles asset.
func (s *Switch) asset() part.VectorAsset {
	if s.effectiveClosed() {
		return switchClosedAsset
	}
	return switchOpenAsset
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
