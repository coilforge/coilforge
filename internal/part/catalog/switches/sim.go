package switches

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Switches) AddConductive(union part.NetUnion, netByPin func(pinID core.PinID) int) {
	if !self.On {
		return
	}
	a := netByPin(self.A)
	b := netByPin(self.B)
	if a < 0 || b < 0 {
		return
	}
	union.Union(a, b)
}

func (self *Switches) Tick(ctx part.SimContext) bool {
	_ = ctx
	// Input transitions are handled by HandleInput/ReleaseMomentary (mouse down/up).
	return false
}

func (self *Switches) HandleInput(active bool) (changed, momentary bool) {
	if self.Mode == ModeMomentary {
		next := active
		if self.On == next {
			return false, true
		}
		self.On = next
		return true, true
	}
	self.On = !self.On
	return true, false
}

func (self *Switches) ReleaseMomentary() bool {
	if self.Mode != ModeMomentary || !self.On {
		return false
	}
	self.On = false
	return true
}
