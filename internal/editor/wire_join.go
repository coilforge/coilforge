package editor

// Wire branching: start routing from a click on an existing wire body (junction snapped to major grid).

import (
	"coilforge/internal/core"
	"math"
)

const wireJoinCoordEps = 1e-6

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
