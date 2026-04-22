package relay

import (
	"coilforge/internal/core"
	"testing"
)

type unionSpy struct {
	pairs [][2]int
}

func (u *unionSpy) Union(a, b int) { u.pairs = append(u.pairs, [2]int{a, b}) }
func (u *unionSpy) Find(a int) int  { return a }

func TestRelayAddConductiveDeenergizedUsesNC(t *testing.T) {
	r := &Relay{RelayPinIDs: RelayPinIDs{
		COM: 1,
		NC:  2,
		NO:  3,
	}}
	u := &unionSpy{}
	netByPin := func(pin core.PinID) int {
		switch pin {
		case r.COM:
			return 10
		case r.NC:
			return 20
		case r.NO:
			return 30
		default:
			return -1
		}
	}

	r.Energized = false
	r.AddConductive(u, netByPin)
	if len(u.pairs) != 1 || u.pairs[0] != [2]int{10, 20} {
		t.Fatalf("expected COM-NC union, got %#v", u.pairs)
	}
}

func TestRelayAddConductiveEnergizedUsesNO(t *testing.T) {
	r := &Relay{RelayPinIDs: RelayPinIDs{
		COM: 1,
		NC:  2,
		NO:  3,
	}}
	u := &unionSpy{}
	netByPin := func(pin core.PinID) int {
		switch pin {
		case r.COM:
			return 10
		case r.NC:
			return 20
		case r.NO:
			return 30
		default:
			return -1
		}
	}

	r.Energized = true
	r.AddConductive(u, netByPin)
	if len(u.pairs) != 1 || u.pairs[0] != [2]int{10, 30} {
		t.Fatalf("expected COM-NO union, got %#v", u.pairs)
	}
}

