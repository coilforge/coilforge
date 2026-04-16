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
func (s *Switch) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
	if !s.effectiveClosed() {
		return
	}
	union.Union(netByPin(s.PinA), netByPin(s.PinB))
}

// HandleInput handles input.
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

// ReleaseMomentary handles release momentary.
func (s *Switch) ReleaseMomentary() bool {
	if !s.Momentary || !s.Pressed {
		return false
	}
	s.Pressed = false
	return true
}

// effectiveClosed handles effective closed.
func (s *Switch) effectiveClosed() bool {
	if s.Momentary {
		return s.Pressed
	}
	return s.Closed
}
