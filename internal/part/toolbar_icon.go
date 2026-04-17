package part

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// LoadToolbarIconPNG decodes embedded PNG bytes into an Ebiten image.
func LoadToolbarIconPNG(pngData []byte, assetName string) *ebiten.Image {
	if len(pngData) == 0 {
		log.Fatalf("part: empty toolbar icon data for %s", assetName)
	}
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		log.Fatalf("part: decode toolbar icon %s: %v", assetName, err)
	}
	return ebiten.NewImageFromImage(invertIconRGB(img))
}

// invertIconRGB flips RGB channels while preserving alpha coverage.
// This lets source artwork stay black-on-transparent but render as white-ready glyphs.
func invertIconRGB(src image.Image) *image.NRGBA {
	b := src.Bounds()
	out := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			c.R = 255 - c.R
			c.G = 255 - c.G
			c.B = 255 - c.B
			out.SetNRGBA(x, y, c)
		}
	}
	return out
}
