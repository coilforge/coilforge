package rch

// File overview:
// sim implements simulation-facing behavior for rch using part sim interfaces.
// Subsystem: part catalog (rch) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import "coilforge/internal/part"

// Tick handles tick.
func (self *RCH) Tick(ctx part.SimContext) bool {
	_ = ctx
	return false
}
