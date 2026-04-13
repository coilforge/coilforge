package indicator

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (ind *Indicator) asset() part.VectorAsset {
	if ind.Lit {
		return indicatorOnAsset
	}
	return indicatorOffAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
