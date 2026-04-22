package part

import (
	"coilforge/internal/core"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawVG* use symbol-centred SVG coordinates mapped with [SVGLocalToWorld].
func DrawVGLine(ctx DrawContext, base core.BasePart, x1, y1, x2, y2, strokeWidth float64, stroke color.RGBA) {
	ax, ay := vgPointToScreen(ctx, base, x1, y1)
	bx, by := vgPointToScreen(ctx, base, x2, y2)
	vector.StrokeLine(ctx.Dst, ax, ay, bx, by, vgStrokeWidth(ctx, strokeWidth), vgColor(ctx, stroke), false)
}

func DrawVGCubicBezier(ctx DrawContext, base core.BasePart, x0, y0, cx1, cy1, cx2, cy2, x3, y3, strokeWidth float64, stroke color.RGBA) {
	sx0, sy0 := vgPointToScreen(ctx, base, x0, y0)
	sx1, sy1 := vgPointToScreen(ctx, base, cx1, cy1)
	sx2, sy2 := vgPointToScreen(ctx, base, cx2, cy2)
	sx3, sy3 := vgPointToScreen(ctx, base, x3, y3)
	var path vector.Path
	path.MoveTo(float32(sx0), float32(sy0))
	path.CubicTo(float32(sx1), float32(sy1), float32(sx2), float32(sy2), float32(sx3), float32(sy3))
	sw := vgStrokeWidth(ctx, strokeWidth)
	so := &vector.StrokeOptions{
		Width:    sw,
		LineJoin: vector.LineJoinRound,
		LineCap:  vector.LineCapRound,
	}
	drawOp := &vector.DrawPathOptions{}
	drawOp.ColorScale.ScaleWithColor(vgColor(ctx, stroke))
	vector.StrokePath(ctx.Dst, &path, so, drawOp)
}

func DrawVGPolyline(ctx DrawContext, base core.BasePart, points []float64, closed bool, strokeWidth float64, stroke color.RGBA) {
	if len(points) < 4 {
		return
	}
	var path vector.Path
	for i := 0; i < len(points); i += 2 {
		sx, sy := vgPointToScreen(ctx, base, points[i], points[i+1])
		if i == 0 {
			path.MoveTo(sx, sy)
		} else {
			path.LineTo(sx, sy)
		}
	}
	if closed {
		path.Close()
	}
	sw := vgStrokeWidth(ctx, strokeWidth)
	so := &vector.StrokeOptions{
		Width:    sw,
		LineJoin: vector.LineJoinRound,
		LineCap:  vector.LineCapRound,
	}
	drawOp := &vector.DrawPathOptions{}
	drawOp.ColorScale.ScaleWithColor(vgColor(ctx, stroke))
	vector.StrokePath(ctx.Dst, &path, so, drawOp)
}

func DrawVGRoundedRect(ctx DrawContext, base core.BasePart, x, y, w, h, rx, ry float64, fill color.RGBA, hasFill bool, stroke color.RGBA, strokeWidth float64, hasStroke bool) {
	rx = math.Min(rx, w*0.5)
	ry = math.Min(ry, h*0.5)
	if rx <= 0 || ry <= 0 {
		DrawVGRect(ctx, base, x, y, w, h, fill, hasFill, stroke, strokeWidth, hasStroke)
		return
	}
	const segPerCorner = 12
	pts := roundedRectOutlinePoints(x, y, w, h, rx, ry, segPerCorner)
	if len(pts) < 6 {
		DrawVGRect(ctx, base, x, y, w, h, fill, hasFill, stroke, strokeWidth, hasStroke)
		return
	}
	var path vector.Path
	first := true
	for i := 0; i < len(pts); i += 2 {
		sx, sy := vgPointToScreen(ctx, base, pts[i], pts[i+1])
		if first {
			path.MoveTo(sx, sy)
			first = false
			continue
		}
		path.LineTo(sx, sy)
	}
	path.Close()
	drawOpFill := &vector.DrawPathOptions{}
	drawOpFill.ColorScale.ScaleWithColor(vgColor(ctx, fill))
	drawOpStroke := &vector.DrawPathOptions{}
	drawOpStroke.ColorScale.ScaleWithColor(vgColor(ctx, stroke))
	if hasFill {
		vector.FillPath(ctx.Dst, &path, &vector.FillOptions{FillRule: vector.FillRuleNonZero}, drawOpFill)
	}
	if hasStroke {
		sw := vgStrokeWidth(ctx, strokeWidth)
		so := &vector.StrokeOptions{
			Width:    sw,
			LineJoin: vector.LineJoinRound,
			LineCap:  vector.LineCapRound,
		}
		vector.StrokePath(ctx.Dst, &path, so, drawOpStroke)
	}
}

func roundedRectOutlinePoints(x, y, w, h, rx, ry float64, seg int) []float64 {
	rx = math.Min(rx, w*0.5)
	ry = math.Min(ry, h*0.5)
	if rx <= 0 || ry <= 0 {
		return nil
	}
	if seg < 2 {
		seg = 2
	}
	var pts []float64
	appendArc := func(cx, cy, a0, a1 float64) {
		for i := 0; i <= seg; i++ {
			t := float64(i) / float64(seg)
			a := a0 + (a1-a0)*t
			pts = append(pts, cx+rx*math.Cos(a), cy+ry*math.Sin(a))
		}
	}
	appendArc(x+rx, y+ry, math.Pi, 1.5*math.Pi)
	pts = append(pts, x+w-rx, y)
	appendArc(x+w-rx, y+ry, -math.Pi/2, 0)
	pts = append(pts, x+w, y+h-ry)
	appendArc(x+w-rx, y+h-ry, 0, 0.5*math.Pi)
	pts = append(pts, x+rx, y+h)
	appendArc(x+rx, y+h-ry, 0.5*math.Pi, math.Pi)
	pts = append(pts, x, y+ry)
	return pts
}

func DrawVGRect(ctx DrawContext, base core.BasePart, x, y, w, h float64, fill color.RGBA, hasFill bool, stroke color.RGBA, strokeWidth float64, hasStroke bool) {
	x0, y0 := vgPointToScreen(ctx, base, x, y)
	x1, y1 := vgPointToScreen(ctx, base, x+w, y+h)
	rectW := x1 - x0
	rectH := y1 - y0
	if hasFill {
		vector.FillRect(ctx.Dst, x0, y0, rectW, rectH, vgColor(ctx, fill), false)
	}
	if hasStroke {
		vector.StrokeRect(ctx.Dst, x0, y0, rectW, rectH, vgStrokeWidth(ctx, strokeWidth), vgColor(ctx, stroke), false)
	}
}

func DrawVGCircle(ctx DrawContext, base core.BasePart, cx, cy, r float64, fill color.RGBA, hasFill bool, stroke color.RGBA, strokeWidth float64, hasStroke bool) {
	scx, scy := vgPointToScreen(ctx, base, cx, cy)
	radiusPx := float32(r * SVGUserUnitToWorld * ctx.Zoom)
	if hasFill {
		vector.FillCircle(ctx.Dst, scx, scy, radiusPx, vgColor(ctx, fill), false)
	}
	if hasStroke {
		vector.StrokeCircle(ctx.Dst, scx, scy, radiusPx, vgStrokeWidth(ctx, strokeWidth), vgColor(ctx, stroke), false)
	}
}

func vgPointToScreen(ctx DrawContext, base core.BasePart, x, y float64) (float32, float32) {
	w := SVGLocalToWorld(base, x, y)
	sx, sy := ctx.WorldToScreen(w)
	return float32(sx), float32(sy)
}

func vgStrokeWidth(ctx DrawContext, strokeWidthSVG float64) float32 {
	px := strokeWidthSVG * SVGUserUnitToWorld * ctx.Zoom
	if px < 0.5 {
		px = 0.5
	}
	return float32(px)
}

func vgColor(ctx DrawContext, c color.RGBA) color.RGBA {
	src := c
	if ctx.DarkMode {
		if v, ok := darkModeSwapBlackWhite(c); ok {
			src = v
		}
	}
	if !ctx.Ghost {
		return src
	}
	tint := color.RGBA{R: 255, G: 255, B: 255, A: 144}
	// Blend toward a light ghost tint instead of multiplying channels.
	mix := float64(tint.A) / 255.0
	inv := 1.0 - mix
	return color.RGBA{
		R: uint8(float64(src.R)*inv + float64(tint.R)*mix),
		G: uint8(float64(src.G)*inv + float64(tint.G)*mix),
		B: uint8(float64(src.B)*inv + float64(tint.B)*mix),
		A: uint8(math.Max(24, float64((uint16(src.A)*uint16(tint.A))/255))),
	}
}

// darkModeSwapBlackWhite swaps near-black and near-white in dark mode only; all other colors unchanged.
// Lets symbols with a white fill and black strokes (e.g. indicator-off) read correctly: black→white on a white→black disc.
func darkModeSwapBlackWhite(c color.RGBA) (color.RGBA, bool) {
	const thr uint8 = 28
	hi := uint8(255 - thr)
	if c.R <= thr && c.G <= thr && c.B <= thr {
		return color.RGBA{R: 255, G: 255, B: 255, A: c.A}, true
	}
	if c.R >= hi && c.G >= hi && c.B >= hi {
		return color.RGBA{R: 0, G: 0, B: 0, A: c.A}, true
	}
	return color.RGBA{}, false
}
