package switches

import (
	"fmt"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Switches) layoutName() string {
	stem := "switch-off"
	if self.On {
		stem = "switch-on"
	}
	r := normalizeRotation(self.Rotation)
	return fmt.Sprintf("%s-%d", stem, r)
}

func (self *Switches) drawBase() core.BasePart {
	b := self.BasePart
	b.Rotation = 0
	return b
}

func (self *Switches) Bounds() core.Rect {
	base := self.drawBase()
	rOn, okOn := part.HitBoundsFromVectorLayout(fmt.Sprintf("switch-on-%d", normalizeRotation(self.Rotation)), base)
	rOff, okOff := part.HitBoundsFromVectorLayout(fmt.Sprintf("switch-off-%d", normalizeRotation(self.Rotation)), base)
	if !okOn && !okOff {
		return core.Rect{}
	}
	if !okOn {
		return rOff
	}
	if !okOff {
		return rOn
	}
	return core.Rect{
		Min: core.Pt{
			X: min(rOn.Min.X, rOff.Min.X),
			Y: min(rOn.Min.Y, rOff.Min.Y),
		},
		Max: core.Pt{
			X: max(rOn.Max.X, rOff.Max.X),
			Y: max(rOn.Max.Y, rOff.Max.Y),
		},
	}
}

func (self *Switches) Anchors() []core.PinAnchor {
	return part.AnchorsFromVectorMarkerIDs(self.layoutName(), self.drawBase(), switchesPinMarkerMap(self))
}

func (self *Switches) HitTest(pt core.Pt) part.HitResult {
	for _, anchor := range self.Anchors() {
		if core.PointNearSeg(pt, core.Seg{A: anchor.Pt, B: anchor.Pt}, 6) {
			return part.HitResult{Hit: true, Kind: part.HitPin, PinID: anchor.PinID}
		}
	}
	if core.PointInRect(pt, self.Bounds()) {
		return part.HitResult{Hit: true, Kind: part.HitBody}
	}
	return part.HitResult{}
}

func (self *Switches) Draw(ctx part.DrawContext) {
	part.VectorAsset{Name: self.layoutName()}.Draw(ctx, self.drawBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
