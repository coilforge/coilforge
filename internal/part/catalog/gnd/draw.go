package gnd

// File overview:
// draw renders gnd geometry and anchors in world space for this part.
// Subsystem: part catalog (gnd) drawing.
// It cooperates with assets selection and is called through the generic part.Draw path.
// Flow position: part-level render leaf beneath render scene orchestration.

import (
	"fmt"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Gnd) layoutName() string {
	slots := RotationSlots
	if slots <= 0 {
		slots = 4
	}
	r := self.Rotation % slots
	if r < 0 {
		r += slots
	}
	return fmt.Sprintf("gnd-%d", r)
}

// drawBase maps SVG through position and mirror only; rotation is baked into layoutName suffix.
func (self *Gnd) drawBase() core.BasePart {
	b := self.BasePart
	b.Rotation = 0
	return b
}

// Bounds handles bounds.
func (self *Gnd) Bounds() core.Rect {
	if r, ok := part.HitBoundsFromVectorLayout(self.layoutName(), self.drawBase()); ok {
		return r
	}
	return core.Rect{}
}

// Anchors handles anchors.
func (self *Gnd) Anchors() []core.PinAnchor {
	return part.AnchorsFromVectorMarkerIDs(self.layoutName(), self.drawBase(), gndPinMarkerMap(self))
}

// HitTest handles hit test.
func (self *Gnd) HitTest(pt core.Pt) part.HitResult {
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
func (self *Gnd) Draw(ctx part.DrawContext) {
	part.VectorAsset{Name: self.layoutName()}.Draw(ctx, self.drawBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}
