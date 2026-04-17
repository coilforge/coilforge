package switches

// File overview:
// sim implements simulation-facing behavior for switches using part sim interfaces.
// Subsystem: part catalog (switches) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// AddConductive adds conductive.
func (self *Switch) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	if !self.effectiveClosed() {
		return
	}
	union.Union(netByPin(self.PinA), netByPin(self.PinB))
}

// HandleInput handles input.
func (self *Switch) HandleInput(active bool) (changed, momentary bool) {
	if self.Momentary {
		prev := self.Pressed
		self.Pressed = active
		return self.Pressed != prev, true
	}
	if !active {
		return false, false
	}
	self.Closed = !self.Closed
	return true, false
}

// ReleaseMomentary handles release momentary.
func (self *Switch) ReleaseMomentary() bool {
	if !self.Momentary || !self.Pressed {
		return false
	}
	self.Pressed = false
	return true
}

// effectiveClosed handles effective closed.
func (self *Switch) effectiveClosed() bool {
	if self.Momentary {
		return self.Pressed
	}
	return self.Closed
}
