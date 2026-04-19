package clock

// File overview:
// sim implements NetSeeder for clock (square wave low/high on OUT vs simulated time).
// Subsystem: part catalog (clock) simulation.
// Flow position: net seeding during resolveNets.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// HalfPeriodMicros is one on or off phase of the output (full cycle = 2×).
// Fixed at 250ms per phase until a properties UI exposes timing.
const HalfPeriodMicros = 250_000

// SeedNets drives the attached net high or low from simulated time (square wave).
func (self *Clock) SeedNets(union part.NetUnion, netByPin func(core.PinID) int, high, low map[int]bool, nowMicros uint64) {
	netID := netByPin(self.OUT)
	if netID < 0 {
		return
	}
	root := union.Find(netID)
	if root < 0 {
		return
	}
	half := uint64(HalfPeriodMicros)
	if half == 0 {
		return
	}
	phase := (nowMicros / half) % 2
	if phase == 0 {
		high[root] = true
	} else {
		low[root] = true
	}
}
