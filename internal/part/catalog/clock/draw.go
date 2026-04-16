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
func (c *Clock) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: c.Pos.X - 18, Y: c.Pos.Y - 12},
		core.Pt{X: c.Pos.X + 18, Y: c.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (c *Clock) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: c.Pos.X + 20, Y: c.Pos.Y},
		PinID: c.PinOut,
	}}
}

// HitTest handles hit test.
func (c *Clock) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, c.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (c *Clock) Draw(ctx part.DrawContext) {
	c.asset().Draw(ctx, c.Bounds())
}
