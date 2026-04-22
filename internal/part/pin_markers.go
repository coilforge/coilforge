package part

import (
	"coilforge/internal/core"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PinMarkerRadiusPx is the marker radius at zoom=1, in screen pixels.
// Increased by 50% from the previous square half-size for better visibility.
const PinMarkerRadiusPx float32 = 1.125

// DrawPinMarkers draws small filled circles at each anchor (pin connection point).
func DrawPinMarkers(ctx DrawContext, anchors []core.PinAnchor) {
	if len(anchors) == 0 || ctx.Dst == nil {
		return
	}
	fill := pinMarkerFill(ctx)
	// Scale markers with zoom so thick strokes don't swallow them when zooming in.
	scale := float32(math.Max(0.25, ctx.Zoom))
	r := PinMarkerRadiusPx * scale
	for _, a := range anchors {
		sx, sy := ctx.WorldToScreen(a.Pt)
		vector.FillCircle(ctx.Dst, float32(sx), float32(sy), r, fill, true)
	}
}

func pinMarkerFill(ctx DrawContext) color.RGBA {
	var c color.RGBA
	if ctx.DarkMode {
		c = color.RGBA{R: 122, G: 185, B: 255, A: 255}
	} else {
		c = color.RGBA{R: 28, G: 112, B: 220, A: 255}
	}
	if ctx.Ghost {
		c.A = 140
	}
	return c
}
