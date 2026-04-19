package flatten

import (
	"coilforge/internal/core"
	"testing"
)

func TestDeriveNetsUnionsPinsThroughWireSegment(t *testing.T) {
	a := core.PinAnchor{Pt: core.Pt{X: 0, Y: 0}, PinID: 1}
	b := core.PinAnchor{Pt: core.Pt{X: 4, Y: 0}, PinID: 2}
	segs := []core.Seg{
		{A: core.Pt{X: 0, Y: 0}, B: core.Pt{X: 4, Y: 0}},
	}
	nets := deriveNets([]core.PinAnchor{a, b}, segs)
	if len(nets) != 1 {
		t.Fatalf("expected 1 net, got %d", len(nets))
	}
	found1, found2 := false, false
	for _, id := range nets[0].Pins {
		if id == 1 {
			found1 = true
		}
		if id == 2 {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Fatalf("expected both pins in same net, pins=%v", nets[0].Pins)
	}
}

func TestDeriveNetsIsolatedPinsStaySeparate(t *testing.T) {
	a := core.PinAnchor{Pt: core.Pt{X: 0, Y: 0}, PinID: 1}
	b := core.PinAnchor{Pt: core.Pt{X: 100, Y: 100}, PinID: 2}
	nets := deriveNets([]core.PinAnchor{a, b}, nil)
	if len(nets) != 2 {
		t.Fatalf("expected 2 nets, got %d", len(nets))
	}
}
