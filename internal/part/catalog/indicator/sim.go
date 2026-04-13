package indicator

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (ind *Indicator) Tick(ctx part.SimContext) bool {
	net := ctx.NetByPin(ind.PinA)
	wasLit := ind.Lit
	ind.Lit = ctx.NetState(net) == core.NetHigh
	return ind.Lit != wasLit
}
