package wire

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (w *Wire) Bounds() core.Rect {
	seg := w.Segments()[0]
	return core.RectFromPoints(seg.A, seg.B)
}

func (w *Wire) Anchors() []core.PinAnchor {
	seg := w.Segments()[0]
	return []core.PinAnchor{
		{Pt: seg.A, PinID: w.PinA},
		{Pt: seg.B, PinID: w.PinB},
	}
}

func (w *Wire) HitTest(pt core.Pt) part.HitResult {
	seg := w.Segments()[0]
	if core.PointNearSeg(pt, seg, 6) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (w *Wire) Segments() []core.Seg {
	return []core.Seg{{
		A: core.Pt{X: w.Pos.X - w.Half.X, Y: w.Pos.Y - w.Half.Y},
		B: core.Pt{X: w.Pos.X + w.Half.X, Y: w.Pos.Y + w.Half.Y},
	}}
}

func (w *Wire) Draw(ctx part.DrawContext) {
	w.asset().Draw(ctx, w.Bounds())
}
