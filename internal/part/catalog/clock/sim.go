package clock

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (c *Clock) Tick(ctx part.SimContext) bool {
	if c.PeriodTick <= 0 {
		return false
	}
	prev := c.OutputHigh
	phase := int(ctx.Tick % uint64(c.PeriodTick))
	c.OutputHigh = phase < c.HighTick
	return c.OutputHigh != prev
}

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
