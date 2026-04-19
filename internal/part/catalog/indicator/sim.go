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
	net := ctx.NetByPin(self.TerminalA)
	wasLit := self.Lit
	self.Lit = ctx.NetState(net) == core.NetHigh
	return self.Lit != wasLit
}
