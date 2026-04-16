package diode

// File overview:
// assets selects pre-generated vectors and icons used by the diode part.
// Subsystem: part catalog (diode) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

// asset handles asset.
func (d *Diode) asset() part.VectorAsset {
	return diodeAsset
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
