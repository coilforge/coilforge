package rch

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (r *RCH) asset() part.VectorAsset {
	if r.Active {
		return rchActiveAsset
	}
	return rchIdleAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
