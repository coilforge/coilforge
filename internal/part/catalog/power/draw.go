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
func (p *Power) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: p.Pos.X - 10, Y: p.Pos.Y - 12},
		core.Pt{X: p.Pos.X + 10, Y: p.Pos.Y + 12},
	)
}

// Anchors handles anchors.
func (p *Power) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: p.Pos.X, Y: p.Pos.Y + 16},
		PinID: p.Pin,
	}}
}

// HitTest handles hit test.
func (p *Power) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, p.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (p *Power) Draw(ctx part.DrawContext) {
	p.asset().Draw(ctx, p.Bounds())
}
