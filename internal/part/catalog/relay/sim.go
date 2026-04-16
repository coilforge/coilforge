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
func (r *Relay) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	r.ensureContactSlices()
	for idx, pole := range r.Poles {
		commonNet := netByPin(pole.PinCommon)
		switch r.Contacts[idx] {
		case ContactNO:
			union.Union(commonNet, netByPin(pole.PinNO))
		default:
			union.Union(commonNet, netByPin(pole.PinNC))
		}
	}
}

// Tick handles tick.
func (r *Relay) Tick(ctx part.SimContext) bool {
	r.ensureContactSlices()

	changed := false
	if r.TransitionScheduled && ctx.Tick >= r.TransitionDueTick {
		copy(r.Contacts, r.PendingContacts)
		r.TransitionScheduled = false
		changed = true
	}

	coilNet := ctx.NetByPin(r.PinCoilA)
	coilHigh := ctx.NetState(coilNet) == core.NetHigh
	if coilHigh != r.CoilActive {
		r.CoilActive = coilHigh
		delay := r.pickupTicks(ctx.TickMicros)
		if !coilHigh {
			delay = r.releaseTicks(ctx.TickMicros)
		}
		if ctx.EnableJitter {
			delay += r.jitterTicks(ctx.TickMicros, ctx.Rand)
		}
		r.computePendingContacts()
		r.TransitionDueTick = ctx.Tick + uint64(delay)
		r.TransitionScheduled = true
	}

	return changed
}

// computePendingContacts handles compute pending contacts.
func (r *Relay) computePendingContacts() {
	target := ContactNC
	if r.CoilActive {
		target = ContactNO
	}
	for i := range r.PendingContacts {
		r.PendingContacts[i] = target
	}
}

// pickupTicks handles pickup ticks.
func (r *Relay) pickupTicks(tickMicros int) int {
	return millisToTicks(r.PickupMs+r.FlightMs, tickMicros)
}

// releaseTicks handles release ticks.
func (r *Relay) releaseTicks(tickMicros int) int {
	return millisToTicks(r.ReleaseMs+r.FlightMs, tickMicros)
}

// jitterTicks handles jitter ticks.
func (r *Relay) jitterTicks(tickMicros int, randSource *rand.Rand) int {
	if r.JitterMs <= 0 || randSource == nil {
		return 0
	}
	return millisToTicks(randSource.Intn(r.JitterMs+1), tickMicros)
}

// millisToTicks handles millis to ticks.
func millisToTicks(ms, tickMicros int) int {
	if ms <= 0 || tickMicros <= 0 {
		return 0
	}
	return (ms * 1000) / tickMicros
}
