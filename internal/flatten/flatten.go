package flatten

// File overview:
// flatten converts placed parts and wires into connectivity data for net solving.
// Subsystem: flatten net derivation.
// It consumes part/world geometry and feeds sim with conductive and state-edge relationships.
// Flow position: preprocessing step between edit-time layout and run-mode simulation.

import (
	"coilforge/internal/core"
	"coilforge/internal/world"
	"fmt"
	"math"
	"sort"
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
	keys := make([]string, 0)
	index := map[string]int{}

	addPt := func(p core.Pt) int {
		k := pointKey(p)
		if i, ok := index[k]; ok {
			return i
		}
		i := len(keys)
		keys = append(keys, k)
		index[k] = i
		return i
	}

	for _, a := range anchors {
		addPt(a.Pt)
	}
	for _, s := range segs {
		addPt(s.A)
		addPt(s.B)
	}

	if len(keys) == 0 {
		return nil
	}

	uf := newUnionFind(len(keys))
	for _, s := range segs {
		uf.union(index[pointKey(s.A)], index[pointKey(s.B)])
	}

	type bucket struct {
		pins []core.PinID
		segs []core.Seg
	}
	byRoot := map[int]*bucket{}

	bucketFor := func(p core.Pt) *bucket {
		i := index[pointKey(p)]
		r := uf.find(i)
		b := byRoot[r]
		if b == nil {
			b = &bucket{}
			byRoot[r] = b
		}
		return b
	}

	for _, a := range anchors {
		b := bucketFor(a.Pt)
		b.pins = append(b.pins, a.PinID)
	}
	for _, s := range segs {
		b := bucketFor(s.A)
		b.segs = append(b.segs, s)
	}

	roots := make([]int, 0, len(byRoot))
	for r := range byRoot {
		roots = append(roots, r)
	}
	sort.Ints(roots)

	nets := make([]core.Net, 0, len(roots))
	for nid, r := range roots {
		b := byRoot[r]
		nets = append(nets, core.Net{
			ID:   nid,
			Pins: append([]core.PinID(nil), b.pins...),
			Segs: append([]core.Seg(nil), b.segs...),
		})
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
