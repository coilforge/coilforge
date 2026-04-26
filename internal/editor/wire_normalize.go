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

const normEps = 1e-9

func junctionKey(p core.Pt) string {
	s := snapToMajorGrid(p)
	x := s.X
	y := s.Y
	return fmt.Sprintf("%.12g,%.12g", x, y)
}

func snapToMajorGrid(p core.Pt) core.Pt {
	g := world.MajorGridWorld
	return core.Pt{
		X: math.Round(p.X/g) * g,
		Y: math.Round(p.Y/g) * g,
	}
}

// collectJunctionPoints returns unique snapped locations: every wire vertex and every non-wire pin anchor.
func collectJunctionPoints(parts []part.Part) []core.Pt {
	seen := make(map[string]bool)
	var out []core.Pt
	add := func(p core.Pt) {
		s := snapToMajorGrid(p)
		k := junctionKey(s)
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, s)
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
			a, b := snapToMajorGrid(pts[i]), snapToMajorGrid(pts[i+1])
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

type axisEdge struct {
	A core.Pt
	B core.Pt
}

type axisInterval struct {
	lo float64
	hi float64
}

func sortedInterval(a, b float64) (float64, float64) {
	if a <= b {
		return a, b
	}
	return b, a
}

func collectCurrentStraightEdges(parts []part.Part) []axisEdge {
	var edges []axisEdge
	for _, p := range parts {
		if p.Base().TypeID != wireToolID {
			continue
		}
		ww := p.(*wire.Wire)
		for i := 0; i < len(ww.Points)-1; i++ {
			a := snapToMajorGrid(ww.Points[i])
			b := snapToMajorGrid(ww.Points[i+1])
			if !axisAlignedNorm(a, b) || approxEqPtNorm(a, b) {
				continue
			}
			edges = append(edges, axisEdge{A: a, B: b})
		}
	}
	return edges
}

func collectExplodedEdges(parts []part.Part, juncs []core.Pt) []axisEdge {
	var out []axisEdge
	for _, p := range parts {
		if p.Base().TypeID != wireToolID {
			continue
		}
		ww := p.(*wire.Wire)
		edges := explodeWireToStraights(ww.Points, juncs)
		for _, e := range edges {
			if len(e) != 2 {
				continue
			}
			a := snapToMajorGrid(e[0])
			b := snapToMajorGrid(e[1])
			if !axisAlignedNorm(a, b) || approxEqPtNorm(a, b) {
				continue
			}
			out = append(out, axisEdge{A: a, B: b})
		}
	}
	return out
}

func mergeIntervals(intervals []axisInterval) []axisInterval {
	if len(intervals) == 0 {
		return nil
	}
	cp := append([]axisInterval(nil), intervals...)
	sort.Slice(cp, func(i, j int) bool { return cp[i].lo < cp[j].lo })
	merged := []axisInterval{cp[0]}
	for i := 1; i < len(cp); i++ {
		cur := cp[i]
		last := &merged[len(merged)-1]
		if cur.lo <= last.hi+normEps {
			if cur.hi > last.hi {
				last.hi = cur.hi
			}
			continue
		}
		merged = append(merged, cur)
	}
	return merged
}

func sortUniqueFloats(vals []float64) []float64 {
	if len(vals) == 0 {
		return nil
	}
	sort.Float64s(vals)
	out := vals[:1]
	for i := 1; i < len(vals); i++ {
		if math.Abs(vals[i]-out[len(out)-1]) < normEps {
			continue
		}
		out = append(out, vals[i])
	}
	return out
}

func collectNonWirePinPoints(parts []part.Part) []core.Pt {
	var pins []core.Pt
	for _, p := range parts {
		if p.Base().TypeID == wireToolID {
			continue
		}
		for _, a := range p.Anchors() {
			pins = append(pins, snapToMajorGrid(a.Pt))
		}
	}
	return pins
}

func collectEndpointPoints(edges []axisEdge) (hEnds []core.Pt, vEnds []core.Pt) {
	for _, e := range edges {
		a := snapToMajorGrid(e.A)
		b := snapToMajorGrid(e.B)
		if math.Abs(a.X-b.X) < wireJoinCoordEps {
			vEnds = append(vEnds, a, b)
		} else if math.Abs(a.Y-b.Y) < wireJoinCoordEps {
			hEnds = append(hEnds, a, b)
		}
	}
	return hEnds, vEnds
}

func splitMergedIntervals(merged []axisInterval, cuts []float64) []axisInterval {
	if len(merged) == 0 {
		return nil
	}
	cuts = sortUniqueFloats(cuts)
	var out []axisInterval
	for _, iv := range merged {
		lineCuts := []float64{iv.lo, iv.hi}
		for _, c := range cuts {
			if c > iv.lo+normEps && c < iv.hi-normEps {
				lineCuts = append(lineCuts, c)
			}
		}
		lineCuts = sortUniqueFloats(lineCuts)
		for i := 0; i < len(lineCuts)-1; i++ {
			lo, hi := lineCuts[i], lineCuts[i+1]
			if hi-lo < normEps {
				continue
			}
			out = append(out, axisInterval{lo: lo, hi: hi})
		}
	}
	return out
}

func normalizeEdge(e axisEdge) axisEdge {
	if e.A.X < e.B.X || (math.Abs(e.A.X-e.B.X) < normEps && e.A.Y <= e.B.Y) {
		return e
	}
	return axisEdge{A: e.B, B: e.A}
}

func edgeSig(e axisEdge) string {
	n := normalizeEdge(e)
	return fmt.Sprintf("%.12g,%.12g->%.12g,%.12g", n.A.X, n.A.Y, n.B.X, n.B.Y)
}

func edgeSetsEqual(a, b []axisEdge) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[string]int, len(a))
	for _, e := range a {
		ma[edgeSig(e)]++
	}
	for _, e := range b {
		k := edgeSig(e)
		if ma[k] == 0 {
			return false
		}
		ma[k]--
	}
	for _, n := range ma {
		if n != 0 {
			return false
		}
	}
	return true
}

