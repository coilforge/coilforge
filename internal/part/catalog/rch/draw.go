package rch

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (r *RCH) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: r.Pos.X - 22, Y: r.Pos.Y - 10},
		core.Pt{X: r.Pos.X + 22, Y: r.Pos.Y + 10},
	)
}

func (r *RCH) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: r.Pos.X - 24, Y: r.Pos.Y}, PinID: r.PinIn},
		{Pt: core.Pt{X: r.Pos.X + 24, Y: r.Pos.Y}, PinID: r.PinOut},
	}
}

func (r *RCH) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, r.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (r *RCH) Draw(ctx part.DrawContext) {
	r.asset().Draw(ctx, r.Bounds())
}
