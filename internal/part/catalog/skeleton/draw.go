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
func (t *Template) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: t.Pos.X - 16, Y: t.Pos.Y - 8},
		core.Pt{X: t.Pos.X + 16, Y: t.Pos.Y + 8},
	)
}

// Anchors handles anchors.
func (t *Template) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: t.Pos.X - 18, Y: t.Pos.Y}, PinID: t.PinA},
		{Pt: core.Pt{X: t.Pos.X + 18, Y: t.Pos.Y}, PinID: t.PinB},
	}
}

// HitTest handles hit test.
func (t *Template) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, t.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (t *Template) Draw(ctx part.DrawContext) {
	t.asset().Draw(ctx, t.Bounds())
}
