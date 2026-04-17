package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/unicode/norm"
)

type glyphMeta struct {
	Rune    int     `json:"rune"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	W       int     `json:"w"`
	H       int     `json:"h"`
	Advance float64 `json:"advance"`
	OffX    float64 `json:"offX"`
	OffY    float64 `json:"offY"`
}

type atlasMeta struct {
	Ascent   float64     `json:"ascent"`
	Descent  float64     `json:"descent"`
	LineH    float64     `json:"lineH"`
	Fallback int         `json:"fallback"`
	Glyphs   []glyphMeta `json:"glyphs"`
}

type renderedGlyph struct {
	r       rune
	img     *image.Alpha
	advance float64
	offX    float64
	offY    float64
	w       int
	h       int
}

const uiLatinBasic = "" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"0123456789" +
	" " +
	"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" +
	"ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝ" +
	"àáâãäåæçèéêëìíîïñòóôõöøùúûüýÿ" +
	"ÐðÞþßŒœŠšŽžŸ" +
	"€£¥°µ§«»–—…"

func main() {
	fontPath := flag.String("font", "internal/render/fonts/Inter-Variable.ttf", "Path to source TTF/OTF")
	outPNG := flag.String("out-png", "internal/render/fonts/ui_label_atlas.png", "Output atlas PNG")
	outJSON := flag.String("out-json", "internal/render/fonts/ui_label_atlas.json", "Output atlas metadata JSON")
	size := flag.Float64("size", 10, "Font size in px-ish units")
	dpi := flag.Float64("dpi", 72, "Font DPI")
	texW := flag.Int("tex-w", 512, "Atlas width")
	padding := flag.Int("padding", 1, "Padding around glyph bitmaps")
	flag.Parse()

	fontBytes, err := os.ReadFile(*fontPath)
	if err != nil {
		fail("read font: %v", err)
	}
	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		fail("parse font: %v", err)
	}
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    *size,
		DPI:     *dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		fail("create face: %v", err)
	}
	defer face.Close()

	metrics := face.Metrics()
	ascent := float64(metrics.Ascent) / 64.0
	descent := float64(metrics.Descent) / 64.0
	lineH := float64(metrics.Height) / 64.0
	if lineH <= 0 {
		lineH = ascent + descent
	}
	if lineH <= 0 {
		lineH = *size * 1.2
	}

	runes := charsetRunes(uiLatinBasic)
	fallback := '?'
	if !containsRune(runes, fallback) {
		runes = append(runes, fallback)
	}

	glyphs := make([]renderedGlyph, 0, len(runes))
	for _, r := range runes {
		g := rasterizeGlyph(face, r, ascent)
		glyphs = append(glyphs, g)
	}

	sort.Slice(glyphs, func(i, j int) bool {
		if glyphs[i].h == glyphs[j].h {
			return glyphs[i].w > glyphs[j].w
		}
		return glyphs[i].h > glyphs[j].h
	})

	meta, atlas := packGlyphs(glyphs, *texW, *padding, ascent, descent, lineH, fallback)

	if err := os.MkdirAll(filepath.Dir(*outPNG), 0o755); err != nil {
		fail("mkdir png dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(*outJSON), 0o755); err != nil {
		fail("mkdir json dir: %v", err)
	}

	pf, err := os.Create(*outPNG)
	if err != nil {
		fail("create png: %v", err)
	}
	if err := png.Encode(pf, atlas); err != nil {
		_ = pf.Close()
		fail("encode png: %v", err)
	}
	_ = pf.Close()

	jf, err := os.Create(*outJSON)
	if err != nil {
		fail("create json: %v", err)
	}
	enc := json.NewEncoder(jf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(meta); err != nil {
		_ = jf.Close()
		fail("encode json: %v", err)
	}
	_ = jf.Close()

	fmt.Printf("Generated atlas: %s (%dx%d), glyphs=%d\n", *outPNG, atlas.Bounds().Dx(), atlas.Bounds().Dy(), len(meta.Glyphs))
	fmt.Printf("Generated meta : %s\n", *outJSON)
}

func rasterizeGlyph(face font.Face, r rune, ascent float64) renderedGlyph {
	const canvasPad = 16
	canvasW := 256
	canvasH := 256

	img := image.NewAlpha(image.Rect(0, 0, canvasW, canvasH))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Alpha{A: 255}),
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(canvasPad),
			Y: fixed.I(canvasPad + int(ascent+0.5)),
		},
	}
	d.DrawString(string(r))

	b := alphaBounds(img)
	adv := float64(d.MeasureString(string(r))) / 64.0
	if b.Empty() {
		return renderedGlyph{
			r:       r,
			img:     image.NewAlpha(image.Rect(0, 0, 1, 1)),
			advance: adv,
			offX:    0,
			offY:    0,
			w:       0,
			h:       0,
		}
	}

	crop := image.NewAlpha(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(crop, crop.Bounds(), img, b.Min, draw.Src)

	baseX := canvasPad
	baseY := canvasPad + int(ascent+0.5)
	return renderedGlyph{
		r:       r,
		img:     crop,
		advance: adv,
		offX:    float64(b.Min.X - baseX),
		offY:    float64(b.Min.Y - baseY),
		w:       b.Dx(),
		h:       b.Dy(),
	}
}

func packGlyphs(glyphs []renderedGlyph, texW, pad int, ascent, descent, lineH float64, fallback rune) (atlasMeta, *image.NRGBA) {
	x, y := pad, pad
	rowH := 0
	placed := make([]glyphMeta, 0, len(glyphs))

	maxY := 0
	for _, g := range glyphs {
		w, h := g.w, g.h
		if w > 0 && h > 0 && x+w+pad > texW {
			x = pad
			y += rowH + pad
			rowH = 0
		}
		if h > rowH {
			rowH = h
		}
		if y+h+pad > maxY {
			maxY = y + h + pad
		}
		placed = append(placed, glyphMeta{
			Rune:    int(g.r),
			X:       x,
			Y:       y,
			W:       w,
			H:       h,
			Advance: g.advance,
			OffX:    g.offX,
			OffY:    g.offY,
		})
		x += w + pad
	}
	if maxY < pad*2 {
		maxY = pad * 2
	}
	atlas := image.NewNRGBA(image.Rect(0, 0, texW, maxY))

	byRune := map[rune]renderedGlyph{}
	for _, g := range glyphs {
		byRune[g.r] = g
	}

	for _, m := range placed {
		g := byRune[rune(m.Rune)]
		if m.W == 0 || m.H == 0 {
			continue
		}
		for yy := 0; yy < m.H; yy++ {
			for xx := 0; xx < m.W; xx++ {
				a := g.img.AlphaAt(xx, yy).A
				i := atlas.PixOffset(m.X+xx, m.Y+yy)
				atlas.Pix[i+0] = 255
				atlas.Pix[i+1] = 255
				atlas.Pix[i+2] = 255
				atlas.Pix[i+3] = a
			}
		}
	}

	meta := atlasMeta{
		Ascent:   ascent,
		Descent:  descent,
		LineH:    lineH,
		Fallback: int(fallback),
		Glyphs:   placed,
	}
	return meta, atlas
}

func alphaBounds(img *image.Alpha) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	found := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.AlphaAt(x, y).A == 0 {
				continue
			}
			found = true
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	if !found {
		return image.Rectangle{}
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func charsetRunes(s string) []rune {
	n := norm.NFC.String(s)
	seen := map[rune]bool{}
	out := make([]rune, 0, len(n))
	for _, r := range strings.ReplaceAll(n, "\n", "") {
		if seen[r] {
			continue
		}
		seen[r] = true
		out = append(out, r)
	}
	return out
}

func containsRune(runes []rune, target rune) bool {
	for _, r := range runes {
		if r == target {
			return true
		}
	}
	return false
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}
