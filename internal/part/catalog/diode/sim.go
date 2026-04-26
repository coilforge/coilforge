package diode

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

func (self *Diode) Tick(ctx part.SimContext) bool {
	_ = ctx
	return false
}

func (self *Diode) AddStateEdges(netByPin func(core.PinID) int, graph *part.StateGraph) {
	anode := netByPin(self.ANODE)
	cathode := netByPin(self.CATHODE)
	if anode < 0 || cathode < 0 {
		return
	}
	graph.Edges = append(graph.Edges,
		part.StateEdge{FromNet: anode, ToNet: cathode, Drive: core.NetHigh},
		part.StateEdge{FromNet: anode, ToNet: cathode, Drive: core.NetLow},
	)
}
