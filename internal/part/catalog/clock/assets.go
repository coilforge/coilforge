package clock

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (c *Clock) asset() part.VectorAsset {
	if c.OutputHigh {
		return clockHighAsset
	}
	return clockLowAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
