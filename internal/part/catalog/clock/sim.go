package clock

// File overview:
// sim implements NetSeeder for clock (square wave low/high on OUT vs simulated time).
// Subsystem: part catalog (clock) simulation.
// Flow position: net seeding during resolveNets.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

const (
	defaultOnMs  = 500
	defaultOffMs = 500
	minPhaseMs   = 1
	maxPhaseMs   = 60_000
)

func clampPhaseMs(v int) int {
	if v == 0 {
		return defaultOnMs
	}
	if v < minPhaseMs {
		return minPhaseMs
	}
	if v > maxPhaseMs {
		return maxPhaseMs
	}
	return v
}

func msToMicros(v int) uint64 {
	return uint64(clampPhaseMs(v)) * 1_000
}

func (self *Clock) onMicros() uint64 {
	return msToMicros(self.OnMs)
}

func (self *Clock) offMicros() uint64 {
	return msToMicros(self.OffMs)
}

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
	on := self.onMicros()
	off := self.offMicros()
	cycle := on + off
	if cycle == 0 {
		return
	}
	phase := nowMicros % cycle
	if phase < on {
		high[root] = true
	} else {
		low[root] = true
	}
}

// HalfPeriodMicros returns the default half-cycle duration used by freshly created clocks.
func HalfPeriodMicros() uint64 {
	return uint64(defaultOnMs) * 1_000
}

// PhaseAt returns phase for the default 500ms/500ms clock.
func PhaseAt(nowMicros uint64) int {
	on := uint64(defaultOnMs) * 1_000
	off := uint64(defaultOffMs) * 1_000
	cycle := on + off
	if cycle == 0 {
		return 0
	}
	if nowMicros%cycle < on {
		return 0
	}
	return 1
}

// PhaseLabel returns "high" or "low" for logging and diagnostics ([PhaseAt] mapping).
func PhaseLabel(nowMicros uint64) string {
	if PhaseAt(nowMicros) == 0 {
		return "high"
	}
	return "low"
}
