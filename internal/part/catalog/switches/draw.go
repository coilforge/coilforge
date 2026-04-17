package switches

// File overview:
// draw renders switches geometry and anchors in world space for this part.
// Subsystem: part catalog (switches) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (self *Switch) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 18, Y: self.Pos.Y - 10},
		core.Pt{X: self.Pos.X + 18, Y: self.Pos.Y + 10},
	)
}

// Anchors handles anchors.
func (self *Switch) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: self.Pos.X - 20, Y: self.Pos.Y}, PinID: self.PinA},
		{Pt: core.Pt{X: self.Pos.X + 20, Y: self.Pos.Y}, PinID: self.PinB},
	}
}

// HitTest handles hit test.
func (self *Switch) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Switch) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
