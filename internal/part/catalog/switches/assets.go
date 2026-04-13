package switches

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (s *Switch) asset() part.VectorAsset {
	if s.effectiveClosed() {
		return switchClosedAsset
	}
	return switchOpenAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
