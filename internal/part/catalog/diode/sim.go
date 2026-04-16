package diode

// File overview:
// sim implements simulation-facing behavior for diode using part sim interfaces.
// Subsystem: part catalog (diode) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// AddStateEdges adds state edges.
func (d *Diode) AddStateEdges(netByPin func(core.PinID) int, graph *part.StateGraph) {
	from := netByPin(d.PinAnode)
	to := netByPin(d.PinCathode)
	if from < 0 || to < 0 {
		return
	}
	graph.Edges = append(graph.Edges, part.StateEdge{
		FromNet: from,
		ToNet:   to,
		Drive:   core.NetHigh,
	})
}
