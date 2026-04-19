package editor

// Wire branching: start routing from a click on an existing wire body (junction snapped to major grid).

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/part/catalog/wire"
	"coilforge/internal/world"
	"math"
)

const wireJoinCoordEps = 1e-6

func wireBranchHit(pt core.Pt) (int, core.Pt, bool) {
	for i := len(world.Parts) - 1; i >= 0; i-- {
		p := world.Parts[i]
		if p.Base().TypeID != wireToolID {
			continue
		}
		hr := p.HitTest(pt)
		if !hr.Hit || hr.Kind != part.HitBody {
			continue
		}
		ww := p.(*wire.Wire)
		junction, ok := junctionOnWire(ww, pt)
		if !ok {
			continue
		}
		return i, junction, true
	}
	return 0, core.Pt{}, false
}

func junctionOnWire(ww *wire.Wire, click core.Pt) (core.Pt, bool) {
	pts := ww.Points
	if len(pts) < 2 {
		return core.Pt{}, false
	}
	var best core.Pt
	bestD := math.MaxFloat64
	for i := 0; i < len(pts)-1; i++ {
		a, b := pts[i], pts[i+1]
		cl := closestPointOnSegment(a, b, click)
		snapped := snapOntoMajorGridAlongSegment(a, b, cl)
		if !pointInsideClosedAxisSeg(a, b, snapped) {
			continue
		}
		d := distSq(snapped, click)
		if d < bestD {
			bestD = d
			best = snapped
		}
	}
	if bestD >= math.MaxFloat64/2 {
		return core.Pt{}, false
	}
	return best, true
}

func closestPointOnSegment(a, b, p core.Pt) core.Pt {
	dx := b.X - a.X
	dy := b.Y - a.Y
	lenSq := dx*dx + dy*dy
	if lenSq == 0 {
		return a
	}
	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / lenSq
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return core.Pt{X: a.X + t*dx, Y: a.Y + t*dy}
}

func snapOntoMajorGridAlongSegment(a, b, close core.Pt) core.Pt {
	g := world.MajorGridWorld
	if math.Abs(a.Y-b.Y) < wireJoinCoordEps {
		x := math.Round(close.X/g) * g
		minX, maxX := math.Min(a.X, b.X), math.Max(a.X, b.X)
		if x < minX {
			x = minX
		}
		if x > maxX {
			x = maxX
		}
		return core.Pt{X: x, Y: a.Y}
	}
	if math.Abs(a.X-b.X) < wireJoinCoordEps {
		y := math.Round(close.Y/g) * g
		minY, maxY := math.Min(a.Y, b.Y), math.Max(a.Y, b.Y)
		if y < minY {
			y = minY
		}
		if y > maxY {
			y = maxY
		}
		return core.Pt{X: a.X, Y: y}
	}
	return close
}

func pointInsideClosedAxisSeg(a, b, p core.Pt) bool {
	minX, maxX := math.Min(a.X, b.X), math.Max(a.X, b.X)
	minY, maxY := math.Min(a.Y, b.Y), math.Max(a.Y, b.Y)
	if p.X < minX-wireJoinCoordEps || p.X > maxX+wireJoinCoordEps {
		return false
	}
	if p.Y < minY-wireJoinCoordEps || p.Y > maxY+wireJoinCoordEps {
		return false
	}
	if math.Abs(a.Y-b.Y) < wireJoinCoordEps {
		return math.Abs(p.Y-a.Y) < wireJoinCoordEps*10
	}
	if math.Abs(a.X-b.X) < wireJoinCoordEps {
		return math.Abs(p.X-a.X) < wireJoinCoordEps*10
	}
	return false
}

func distSq(a, b core.Pt) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}
