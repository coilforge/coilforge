package rch

// File overview:
// assets selects pre-generated vectors and icons used by the rch part.
// Subsystem: part catalog (rch) assets.
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
var assetIdle = part.VectorAsset{Name: "rch-idle"}
var assetActive = part.VectorAsset{Name: "rch-active"}

// asset handles asset.
func (self *RCH) asset() part.VectorAsset {
	if self.Active {
		return assetActive
	}
	return assetIdle
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
