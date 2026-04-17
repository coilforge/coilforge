package clock

// File overview:
// draw renders clock geometry and anchors in world space for this part.
// Subsystem: part catalog (clock) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (self *Clock) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 18, Y: self.Pos.Y - 12},
		core.Pt{X: self.Pos.X + 18, Y: self.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (self *Clock) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: self.Pos.X + 20, Y: self.Pos.Y},
		PinID: self.PinOut,
	}}
}

// HitTest handles hit test.
func (self *Clock) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Clock) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
