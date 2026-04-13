package diode

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

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
