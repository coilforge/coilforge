package power

// File overview:
// assets selects pre-generated vectors and icons used by the power part.
// Subsystem: part catalog (power) assets.
// It supports draw and toolbar registration without runtime SVG parsing.
// Flow position: static visual resource selector beneath part drawing.

import (
	_ "embed"

	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed toolbar_icon_vcc.png
var toolbarIconVCCPNG []byte

//go:embed toolbar_icon_gnd.png
var toolbarIconGNDPNG []byte

var toolbarIconVCCImage = part.LoadToolbarIconPNG(toolbarIconVCCPNG, "power/toolbar_icon_vcc.png")
var toolbarIconGNDImage = part.LoadToolbarIconPNG(toolbarIconGNDPNG, "power/toolbar_icon_gnd.png")

// asset handles asset.
func (p *Power) asset() part.VectorAsset {
	if p.Kind == "gnd" {
		return gndAsset
	}
	return vccAsset
}

// toolbarIconVCC handles toolbar icon.
func toolbarIconVCC() *ebiten.Image {
	return toolbarIconVCCImage
}

// toolbarIconGND handles toolbar icon.
func toolbarIconGND() *ebiten.Image {
	return toolbarIconGNDImage
}
