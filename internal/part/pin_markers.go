package part

import (
	"coilforge/internal/core"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PinMarkerHalfPx is half the square edge length at zoom=1, in screen pixels.
const PinMarkerHalfPx float32 = 0.75

// DrawPinMarkers draws small filled squares at each anchor (pin connection point).
func DrawPinMarkers(ctx DrawContext, anchors []core.PinAnchor) {
	if len(anchors) == 0 || ctx.Dst == nil {
		return
	}
	fill := pinMarkerFill(ctx)
	// Scale markers with zoom so thick strokes don't swallow them when zooming in.
	scale := float32(math.Max(0.25, ctx.Zoom))
	h := PinMarkerHalfPx * scale
	w := h * 2
	for _, a := range anchors {
		sx, sy := ctx.WorldToScreen(a.Pt)
		x := float32(sx) - h
		y := float32(sy) - h
		vector.FillRect(ctx.Dst, x, y, w, w, fill, true)
	}
}

func pinMarkerFill(ctx DrawContext) color.RGBA {
	var c color.RGBA
	if ctx.DarkMode {
		c = color.RGBA{R: 232, G: 236, B: 244, A: 255}
	} else {
		c = color.RGBA{R: 48, G: 52, B: 62, A: 255}
	}
	if ctx.Ghost {
		c.A = 140
	}
	return c
}
