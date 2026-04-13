package switches

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (s *Switch) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: s.Pos.X - 18, Y: s.Pos.Y - 10},
		core.Pt{X: s.Pos.X + 18, Y: s.Pos.Y + 10},
	)
}

func (s *Switch) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: s.Pos.X - 20, Y: s.Pos.Y}, PinID: s.PinA},
		{Pt: core.Pt{X: s.Pos.X + 20, Y: s.Pos.Y}, PinID: s.PinB},
	}
}

func (s *Switch) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, s.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (s *Switch) Draw(ctx part.DrawContext) {
	s.asset().Draw(ctx, s.Bounds())
}
