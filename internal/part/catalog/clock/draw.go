package clock

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (c *Clock) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: c.Pos.X - 18, Y: c.Pos.Y - 12},
		core.Pt{X: c.Pos.X + 18, Y: c.Pos.Y + 12},
	)
}

func (c *Clock) Anchors() []core.PinAnchor {
	return []core.PinAnchor{{
		Pt:    core.Pt{X: c.Pos.X + 20, Y: c.Pos.Y},
		PinID: c.PinOut,
	}}
}

func (c *Clock) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, c.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (c *Clock) Draw(ctx part.DrawContext) {
	c.asset().Draw(ctx, c.Bounds())
}
