package flatten

// File overview:
// flatten converts placed parts and wires into connectivity data for net solving.
// Subsystem: flatten net derivation.
// It consumes part/world geometry and feeds sim with conductive and state-edge relationships.
// Flow position: preprocessing step between edit-time layout and run-time simulation.

import (
	"coilforge/internal/core"
	"coilforge/internal/world"
	"fmt"
	"math"
)

// BuildNets builds nets.
func BuildNets() {
	var anchors []core.PinAnchor
	var segs []core.Seg

	for _, p := range world.Parts {
		anchors = append(anchors, p.Anchors()...)
		segs = append(segs, p.Segments()...)
	}

	world.Nets = deriveNets(anchors, segs)
	world.PinNet = BuildPinNetMap(world.Nets)
}

// BuildPinNetMap builds pin net map.
func BuildPinNetMap(nets []core.Net) map[core.PinID]int {
	out := make(map[core.PinID]int, len(nets))
	for _, net := range nets {
		for _, pin := range net.Pins {
			out[pin] = net.ID
		}
	}
	return out
}

// deriveNets handles derive nets.
func deriveNets(anchors []core.PinAnchor, segs []core.Seg) []core.Net {
	type bucket struct {
		pins []core.PinID // pins value.
		segs []core.Seg   // segs value.
	}

	grouped := map[string]*bucket{}
	for _, anchor := range anchors {
		key := pointKey(anchor.Pt)
		b, ok := grouped[key]
		if !ok {
			b = &bucket{}
			grouped[key] = b
		}
		b.pins = append(b.pins, anchor.PinID)
	}

	for _, seg := range segs {
		key := pointKey(seg.A) + "->" + pointKey(seg.B)
		b, ok := grouped[key]
		if !ok {
			b = &bucket{}
			grouped[key] = b
		}
		b.segs = append(b.segs, seg)
	}

	nets := make([]core.Net, 0, len(grouped))
	id := 0
	for _, bucket := range grouped {
		nets = append(nets, core.Net{
			ID:   id,
			Pins: append([]core.PinID(nil), bucket.pins...),
			Segs: append([]core.Seg(nil), bucket.segs...),
		})
		id++
	}

	return nets
}

// pointKey handles point key.
func pointKey(pt core.Pt) string {
	return fmt.Sprintf("%d:%d", quantize(pt.X), quantize(pt.Y))
}

// quantize handles quantize.
func quantize(v float64) int64 {
	return int64(math.Round(v * 1000))
}
