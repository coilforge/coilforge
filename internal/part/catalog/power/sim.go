package power

// File overview:
// sim implements simulation-facing behavior for power using part sim interfaces.
// Subsystem: part catalog (power) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: part sim logic executed during run-mode ticks and net solving.

import "coilforge/internal/core"

// SeedNets seeds nets.
func (p *Power) SeedNets(netByPin func(core.PinID) int, high, low map[int]bool) {
	net := netByPin(p.Pin)
	if net < 0 {
		return
	}
	if p.Kind == "gnd" {
		low[net] = true
		return
	}
	high[net] = true
}
