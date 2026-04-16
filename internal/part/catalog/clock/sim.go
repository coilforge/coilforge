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
func (c *Clock) Tick(ctx part.SimContext) bool {
	if c.PeriodTick <= 0 {
		return false
	}
	prev := c.OutputHigh
	phase := int(ctx.Tick % uint64(c.PeriodTick))
	c.OutputHigh = phase < c.HighTick
	return c.OutputHigh != prev
}

// SeedNets seeds nets.
func (c *Clock) SeedNets(netByPin func(core.PinID) int, high, low map[int]bool) {
	net := netByPin(c.PinOut)
	if net < 0 {
		return
	}
	if c.OutputHigh {
		high[net] = true
		return
	}
	low[net] = true
}
