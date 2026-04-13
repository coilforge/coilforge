package wire

import (
	"coilforge/internal/core"
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (w *Wire) asset() part.VectorAsset {
	switch w.State {
	case core.NetHigh:
		return wireHighAsset
	case core.NetLow:
		return wireLowAsset
	case core.NetShort:
		return wireShortAsset
	default:
		return wireFloatAsset
	}
}

func toolbarIcon() *ebiten.Image {
	return nil
}
