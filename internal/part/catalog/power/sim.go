package power

import "coilforge/internal/core"

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
