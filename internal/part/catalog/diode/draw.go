package diode

import (
	"fmt"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Diode) layoutName() string {
	return fmt.Sprintf("diode-%d", normalizeRotation(self.Rotation))
}

func (self *Diode) drawBase() core.BasePart {
	b := self.BasePart
	b.Rotation = 0
	return b
}

func (self *Diode) Bounds() core.Rect {
	if r, ok := part.HitBoundsFromVectorLayout(self.layoutName(), self.drawBase()); ok {
		return r
	}
	return core.Rect{}
}

func (self *Diode) Anchors() []core.PinAnchor {
	return part.AnchorsFromVectorMarkerIDs(self.layoutName(), self.drawBase(), diodePinMarkerMap(self))
}

func (self *Diode) HitTest(pt core.Pt) part.HitResult {
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

func (self *Diode) Draw(ctx part.DrawContext) {
	part.VectorAsset{Name: self.layoutName()}.Draw(ctx, self.drawBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}
