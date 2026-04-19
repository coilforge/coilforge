package wire

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	toolbarIconOnce sync.Once
	toolbarIconImg  *ebiten.Image
)

func toolbarIcon() *ebiten.Image {
	toolbarIconOnce.Do(func() {
		const sz = 56
		img := ebiten.NewImage(sz, sz)
		line := color.RGBA{R: 220, G: 226, B: 236, A: 255}
		vector.StrokeLine(img, 10, float32(sz)/2, float32(sz)-10, float32(sz)/2, 3.5, line, false)
		vector.StrokeLine(img, float32(sz)/2, 10, float32(sz)/2, float32(sz)-10, 3.5, line, false)
		toolbarIconImg = img
	})
	return toolbarIconImg
}
