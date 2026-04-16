package wire

// File overview:
// assets selects pre-generated vectors and icons used by the wire part.
// Subsystem: part catalog (wire) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

// asset handles asset.
func (w *Wire) asset() part.VectorAsset {
	switch w.State {
	case core.NetHigh:
		return wireHighAsset
	case core.NetLow:
		return wireLowAsset
	case core.NetShort:
		return wireShortAsset
	default:
		return wireFloatAsset
	}
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
