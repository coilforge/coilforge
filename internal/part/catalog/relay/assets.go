package relay

// File overview:
// assets selects pre-generated vectors and icons used by the relay part.
// Subsystem: part catalog (relay) assets.
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
var assetIdle = part.VectorAsset{Name: "relay-idle"}
var assetActive = part.VectorAsset{Name: "relay-active"}

// asset handles asset.
func (self *Relay) asset() part.VectorAsset {
	if self.CoilActive {
		return assetActive
	}
	return assetIdle
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
