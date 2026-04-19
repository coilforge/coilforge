package editor

// Post-pass wire normalization: split axis-aligned segments where another endpoint or a component pin
// lies strictly in the interior. One pass replaces all wire parts with 2-point straights at junctions.
//
// Future (same pass file): merge colinear overlapping segments; merge duplicate coincident vertices.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/part/catalog/wire"
	"coilforge/internal/world"
	"fmt"
	"math"
	"sort"
)

func junctionKey(p core.Pt) string {
	g := world.MajorGridWorld
	x := math.Round(p.X/g) * g
	y := math.Round(p.Y/g) * g
	return fmt.Sprintf("%.12g,%.12g", x, y)
}

// collectJunctionPoints returns unique snapped locations: every wire vertex and every non-wire pin anchor.
func collectJunctionPoints(parts []part.Part) []core.Pt {
	seen := make(map[string]bool)
	var out []core.Pt
	add := func(p core.Pt) {
		k := junctionKey(p)
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, p)
	}
	for _, p := range parts {
		if p.Base().TypeID == wireToolID {
			ww := p.(*wire.Wire)
			for _, q := range ww.Points {
				add(q)
			}
			continue
		}
		for _, a := range p.Anchors() {
			add(a.Pt)
		}
	}
	return out
}

func anyInteriorJunction(parts []part.Part, juncs []core.Pt) bool {
	for _, p := range parts {
		if p.Base().TypeID != wireToolID {
			continue
		}
		ww := p.(*wire.Wire)
		pts := ww.Points
		for i := 0; i < len(pts)-1; i++ {
			a, b := pts[i], pts[i+1]
			for _, q := range juncs {
				if strictInteriorOnAxisSeg(a, b, q) {
					return true
				}
			}
		}
	}
	return false
}

func strictInteriorOnAxisSeg(a, b, q core.Pt) bool {
	if approxEqPtNorm(q, a) || approxEqPtNorm(q, b) {
		return false
	}
	if !axisAlignedNorm(a, b) {
		return false
	}
	return pointInsideClosedAxisSeg(a, b, q)
}

func approxEqPtNorm(a, b core.Pt) bool {
	return math.Abs(a.X-b.X) < wireJoinCoordEps && math.Abs(a.Y-b.Y) < wireJoinCoordEps
}

func axisAlignedNorm(a, b core.Pt) bool {
	return math.Abs(a.X-b.X) < wireJoinCoordEps || math.Abs(a.Y-b.Y) < wireJoinCoordEps
}

func sortInteriorsAlongSeg(a, b core.Pt, interiors []core.Pt) []core.Pt {
	if len(interiors) == 0 {
		return interiors
	}
	cp := append([]core.Pt(nil), interiors...)
	if math.Abs(a.Y-b.Y) < wireJoinCoordEps {
		sort.Slice(cp, func(i, j int) bool { return cp[i].X < cp[j].X })
		return cp
	}
	sort.Slice(cp, func(i, j int) bool { return cp[i].Y < cp[j].Y })
	return cp
}

func dedupeConsecutivePts(pts []core.Pt) []core.Pt {
	if len(pts) == 0 {
		return pts
	}
	out := []core.Pt{pts[0]}
	for i := 1; i < len(pts); i++ {
		if approxEqPtNorm(pts[i], out[len(out)-1]) {
			continue
		}
		out = append(out, pts[i])
	}
	return out
}

// explodeWireToStraights returns 2-point edges covering the same geometry after inserting interior junctions.
func explodeWireToStraights(pts []core.Pt, juncs []core.Pt) [][]core.Pt {
	var edges [][]core.Pt
	for i := 0; i < len(pts)-1; i++ {
		a, b := pts[i], pts[i+1]
		var inter []core.Pt
		for _, q := range juncs {
			if strictInteriorOnAxisSeg(a, b, q) {
				inter = append(inter, q)
			}
		}
		inter = sortInteriorsAlongSeg(a, b, inter)
		chain := []core.Pt{a}
		chain = append(chain, inter...)
		chain = append(chain, b)
		chain = dedupeConsecutivePts(chain)
		for j := 0; j < len(chain)-1; j++ {
			if approxEqPtNorm(chain[j], chain[j+1]) {
				continue
			}
			edges = append(edges, []core.Pt{chain[j], chain[j+1]})
		}
	}
	return edges
}

// applyWireSegmentNormalization splits axis-aligned wire segments at interior junctions (pins + wire vertices).
//
// When recordUndo is true, pushes one undo snapshot before mutating (after wire strokes, where normalization is its own batch).
// When false, the caller must already own undo for the surrounding edit (paste, placement, drag) so revert restores geometry + splits together.
func applyWireSegmentNormalization(recordUndo bool) {
	juncs := collectJunctionPoints(world.Parts)
	if !anyInteriorJunction(world.Parts, juncs) {
		return
	}
	if recordUndo {
		pushUndo()
	}
	var out []part.Part
	for _, p := range world.Parts {
		if p.Base().TypeID != wireToolID {
			out = append(out, p)
			continue
		}
		ww := p.(*wire.Wire)
		edges := explodeWireToStraights(ww.Points, juncs)
		for _, e := range edges {
			if len(e) != 2 {
				continue
			}
			nw := wire.NewStraightWire(world.AllocPartID(), e[0], e[1], world.AllocPinID)
			if nw != nil {
				out = append(out, nw)
			}
		}
	}
	world.Parts = out
}
