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

var assetFloat = part.VectorAsset{Name: "wire-float"}
var assetLow = part.VectorAsset{Name: "wire-low"}
var assetHigh = part.VectorAsset{Name: "wire-high"}
var assetShort = part.VectorAsset{Name: "wire-short"}

// asset handles asset.
func (self *Wire) asset() part.VectorAsset {
	switch self.State {
	case core.NetHigh:
		return assetHigh
	case core.NetLow:
		return assetLow
	case core.NetShort:
		return assetShort
	default:
		return assetFloat
	}
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return nil
}
