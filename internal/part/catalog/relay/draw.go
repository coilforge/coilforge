package relay

import (
	"fmt"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Relay) layoutName(stem string) string {
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

func (self *Relay) drawBase() core.BasePart {
	b := self.BasePart
	b.Rotation = 0
	return b
}

func (self *Relay) coilStem() string {
	if self.Energized {
		return "coil-on"
	}
	return "coil-off"
}

func (self *Relay) midStem() string {
	if self.Energized {
		return "mid-on"
	}
	return "mid-off"
}

// Vertical stack offsets (world Y, increasing Y = down on screen): mid sits above coil, top above mid.
// majorGridWorld matches [world.MajorGridWorld]; catalog cannot import world per package boundaries.
const majorGridWorld = 4.0

const (
	relayMidAboveCoilMajors = 6
	relayPolePitchMajors    = 6
	relayTopAboveMidMajors  = 3
)

func (self *Relay) relayCoilBase() core.BasePart {
	return self.drawBase()
}

func (self *Relay) relayMidBaseForPole(pole int) core.BasePart {
	b := self.drawBase()
	if pole < 1 {
		pole = 1
	}
	midMajors := relayMidAboveCoilMajors + (pole-1)*relayPolePitchMajors
	b.Pos.Y -= float64(midMajors) * majorGridWorld
	// Mid flip is a left/right swap of contact geometry only (COM side changes edge).
	if self.isPoleFlipped(pole) {
		b.Mirror = !b.Mirror
	}
	return b
}

func (self *Relay) relayTopBase() core.BasePart {
	b := self.drawBase()
	topMajors := relayMidAboveCoilMajors + (self.poleCountClamped()-1)*relayPolePitchMajors + relayTopAboveMidMajors
	b.Pos.Y -= float64(topMajors) * majorGridWorld
	return b
}

func (self *Relay) poleCountClamped() int {
	return clampRelayPoleCount(self.PoleCount)
}

func (self *Relay) isPoleFlipped(pole int) bool {
	if pole < 1 || pole > 8 {
		return false
	}
	return (self.MidFlipMask & uint8(1<<(pole-1))) != 0
}

func (self *Relay) contactPinsForPole(pole int) (com, nc, no core.PinID) {
	switch pole {
	case 1:
		return self.COM1, self.NC1, self.NO1
	case 2:
		return self.COM2, self.NC2, self.NO2
	case 3:
		return self.COM3, self.NC3, self.NO3
	case 4:
		return self.COM4, self.NC4, self.NO4
	case 5:
		return self.COM5, self.NC5, self.NO5
	case 6:
		return self.COM6, self.NC6, self.NO6
	case 7:
		return self.COM7, self.NC7, self.NO7
	case 8:
		return self.COM8, self.NC8, self.NO8
	default:
		return 0, 0, 0
	}
}

func rectUnion(a, b core.Rect) core.Rect {
	return core.Rect{
		Min: core.Pt{X: min(a.Min.X, b.Min.X), Y: min(a.Min.Y, b.Min.Y)},
		Max: core.Pt{X: max(a.Max.X, b.Max.X), Y: max(a.Max.Y, b.Max.Y)},
	}
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

func (self *Relay) Bounds() core.Rect {
	type namedBase struct {
		name string
		base core.BasePart
	}
	layers := []namedBase{
		{self.layoutName("top"), self.relayTopBase()},
		{self.layoutName(self.coilStem()), self.relayCoilBase()},
	}
	for pole := 1; pole <= self.poleCountClamped(); pole++ {
		layers = append(layers, namedBase{
			name: self.layoutName(self.midStem()),
			base: self.relayMidBaseForPole(pole),
		})
	}
	var out core.Rect
	okAny := false
	for _, layer := range layers {
		r, ok := part.HitBoundsFromVectorLayout(layer.name, layer.base)
		if !ok {
			continue
		}
		if !okAny {
			out = r
			okAny = true
			continue
		}
		out = rectUnion(out, r)
	}
	if !okAny {
		return core.Rect{}
	}
	return out
}

func (self *Relay) Anchors() []core.PinAnchor {
	var out []core.PinAnchor
	for pole := 1; pole <= self.poleCountClamped(); pole++ {
		out = append(out, part.AnchorsFromVectorMarkerIDs(
			self.layoutName(self.midStem()),
			self.relayMidBaseForPole(pole),
			relayMidPinMarkerMap(self, pole),
		)...)
	}
	out = append(out, part.AnchorsFromVectorMarkerIDs(self.layoutName(self.coilStem()), self.relayCoilBase(), relayCoilPinMarkerMap(self))...)
	return out
}

func relayMidPinMarkerMap(self *Relay, pole int) map[string]core.PinID {
	com, nc, no := self.contactPinsForPole(pole)
	return map[string]core.PinID{
		"COM": com,
		"NC":  nc,
		"NO":  no,
	}
}

func relayCoilPinMarkerMap(self *Relay) map[string]core.PinID {
	return map[string]core.PinID{
		"CoilA": self.CoilA,
		"CoilB": self.CoilB,
	}
}

func (self *Relay) HitTest(pt core.Pt) part.HitResult {
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

func (self *Relay) Draw(ctx part.DrawContext) {
	part.VectorAsset{Name: self.layoutName(self.coilStem())}.Draw(ctx, self.relayCoilBase())
	for pole := 1; pole <= self.poleCountClamped(); pole++ {
		part.VectorAsset{Name: self.layoutName(self.midStem())}.Draw(ctx, self.relayMidBaseForPole(pole))
	}
	part.VectorAsset{Name: self.layoutName("top")}.Draw(ctx, self.relayTopBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}

func (self *Relay) poleAtPoint(pt core.Pt) int {
	for pole := 1; pole <= self.poleCountClamped(); pole++ {
		r, ok := part.HitBoundsFromVectorLayout(self.layoutName(self.midStem()), self.relayMidBaseForPole(pole))
		if ok && core.PointInRect(pt, r) {
			return pole
		}
	}
	return 0
}

// ToggleMidFlipAt toggles the clicked pole's COM-side orientation.
func (self *Relay) ToggleMidFlipAt(pt core.Pt) bool {
	pole := self.poleAtPoint(pt)
	if pole == 0 {
		return false
	}
	return self.TogglePoleFlip(pole)
}
