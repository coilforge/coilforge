package wire

// File overview:
// sim implements simulation-facing behavior for wire using part sim interfaces.
// Subsystem: part catalog (wire) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import "coilforge/internal/part"

// Tick handles tick.
func (w *Wire) Tick(ctx part.SimContext) bool {
	prev := w.State
	if net := ctx.NetByPin(w.PinA); net >= 0 {
		w.State = ctx.NetState(net)
	}
	return w.State != prev
}
