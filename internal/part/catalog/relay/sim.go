package relay

// File overview:
// sim implements simulation-facing behavior for relay using part sim interfaces.
// Subsystem: part catalog (relay) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"math/rand"
)

// AddConductive adds conductive.
func (self *Relay) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	self.ensureContactSlices()
	for idx, pole := range self.Poles {
		commonNet := netByPin(pole.PinCommon)
		switch self.Contacts[idx] {
		case ContactNO:
			union.Union(commonNet, netByPin(pole.PinNO))
		default:
			union.Union(commonNet, netByPin(pole.PinNC))
		}
	}
}

// Tick handles tick.
func (self *Relay) Tick(ctx part.SimContext) bool {
	self.ensureContactSlices()

	changed := false
	if self.TransitionScheduled && ctx.Tick >= self.TransitionDueTick {
		copy(self.Contacts, self.PendingContacts)
		self.TransitionScheduled = false
		changed = true
	}

	coilNet := ctx.NetByPin(self.PinCoilA)
	coilHigh := ctx.NetState(coilNet) == core.NetHigh
	if coilHigh != self.CoilActive {
		self.CoilActive = coilHigh
		delay := self.pickupTicks(ctx.TickMicros)
		if !coilHigh {
			delay = self.releaseTicks(ctx.TickMicros)
		}
		if ctx.EnableJitter {
			delay += self.jitterTicks(ctx.TickMicros, ctx.Rand)
		}
		self.computePendingContacts()
		self.TransitionDueTick = ctx.Tick + uint64(delay)
		self.TransitionScheduled = true
	}

	return changed
}

// computePendingContacts handles compute pending contacts.
func (self *Relay) computePendingContacts() {
	target := ContactNC
	if self.CoilActive {
		target = ContactNO
	}
	for i := range self.PendingContacts {
		self.PendingContacts[i] = target
	}
}

// pickupTicks handles pickup ticks.
func (self *Relay) pickupTicks(tickMicros int) int {
	return millisToTicks(self.PickupMs+self.FlightMs, tickMicros)
}

// releaseTicks handles release ticks.
func (self *Relay) releaseTicks(tickMicros int) int {
	return millisToTicks(self.ReleaseMs+self.FlightMs, tickMicros)
}

// jitterTicks handles jitter ticks.
func (self *Relay) jitterTicks(tickMicros int, randSource *rand.Rand) int {
	if self.JitterMs <= 0 || randSource == nil {
		return 0
	}
	return millisToTicks(randSource.Intn(self.JitterMs+1), tickMicros)
}

// millisToTicks handles millis to ticks.
func millisToTicks(ms, tickMicros int) int {
	if ms <= 0 || tickMicros <= 0 {
		return 0
	}
	return (ms * 1000) / tickMicros
}
