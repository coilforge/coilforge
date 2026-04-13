package diode

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (d *Diode) Bounds() core.Rect {
	return core.RectFromPoints(
		core.Pt{X: d.Pos.X - 18, Y: d.Pos.Y - 8},
		core.Pt{X: d.Pos.X + 18, Y: d.Pos.Y + 8},
	)
}

func (d *Diode) Anchors() []core.PinAnchor {
	return []core.PinAnchor{
		{Pt: core.Pt{X: d.Pos.X - 20, Y: d.Pos.Y}, PinID: d.PinAnode},
		{Pt: core.Pt{X: d.Pos.X + 20, Y: d.Pos.Y}, PinID: d.PinCathode},
	}
}

func (d *Diode) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, d.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (d *Diode) Draw(ctx part.DrawContext) {
	d.asset().Draw(ctx, d.Bounds())
}
