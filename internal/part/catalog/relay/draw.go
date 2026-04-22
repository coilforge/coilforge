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
	relayTopAboveMidMajors  = 3
)

func (self *Relay) relayCoilBase() core.BasePart {
	return self.drawBase()
}

func (self *Relay) relayMidBase() core.BasePart {
	b := self.drawBase()
	b.Pos.Y -= float64(relayMidAboveCoilMajors) * majorGridWorld
	return b
}

func (self *Relay) relayTopBase() core.BasePart {
	b := self.drawBase()
	b.Pos.Y -= float64(relayMidAboveCoilMajors+relayTopAboveMidMajors) * majorGridWorld
	return b
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
		{self.layoutName(self.midStem()), self.relayMidBase()},
		{self.layoutName(self.coilStem()), self.relayCoilBase()},
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
	out = append(out, part.AnchorsFromVectorMarkerIDs(self.layoutName(self.midStem()), self.relayMidBase(), relayMidPinMarkerMap(self))...)
	out = append(out, part.AnchorsFromVectorMarkerIDs(self.layoutName(self.coilStem()), self.relayCoilBase(), relayCoilPinMarkerMap(self))...)
	return out
}

func relayMidPinMarkerMap(self *Relay) map[string]core.PinID {
	return map[string]core.PinID{
		"COM": self.COM,
		"NC":  self.NC,
		"NO":  self.NO,
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
	part.VectorAsset{Name: self.layoutName(self.midStem())}.Draw(ctx, self.relayMidBase())
	part.VectorAsset{Name: self.layoutName("top")}.Draw(ctx, self.relayTopBase())
	part.DrawPinMarkers(ctx, self.Anchors())
}
