package indicator

// File overview:
// sim implements simulation-facing behavior for indicator using part sim interfaces.
// Subsystem: part catalog (indicator) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// Tick handles tick.
// Pin fields TerminalA / TerminalB come from generated IndicatorPinIDs (SVG ids); compare nets or states as needed.
func (self *Indicator) Tick(ctx part.SimContext) bool {
	wasLit := self.Lit
	a := ctx.PinNetState(self.TerminalA)
	b := ctx.PinNetState(self.TerminalB)
	if a != core.NetFloat && b != core.NetFloat && a != b {
		self.Lit = true
	} else {
		self.Lit = false
	}
	return self.Lit != wasLit
}
