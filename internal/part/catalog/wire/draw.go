package wire

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/render"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Pick tolerance roughly matches on-screen stroke thickness at typical zoom.
const hitTolWorld = 6.0

// DrawOrthogonalPolyline strokes a polyline in world space with schematic stroke width (ghost preview & wires).
// Each segment is stroked with round caps so ends overlap at bends (butt caps leave hollow corners when segments are drawn separately).
func DrawOrthogonalPolyline(ctx part.DrawContext, pts []core.Pt, clr color.Color) {
	if len(pts) < 2 {
		return
	}
	sw := part.SchematicStrokeScreenPx(ctx.Zoom)
	strokeOp := &vector.StrokeOptions{
		Width:   sw,
		LineCap: vector.LineCapRound,
	}
	drawOp := &vector.DrawPathOptions{
		AntiAlias: true,
	}
	drawOp.ColorScale.ScaleWithColor(clr)

	for i := 0; i < len(pts)-1; i++ {
		x0, y0 := ctx.WorldToScreen(pts[i])
		x1, y1 := ctx.WorldToScreen(pts[i+1])
		var path vector.Path
		path.MoveTo(float32(x0), float32(y0))
		path.LineTo(float32(x1), float32(y1))
		vector.StrokePath(ctx.Dst, &path, strokeOp, drawOp)
	}
}

// Draw renders the wire polyline using net color when sim context provides levels.
func (self *Wire) Draw(ctx part.DrawContext) {
	state := core.NetFloat
	if ctx.NetState != nil {
		state = ctx.NetState(self.PinA)
	}
	DrawOrthogonalPolyline(ctx, self.Points, render.WireColor(state))
}

// HitTest treats the polyline as a thin stroke in world space.
func (self *Wire) HitTest(pt core.Pt) part.HitResult {
	for _, s := range self.Segments() {
		if core.PointNearSeg(pt, s, hitTolWorld) {
			return part.HitResult{Hit: true, Kind: part.HitBody}
		}
	}
	return part.HitResult{}
}
