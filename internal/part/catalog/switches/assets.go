package switches

// File overview:
// assets selects pre-generated vectors and icons used by the switches part.
// Subsystem: part catalog (switches) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	_ "embed"

	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed toolbar_icon.png
var toolbarIconPNG []byte

var toolbarIconImage = part.LoadToolbarIconPNG(toolbarIconPNG, "switches/toolbar_icon.png")

// asset handles asset.
func (s *Switch) asset() part.VectorAsset {
	if s.effectiveClosed() {
		return switchClosedAsset
	}
	return switchOpenAsset
}

// toolbarIcon handles toolbar icon.
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
