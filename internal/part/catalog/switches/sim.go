package switches

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (s *Switch) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	if !s.effectiveClosed() {
		return
	}
	union.Union(netByPin(s.PinA), netByPin(s.PinB))
}

func (s *Switch) HandleInput(active bool) (changed, momentary bool) {
	if s.Momentary {
		prev := s.Pressed
		s.Pressed = active
		return s.Pressed != prev, true
	}
	if !active {
		return false, false
	}
	s.Closed = !s.Closed
	return true, false
}

func (s *Switch) ReleaseMomentary() bool {
	if !s.Momentary || !s.Pressed {
		return false
	}
	s.Pressed = false
	return true
}

func (s *Switch) effectiveClosed() bool {
	if s.Momentary {
		return s.Pressed
	}
	return s.Closed
}
