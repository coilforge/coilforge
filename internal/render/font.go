package render

import (
	"log"
	"math"
	"sync"

	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

var (
	toolbarLabelTTFOnce sync.Once
	toolbarLabelTTF     *opentype.Font
	toolbarFaceMu       sync.Mutex
	toolbarFaceCache    = map[int]textv2.Face{}
)

func toolbarLabelFont() *opentype.Font {
	toolbarLabelTTFOnce.Do(func() {
		ttf, err := opentype.Parse(goregular.TTF)
		if err != nil {
			log.Fatalf("render: parse toolbar label TTF: %v", err)
		}
		toolbarLabelTTF = ttf
	})
	return toolbarLabelTTF
}

// toolbarLabelFaceForScale returns a cached face rasterized for the current device scale.
// This keeps text crisp on HiDPI displays while preserving logical UI sizing.
func toolbarLabelFaceForScale(deviceScale float64) textv2.Face {
	if deviceScale <= 0 {
		deviceScale = 1
	}
	scaleKey := int(math.Round(deviceScale * 100))

	toolbarFaceMu.Lock()
	defer toolbarFaceMu.Unlock()

	if face, ok := toolbarFaceCache[scaleKey]; ok {
		return face
	}
	ttf := toolbarLabelFont()
	goFace, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    10 * deviceScale,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("render: create toolbar label face (scale %.2f): %v", deviceScale, err)
	}
	face := textv2.NewGoXFace(goFace)
	toolbarFaceCache[scaleKey] = face
	return face
}
