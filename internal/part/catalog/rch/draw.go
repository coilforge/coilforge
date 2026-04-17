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
func (self *RCH) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 22, Y: self.Pos.Y - 10},
		core.Pt{X: self.Pos.X + 22, Y: self.Pos.Y + 10},
	)
}

// Anchors handles anchors.
func (self *RCH) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: self.Pos.X - 24, Y: self.Pos.Y}, PinID: self.PinIn},
		{Pt: core.Pt{X: self.Pos.X + 24, Y: self.Pos.Y}, PinID: self.PinOut},
	}
}

// HitTest handles hit test.
func (self *RCH) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *RCH) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
