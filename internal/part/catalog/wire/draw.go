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
func (self *Wire) Bounds() core.Rect {
	seg := self.Segments()[0]
	return core.RectFromPoints(seg.A, seg.B)
}

// Anchors handles anchors.
func (self *Wire) Anchors() []core.PinAnchor {
	seg := self.Segments()[0]
	return []core.PinAnchor{
		{Pt: seg.A, PinID: self.PinA},
		{Pt: seg.B, PinID: self.PinB},
	}
}

// HitTest handles hit test.
func (self *Wire) HitTest(pt core.Pt) part.HitResult {
	seg := self.Segments()[0]
	if core.PointNearSeg(pt, seg, 6) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Segments handles segments.
func (self *Wire) Segments() []core.Seg {
	return []core.Seg{{
		A: core.Pt{X: self.Pos.X - self.Half.X, Y: self.Pos.Y - self.Half.Y},
		B: core.Pt{X: self.Pos.X + self.Half.X, Y: self.Pos.Y + self.Half.Y},
	}}
}

// Draw draws its work.
func (self *Wire) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
