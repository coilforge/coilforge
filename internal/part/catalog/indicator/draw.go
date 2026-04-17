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
func (self *Indicator) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 12, Y: self.Pos.Y - 12},
		core.Pt{X: self.Pos.X + 12, Y: self.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (self *Indicator) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: self.Pos.X, Y: self.Pos.Y + 16},
		PinID: self.PinA,
	}}
}

// HitTest handles hit test.
func (self *Indicator) HitTest(pt core.Pt) part.HitResult {
	for _, anchor := range self.Anchors() {
		if core.PointNearSeg(pt, core.Seg{A: anchor.Pt, B: anchor.Pt}, 6) {
			return part.HitResult{Hit: true, Kind: part.HitPin, PinID: anchor.PinID}
		}
	}
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Indicator) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
