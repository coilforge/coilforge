package indicator

// File overview:
// draw renders indicator geometry and anchors in world space for this part.
// Subsystem: part catalog (indicator) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"fmt"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Indicator) layoutName() string {
	stem := "indicator-off"
	if self.Lit {
		stem = "indicator-on"
	}
	slots := RotationSlots
	if slots <= 0 {
		slots = 4
	}
	r := self.Rotation % slots
	if r < 0 {
		r += slots
	}
	return fmt.Sprintf("%s-%d", stem, r)
}

// drawBase maps SVG through position and mirror only; rotation is baked into layoutName suffix.
func (self *Indicator) drawBase() core.BasePart {
	b := self.BasePart
	b.Rotation = 0
	return b
}

// Bounds handles bounds.
func (self *Indicator) Bounds() core.Rect {
	if r, ok := part.HitBoundsFromVectorLayout(self.layoutName(), self.drawBase()); ok {
		return r
	}
	return core.Rect{}
}

// Anchors handles anchors.
func (self *Indicator) Anchors() []core.PinAnchor {
	return part.AnchorsFromVectorMarkerIDs(self.layoutName(), self.drawBase(), map[string]core.PinID{
		"PinA": self.PinA,
		"PinB": self.PinB,
	})
}

// HitTest handles hit test.
func (self *Indicator) HitTest(pt core.Pt) part.HitResult {
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

// Draw draws its work.
func (self *Indicator) Draw(ctx part.DrawContext) {
	part.VectorAsset{Name: self.layoutName()}.Draw(ctx, self.drawBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}
