package relay

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// AddConductive closes the contact pair represented by the mid section.
// De-energized: COM <-> NC, energized: COM <-> NO.
func (self *Relay) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	for pole := 1; pole <= self.poleCountClamped(); pole++ {
		comPin, ncPin, noPin := self.contactPinsForPole(pole)
		com := netByPin(comPin)
		if com < 0 {
			continue
		}
		var other int
		if self.Energized {
			other = netByPin(noPin)
		} else {
			other = netByPin(ncPin)
		}
		if other < 0 {
			continue
		}
		union.Union(com, other)
	}
}

// Tick updates relay energized state from coil pin net states.
// Same logic as indicator: energized only when both coil terminals are driven and opposite.
func (self *Relay) Tick(ctx part.SimContext) bool {
	was := self.Energized
	a := ctx.PinNetState(self.CoilA)
	b := ctx.PinNetState(self.CoilB)
	self.Energized = a != core.NetFloat && b != core.NetFloat && a != b
	return self.Energized != was
}
