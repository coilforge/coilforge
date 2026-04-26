package clock

// File overview:
// sim implements NetSeeder for clock (square wave low/high on OUT vs simulated time).
// Subsystem: part catalog (clock) simulation.
// Flow position: net seeding during resolveNets.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// clockHalfPhaseIters is sim-loop iterations per on/off phase. Chosen so, at ~1:1 sim-µs vs wall-µs (e.g. 10k ticks/s with
// 100 µs [core.SimStepMicros]), one half-phase is ~0.5 s wall → ~1 Hz full square wave — visible without full-speed fast-forward.
const clockHalfPhaseIters = 5_000

// SeedNets drives the attached net high or low from simulated time (square wave).
func (self *Clock) SeedNets(union part.NetUnion, netByPin func(core.PinID) int, high, low map[int]bool, nowMicros uint64) {
	netID := netByPin(self.CLK)
	if netID < 0 {
		return
	}
	root := union.Find(netID)
	if root < 0 {
		return
	}
	half := uint64(clockHalfPhaseIters) * uint64(core.SimStepMicros)
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

// HalfPeriodMicros returns one clock half-cycle duration in simulated microseconds (matches [SeedNets] timing).
func HalfPeriodMicros() uint64 {
	return uint64(clockHalfPhaseIters) * uint64(core.SimStepMicros)
}

// PhaseAt returns 0 while OUT is seeded high and 1 while seeded low ([SeedNets] semantics).
func PhaseAt(nowMicros uint64) int {
	half := HalfPeriodMicros()
	if half == 0 {
		return 0
	}
	return int((nowMicros / half) % 2)
}

// PhaseLabel returns "high" or "low" for logging and diagnostics ([PhaseAt] mapping).
func PhaseLabel(nowMicros uint64) string {
	if PhaseAt(nowMicros) == 0 {
		return "high"
	}
	return "low"
}
