package relay

// File overview:
// draw renders relay geometry and anchors in world space for this part.
// Subsystem: part catalog (relay) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Bounds handles bounds.
func (self *Relay) Bounds() core.Rect {
	rows := len(self.Poles)
	if rows < 1 {
		rows = 1
	}
	height := float64(rows*20 + 24)
	return core.RectFromPoints(
		core.Pt{X: self.Pos.X - 28, Y: self.Pos.Y - height/2},
		core.Pt{X: self.Pos.X + 28, Y: self.Pos.Y + height/2},
	)
}

// Anchors handles anchors.
func (self *Relay) Anchors() []core.PinAnchor {
	self.ensureContactSlices()

	anchors := []core.PinAnchor{
		{Pt: core.Pt{X: self.Pos.X - 12, Y: self.Pos.Y + 24}, PinID: self.PinCoilA},
		{Pt: core.Pt{X: self.Pos.X + 12, Y: self.Pos.Y + 24}, PinID: self.PinCoilB},
	}

	startY := self.Pos.Y - float64((len(self.Poles)-1)*20)/2
	for i, pole := range self.Poles {
		y := startY + float64(i*20)
		anchors = append(anchors,
			core.PinAnchor{Pt: core.Pt{X: self.Pos.X - 28, Y: y}, PinID: pole.PinNC},
			core.PinAnchor{Pt: core.Pt{X: self.Pos.X, Y: y}, PinID: pole.PinCommon},
			core.PinAnchor{Pt: core.Pt{X: self.Pos.X + 28, Y: y}, PinID: pole.PinNO},
		)
	}

	return anchors
}

// HitTest handles hit test.
func (self *Relay) HitTest(pt core.Pt) part.HitResult {
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

// Draw draws its work.
func (self *Relay) Draw(ctx part.DrawContext) {
	self.asset().Draw(ctx, self.Bounds())
}