func canonicalizeColinearEdges(edges []axisEdge, pinPts []core.Pt) []axisEdge {
	vertical, horizontal := groupAxisEdges(edges)
	hEnds, vEnds := collectEndpointPoints(edges)
	out := make([]axisEdge, 0, len(edges))
	out = append(out, expandVerticalCanonicalEdges(vertical, pinPts, hEnds)...)
	out = append(out, expandHorizontalCanonicalEdges(horizontal, pinPts, vEnds)...)
	sortCanonicalEdges(out)
	return out
}

func groupAxisEdges(edges []axisEdge) (
	map[string]struct {
		x         float64
		intervals []axisInterval
	},
	map[string]struct {
		y         float64
		intervals []axisInterval
	},
) {
	vertical := make(map[string]struct {
		x         float64
		intervals []axisInterval
	})
	horizontal := make(map[string]struct {
		y         float64
		intervals []axisInterval
	})
	for _, e := range edges {
		if math.Abs(e.A.X-e.B.X) < wireJoinCoordEps {
			lo, hi := sortedInterval(e.A.Y, e.B.Y)
			if hi-lo < normEps {
				continue
			}
			k := fmt.Sprintf("v:%.12g", e.A.X)
			entry := vertical[k]
			entry.x = e.A.X
			entry.intervals = append(entry.intervals, axisInterval{lo: lo, hi: hi})
			vertical[k] = entry
		} else if math.Abs(e.A.Y-e.B.Y) < wireJoinCoordEps {
			lo, hi := sortedInterval(e.A.X, e.B.X)
			if hi-lo < normEps {
				continue
			}
			k := fmt.Sprintf("h:%.12g", e.A.Y)
			entry := horizontal[k]
			entry.y = e.A.Y
			entry.intervals = append(entry.intervals, axisInterval{lo: lo, hi: hi})
			horizontal[k] = entry
		}
	}
	return vertical, horizontal
}

