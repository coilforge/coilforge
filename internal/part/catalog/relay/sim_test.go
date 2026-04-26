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
		COM1: 1,
		NC1:  2,
		NO1:  3,
	}, PoleCount: 1}
	u := &unionSpy{}
	netByPin := func(pin core.PinID) int {
		switch pin {
		case r.COM1:
			return 10
		case r.NC1:
			return 20
		case r.NO1:
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
		COM1: 1,
		NC1:  2,
		NO1:  3,
	}, PoleCount: 1}
	u := &unionSpy{}
	netByPin := func(pin core.PinID) int {
		switch pin {
		case r.COM1:
			return 10
		case r.NC1:
			return 20
		case r.NO1:
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

func TestRelayAddConductiveMultiPole(t *testing.T) {
	r := &Relay{RelayPinIDs: RelayPinIDs{
		COM1: 1, NC1: 2, NO1: 3,
		COM2: 4, NC2: 5, NO2: 6,
	}, PoleCount: 2}
	u := &unionSpy{}
	netByPin := func(pin core.PinID) int {
		switch pin {
		case r.COM1:
			return 10
		case r.NC1:
			return 20
		case r.NO1:
			return 30
		case r.COM2:
			return 40
		case r.NC2:
			return 50
		case r.NO2:
			return 60
		default:
			return -1
		}
	}
	r.Energized = false
	r.AddConductive(u, netByPin)
	if len(u.pairs) != 2 {
		t.Fatalf("expected two unions, got %#v", u.pairs)
	}
	if u.pairs[0] != [2]int{10, 20} || u.pairs[1] != [2]int{40, 50} {
		t.Fatalf("expected COMi-NCi unions, got %#v", u.pairs)
	}
}

