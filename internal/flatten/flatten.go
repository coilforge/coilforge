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
	"log"
	"math"
	"os"
	"sort"
)

// BuildNets builds nets.
func BuildNets() {
	var anchors []core.PinAnchor
	var segs []core.Seg

	for i, p := range world.Parts {
		a := p.Anchors()
		anchors = append(anchors, a...)
		segs = append(segs, p.Segments()...)
		if flattenTrace() {
			b := p.Base()
			log.Printf("[flatten] part[%d] type=%s id=%d pos=(%.4f,%.4f) rot=%d mirror=%v anchors=%d segs=%d",
				i, b.TypeID, b.ID, b.Pos.X, b.Pos.Y, b.Rotation, b.Mirror, len(a), len(p.Segments()))
		}
	}

	if flattenTrace() {
		log.Printf("[flatten] BuildNets: %d parts, %d anchors, %d seg groups (wires)", len(world.Parts), len(anchors), len(segs))
		for _, an := range anchors {
			log.Printf("[flatten]   anchor pin=%d at (%.4f,%.4f) key=%q", an.PinID, an.Pt.X, an.Pt.Y, pointKey(an.Pt))
		}
	}

	world.Nets = deriveNets(anchors, segs)
	world.PinNet = BuildPinNetMap(world.Nets)

	if flattenTrace() {
		log.Printf("[flatten] derived %d nets", len(world.Nets))
		for _, n := range world.Nets {
			log.Printf("[flatten]   net id=%d pins=%v segs=%d", n.ID, n.Pins, len(n.Segs))
		}
		pinIDs := make([]int, 0, len(world.PinNet))
		for pid := range world.PinNet {
			pinIDs = append(pinIDs, int(pid))
		}
		sort.Ints(pinIDs)
		for _, ip := range pinIDs {
			pid := core.PinID(ip)
			log.Printf("[flatten]   PinNet pin=%d -> net %d", pid, world.PinNet[pid])
		}
	}
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

// flattenTrace is true when COILFORGE_FLATTEN_TRACE is set (any non-empty value).
func flattenTrace() bool {
	return os.Getenv("COILFORGE_FLATTEN_TRACE") != ""
}
