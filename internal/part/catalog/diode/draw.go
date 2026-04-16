package diode

// File overview:
// draw renders diode geometry and anchors in world space for this part.
// Subsystem: part catalog (diode) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (d *Diode) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: d.Pos.X - 18, Y: d.Pos.Y - 8},
		core.Pt{X: d.Pos.X + 18, Y: d.Pos.Y + 8},
	)
}

// Anchors handles anchors.
func (d *Diode) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: d.Pos.X - 20, Y: d.Pos.Y}, PinID: d.PinAnode},
		{Pt: core.Pt{X: d.Pos.X + 20, Y: d.Pos.Y}, PinID: d.PinCathode},
	}
}

// HitTest handles hit test.
func (d *Diode) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, d.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (d *Diode) Draw(ctx part.DrawContext) {
	d.asset().Draw(ctx, d.Bounds())
}
