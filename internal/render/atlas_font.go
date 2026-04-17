package render

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/unicode/norm"
)

//go:embed fonts/ui_label_atlas.png
var uiLabelAtlasPNG []byte

//go:embed fonts/ui_label_atlas.json
var uiLabelAtlasJSON []byte

type atlasGlyph struct {
	X       int     `json:"x"`
	Y       int     `json:"y"`
	W       int     `json:"w"`
	H       int     `json:"h"`
	Advance float64 `json:"advance"`
	OffX    float64 `json:"offX"`
	OffY    float64 `json:"offY"`
}

type atlasMeta struct {
	Ascent   float64 `json:"ascent"`
	Descent  float64 `json:"descent"`
	LineH    float64 `json:"lineH"`
	Fallback int     `json:"fallback"`
	Glyphs   []struct {
		Rune int `json:"rune"`
		atlasGlyph
	} `json:"glyphs"`
}

type atlasFont struct {
	img      *ebiten.Image
	ascent   float64
	lineH    float64
	fallback rune
	glyphs   map[rune]atlasGlyph
}

var (
	uiAtlasOnce sync.Once
	uiAtlas     atlasFont
)

func uiLabelAtlas() atlasFont {
	uiAtlasOnce.Do(func() {
		var meta atlasMeta
		if err := json.Unmarshal(uiLabelAtlasJSON, &meta); err != nil {
			panic("render: decode ui label atlas json failed: " + err.Error())
		}
		imgDecoded, err := png.Decode(bytes.NewReader(uiLabelAtlasPNG))
		if err != nil {
			panic("render: decode ui label atlas png failed: " + err.Error())
		}
		img := ebiten.NewImageFromImage(imgDecoded)
		gm := make(map[rune]atlasGlyph, len(meta.Glyphs))
		for _, g := range meta.Glyphs {
			gm[rune(g.Rune)] = g.atlasGlyph
		}
		uiAtlas = atlasFont{
			img:      img,
			ascent:   meta.Ascent,
			lineH:    meta.LineH,
			fallback: rune(meta.Fallback),
			glyphs:   gm,
		}
	})
	return uiAtlas
}

func atlasMeasure(text string, f atlasFont) (float64, float64) {
	n := norm.NFC.String(text)
	maxW := 0.0
	lineW := 0.0
	lines := 1
	for _, r := range n {
		if r == '\n' {
			if lineW > maxW {
				maxW = lineW
			}
			lineW = 0
			lines++
			continue
		}
		g, ok := f.glyphs[r]
		if !ok {
			g = f.glyphs[f.fallback]
		}
		lineW += g.Advance
	}
	if lineW > maxW {
		maxW = lineW
	}
	return maxW, float64(lines) * f.lineH
}

func drawAtlasText(dst *ebiten.Image, text string, x, y float64, clr color.Color) {
	f := uiLabelAtlas()
	n := norm.NFC.String(text)
	penX := x
	baselineY := y + f.ascent
	for _, r := range n {
		if r == '\n' {
			penX = x
			baselineY += f.lineH
			continue
		}
		g, ok := f.glyphs[r]
		if !ok {
			g = f.glyphs[f.fallback]
		}
		if g.W > 0 && g.H > 0 {
			sub := f.img.SubImage(imageRect(g.X, g.Y, g.W, g.H)).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(penX+g.OffX, baselineY+g.OffY)
			op.ColorScale.ScaleWithColor(clr)
			dst.DrawImage(sub, op)
		}
		penX += g.Advance
	}
}

func imageRect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

func normalizeUIString(s string) string {
	return norm.NFC.String(s)
}
