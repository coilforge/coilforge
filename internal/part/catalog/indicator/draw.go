package indicator

// File overview:
// draw renders indicator geometry and anchors in world space for this part.
// Subsystem: part catalog (indicator) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (ind *Indicator) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: ind.Pos.X - 12, Y: ind.Pos.Y - 12},
		core.Pt{X: ind.Pos.X + 12, Y: ind.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (ind *Indicator) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: ind.Pos.X, Y: ind.Pos.Y + 16},
		PinID: ind.PinA,
	}}
}

// HitTest handles hit test.
func (ind *Indicator) HitTest(pt core.Pt) part.HitResult {
	for _, anchor := range ind.Anchors() {
		if core.PointNearSeg(pt, core.Seg{A: anchor.Pt, B: anchor.Pt}, 6) {
			return part.HitResult{Hit: true, Kind: part.HitPin, PinID: anchor.PinID}
		}
	}
	if core.PointInRect(pt, ind.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (ind *Indicator) Draw(ctx part.DrawContext) {
	ind.asset().Draw(ctx, ind.Bounds())
}
