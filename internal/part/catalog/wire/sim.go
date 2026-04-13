package wire

import "coilforge/internal/part"

func (w *Wire) Tick(ctx part.SimContext) bool {
	prev := w.State
	if net := ctx.NetByPin(w.PinA); net >= 0 {
		w.State = ctx.NetState(net)
	}
	return w.State != prev
}
