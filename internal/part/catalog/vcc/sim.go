package vcc

// File overview:
// sim implements NetSeeder for VCC (drive attached net high).
// Subsystem: part catalog (vcc) simulation.
// It is invoked by the sim engine through interfaces and does not depend on editor.
// Flow position: net seeding during resolveNets.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

// SeedNets marks the union root for this part's terminal as a forced-high seed.
func (self *Vcc) SeedNets(union part.NetUnion, netByPin func(core.PinID) int, high, low map[int]bool, nowMicros uint64) {
	_ = nowMicros
	netID := netByPin(self.VCC)
	if netID < 0 {
		return
	}
	root := union.Find(netID)
	if root < 0 {
		return
	}
	high[root] = true
}
