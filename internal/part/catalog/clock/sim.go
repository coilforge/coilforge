package clock

// File overview:
// sim implements simulation-facing behavior for clock using part sim interfaces.
// Subsystem: part catalog (clock) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Tick handles tick.
func (self *Clock) Tick(ctx part.SimContext) bool {
	if self.PeriodTick <= 0 {
		return false
	}
	prev := self.OutputHigh
	phase := int(ctx.Tick % uint64(self.PeriodTick))
	self.OutputHigh = phase < self.HighTick
	return self.OutputHigh != prev
}

// SeedNets seeds nets.
func (self *Clock) SeedNets(netByPin func(core.PinID) int, high, low map[int]bool) {
	net := netByPin(self.PinOut)
	if net < 0 {
		return
	}
	if self.OutputHigh {
		high[net] = true
		return
	}
	low[net] = true
}
