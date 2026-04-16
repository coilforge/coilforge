package wire

// File overview:
// draw renders wire geometry and anchors in world space for this part.
// Subsystem: part catalog (wire) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (w *Wire) Bounds() core.Rect {
	seg := w.Segments()[0]
	return core.RectFromPoints(seg.A, seg.B)
}

// Anchors handles anchors.
func (w *Wire) Anchors() []core.PinAnchor {
	seg := w.Segments()[0]
	return []core.PinAnchor{
		{Pt: seg.A, PinID: w.PinA},
		{Pt: seg.B, PinID: w.PinB},
	}
}

// HitTest handles hit test.
func (w *Wire) HitTest(pt core.Pt) part.HitResult {
	seg := w.Segments()[0]
	if core.PointNearSeg(pt, seg, 6) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Segments handles segments.
func (w *Wire) Segments() []core.Seg {
	return []core.Seg{{
		A: core.Pt{X: w.Pos.X - w.Half.X, Y: w.Pos.Y - w.Half.Y},
		B: core.Pt{X: w.Pos.X + w.Half.X, Y: w.Pos.Y + w.Half.Y},
	}}
}

// Draw draws its work.
func (w *Wire) Draw(ctx part.DrawContext) {
	w.asset().Draw(ctx, w.Bounds())
}
