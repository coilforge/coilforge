package relay

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (r *Relay) Bounds() core.Rect {
	rows := len(r.Poles)
	if rows < 1 {
		rows = 1
	}
	height := float64(rows*20 + 24)
	return core.RectFromPoints(
		core.Pt{X: r.Pos.X - 28, Y: r.Pos.Y - height/2},
		core.Pt{X: r.Pos.X + 28, Y: r.Pos.Y + height/2},
	)
}

func (r *Relay) Anchors() []core.PinAnchor {
	r.ensureContactSlices()

	anchors := []core.PinAnchor{
		{Pt: core.Pt{X: r.Pos.X - 12, Y: r.Pos.Y + 24}, PinID: r.PinCoilA},
		{Pt: core.Pt{X: r.Pos.X + 12, Y: r.Pos.Y + 24}, PinID: r.PinCoilB},
	}

	startY := r.Pos.Y - float64((len(r.Poles)-1)*20)/2
	for i, pole := range r.Poles {
		y := startY + float64(i*20)
		anchors = append(anchors,
			core.PinAnchor{Pt: core.Pt{X: r.Pos.X - 28, Y: y}, PinID: pole.PinNC},
			core.PinAnchor{Pt: core.Pt{X: r.Pos.X, Y: y}, PinID: pole.PinCommon},
			core.PinAnchor{Pt: core.Pt{X: r.Pos.X + 28, Y: y}, PinID: pole.PinNO},
		)
	}

	return anchors
}

func (r *Relay) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, r.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (r *Relay) Draw(ctx part.DrawContext) {
	r.asset().Draw(ctx, r.Bounds())
}
