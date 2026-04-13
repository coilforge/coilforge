package diode

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (d *Diode) asset() part.VectorAsset {
	return diodeAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
