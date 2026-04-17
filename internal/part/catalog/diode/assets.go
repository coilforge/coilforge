package diode

// File overview:
// assets selects pre-generated vectors and icons used by the diode part.
// Subsystem: part catalog (diode) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	_ "embed"

	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed toolbar_icon.png
var toolbarIconPNG []byte

var toolbarIconImage = part.LoadToolbarIconPNG(toolbarIconPNG, "diode/toolbar_icon.png")

// asset handles asset.
func (d *Diode) asset() part.VectorAsset {
	return diodeAsset
}

// toolbarIcon returns the pre-rasterized toolbar bitmap (e.g. 84×84 source, scaled in chrome).
func toolbarIcon() *ebiten.Image {
	return toolbarIconImage
}
