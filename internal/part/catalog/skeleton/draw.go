package skeleton

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (t *Template) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: t.Pos.X - 16, Y: t.Pos.Y - 8},
		core.Pt{X: t.Pos.X + 16, Y: t.Pos.Y + 8},
	)
}

func (t *Template) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: t.Pos.X - 18, Y: t.Pos.Y}, PinID: t.PinA},
		{Pt: core.Pt{X: t.Pos.X + 18, Y: t.Pos.Y}, PinID: t.PinB},
	}
}

func (t *Template) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, t.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (t *Template) Draw(ctx part.DrawContext) {
	t.asset().Draw(ctx, t.Bounds())
}
