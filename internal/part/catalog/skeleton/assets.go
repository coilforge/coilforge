package skeleton

import (
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

func (t *Template) asset() part.VectorAsset {
	return templateAsset
}

func ToolbarIcon() *ebiten.Image {
	return nil
}
