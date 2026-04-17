package power

// File overview:
// draw renders power geometry and anchors in world space for this part.
// Subsystem: part catalog (power) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (self *Power) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 10, Y: self.Pos.Y - 12},
		core.Pt{X: self.Pos.X + 10, Y: self.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (self *Power) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: self.Pos.X, Y: self.Pos.Y + 16},
		PinID: self.Pin,
	}}
}

// HitTest handles hit test.
func (self *Power) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Power) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
