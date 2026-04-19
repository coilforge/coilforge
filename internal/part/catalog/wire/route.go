package wire

import "coilforge/internal/core"

// OrthogonalRoute returns rectilinear waypoints from a to b: straight segment if aligned,
// otherwise an L-shape with a corner at (b.X, a.Y).
func OrthogonalRoute(a, b core.Pt) []core.Pt {
	if a.X == b.X && a.Y == b.Y {
		return []core.Pt{a}
	}
	if a.X == b.X || a.Y == b.Y {
		return []core.Pt{a, b}
	}
	return []core.Pt{a, {X: b.X, Y: a.Y}, b}
}
