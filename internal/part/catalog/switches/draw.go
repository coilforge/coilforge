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
func (s *Switch) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: s.Pos.X - 18, Y: s.Pos.Y - 10},
		core.Pt{X: s.Pos.X + 18, Y: s.Pos.Y + 10},
	)
}

// Anchors handles anchors.
func (s *Switch) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: s.Pos.X - 20, Y: s.Pos.Y}, PinID: s.PinA},
		{Pt: core.Pt{X: s.Pos.X + 20, Y: s.Pos.Y}, PinID: s.PinB},
	}
}

// HitTest handles hit test.
func (s *Switch) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, s.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (s *Switch) Draw(ctx part.DrawContext) {
	s.asset().Draw(ctx, s.Bounds())
}
