package power

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (p *Power) asset() part.VectorAsset {
	if p.Kind == "gnd" {
		return gndAsset
	}
	return vccAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