func expandVerticalCanonicalEdges(
	vertical map[string]struct {
		x         float64
		intervals []axisInterval
	},
	pinPts []core.Pt,
	hEnds []core.Pt,
) []axisEdge {
	var out []axisEdge
	for _, entry := range vertical {
		merged := mergeIntervals(entry.intervals)
		var cuts []float64
		for _, p := range pinPts {
			if math.Abs(p.X-entry.x) < wireJoinCoordEps {
				cuts = append(cuts, p.Y)
			}
		}
		for _, p := range hEnds {
			if math.Abs(p.X-entry.x) < wireJoinCoordEps {
				cuts = append(cuts, p.Y)
			}
		}
		for _, iv := range splitMergedIntervals(merged, cuts) {
			out = append(out, axisEdge{
				A: core.Pt{X: entry.x, Y: iv.lo},
				B: core.Pt{X: entry.x, Y: iv.hi},
			})
		}
	}
	return out
}

func expandHorizontalCanonicalEdges(
	horizontal map[string]struct {
		y         float64
		intervals []axisInterval
	},
	pinPts []core.Pt,
	vEnds []core.Pt,
) []axisEdge {
	var out []axisEdge
	for _, entry := range horizontal {
		merged := mergeIntervals(entry.intervals)
		var cuts []float64
		for _, p := range pinPts {
			if math.Abs(p.Y-entry.y) < wireJoinCoordEps {
				cuts = append(cuts, p.X)
			}
		}
		for _, p := range vEnds {
			if math.Abs(p.Y-entry.y) < wireJoinCoordEps {
				cuts = append(cuts, p.X)
			}
		}
		for _, iv := range splitMergedIntervals(merged, cuts) {
			out = append(out, axisEdge{
				A: core.Pt{X: iv.lo, Y: entry.y},
				B: core.Pt{X: iv.hi, Y: entry.y},
			})
		}
	}
	return out
}

func sortCanonicalEdges(out []axisEdge) {
	sort.Slice(out, func(i, j int) bool {
		if out[i].A.X != out[j].A.X {
			return out[i].A.X < out[j].A.X
		}
		if out[i].A.Y != out[j].A.Y {
			return out[i].A.Y < out[j].A.Y
		}
		if out[i].B.X != out[j].B.X {
			return out[i].B.X < out[j].B.X
		}
		return out[i].B.Y < out[j].B.Y
	})
}

// explodeWireToStraights returns 2-point edges covering the same geometry after inserting interior junctions.
func explodeWireToStraights(pts []core.Pt, juncs []core.Pt) [][]core.Pt {
	var edges [][]core.Pt
	for i := 0; i < len(pts)-1; i++ {
		a, b := snapToMajorGrid(pts[i]), snapToMajorGrid(pts[i+1])
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
	currentEdges := collectCurrentStraightEdges(world.Parts)
	needsSplit := anyInteriorJunction(world.Parts, juncs)
	needsPolylineCleanup := false
	for _, p := range world.Parts {
		if p.Base().TypeID == wireToolID && len(p.(*wire.Wire).Points) > 2 {
			needsPolylineCleanup = true
			break
		}
	}
	canonicalEdges := canonicalizeColinearEdges(
		collectExplodedEdges(world.Parts, juncs),
		collectNonWirePinPoints(world.Parts),
	)
	needsGeometryRewrite := !edgeSetsEqual(currentEdges, canonicalEdges)
	if !needsSplit && !needsPolylineCleanup && !needsGeometryRewrite {
		return
	}
	if recordUndo {
		pushUndo()
	}
	var out []part.Part
	for _, p := range world.Parts {
		if p.Base().TypeID != wireToolID {
			out = append(out, p)
		}
	}
	for _, e := range canonicalEdges {
		nw := wire.NewStraightWire(world.AllocPartID(), e.A, e.B, world.AllocPinID)
		if nw != nil {
			out = append(out, nw)
		}
	}
	world.Parts = out
}

// NormalizeWiresInPlace canonicalizes wire geometry in [world.Parts] without recording undo.
// Used by project load to clean overlapping legacy segments before interaction starts.
func NormalizeWiresInPlace() {
	applyWireSegmentNormalization(false)
}
