package skeleton

// File overview:
// draw renders skeleton geometry and anchors in world space for this part.
// Subsystem: part catalog (skeleton) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (self *Template) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 16, Y: self.Pos.Y - 8},
		core.Pt{X: self.Pos.X + 16, Y: self.Pos.Y + 8},
	)
}

// Anchors handles anchors.
func (self *Template) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: self.Pos.X - 18, Y: self.Pos.Y}, PinID: self.PinA},
		{Pt: core.Pt{X: self.Pos.X + 18, Y: self.Pos.Y}, PinID: self.PinB},
	}
}

// HitTest handles hit test.
func (self *Template) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Template) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
