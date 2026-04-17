package diode

// File overview:
// sim implements simulation-facing behavior using part sim interfaces.
// Subsystem: part catalog simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// AddStateEdges adds state edges.
func (self *Diode) AddStateEdges(netByPin func(core.PinID) int, graph *part.StateGraph) {
	from := netByPin(self.PinAnode)
	to := netByPin(self.PinCathode)
	if from < 0 || to < 0 {
		return
	}
	graph.Edges = append(graph.Edges, part.StateEdge{
		FromNet: from,
		ToNet:   to,
		Drive:   core.NetHigh,
	})
}
