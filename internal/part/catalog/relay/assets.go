package relay

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (r *Relay) asset() part.VectorAsset {
	if r.CoilActive {
		return relayActiveAsset
	}
	return relayIdleAsset
}

func toolbarIcon() *ebiten.Image {
	return nil
}
