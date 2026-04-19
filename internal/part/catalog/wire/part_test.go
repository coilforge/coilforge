package wire

import (
	"coilforge/internal/core"
	"testing"
)

func TestNewStraightWireKeepsTwoPointsWhenOrthogonalRouteAddsCorner(t *testing.T) {
	allocPin := func() func() core.PinID {
		n := core.PinID(9000)
		return func() core.PinID {
			n++
			return n
		}
	}()
	a := core.Pt{X: 0, Y: 0}
	// Tiny Y drift across a long span: strict axis checks fail; OrthogonalRoute inserts an elbow.
	b := core.Pt{X: 1_000_000, Y: 1e-12}
	if len(OrthogonalRoute(a, b)) != 3 {
		t.Fatalf("OrthogonalRoute got %d points, want 3 (diagonal-ish pair)", len(OrthogonalRoute(a, b)))
	}
	w := NewStraightWire(1, a, b, allocPin).(*Wire)
	if len(w.Points) != 2 {
		t.Fatalf("NewStraightWire Points=%d want 2", len(w.Points))
	}
}
