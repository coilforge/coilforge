package rch

// File overview:
// draw renders rch geometry and anchors in world space for this part.
// Subsystem: part catalog (rch) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (r *RCH) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: r.Pos.X - 22, Y: r.Pos.Y - 10},
		core.Pt{X: r.Pos.X + 22, Y: r.Pos.Y + 10},
	)
}

// Anchors handles anchors.
func (r *RCH) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: r.Pos.X - 24, Y: r.Pos.Y}, PinID: r.PinIn},
		{Pt: core.Pt{X: r.Pos.X + 24, Y: r.Pos.Y}, PinID: r.PinOut},
	}
}

// HitTest handles hit test.
func (r *RCH) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, r.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (r *RCH) Draw(ctx part.DrawContext) {
	r.asset().Draw(ctx, r.Bounds())
}
