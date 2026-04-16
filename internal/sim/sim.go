package sim

// File overview:
// sim advances run-mode logic by deriving nets and ticking part simulation interfaces.
// Subsystem: simulation engine.
// It reads world/part contracts and flatten output while staying independent from editor.
// Flow position: run pipeline branch invoked by app during active simulation mode.

import (
	"coilforge/internal/core"
	"coilforge/internal/flatten"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"math/rand"
)

// TickMicros defines a package-level constant.
const TickMicros = 10

// maxResolveIterations defines a package-level constant.
const maxResolveIterations = 8

// simRand stores package-level state.
var simRand = rand.New(rand.NewSource(1))

// Start starts its work.
func Start() {
	flatten.BuildNets()
	world.SimTick = 0
	resolveAndTick()
}

// Stop stops its work.
func Stop() {
	world.Nets = nil
	world.NetStates = nil
	world.PinNet = nil
	world.SimTick = 0
}

// AdvanceFrame handles advance frame.
func AdvanceFrame() {
	target := world.SimTick + 1000
	for world.SimTick < target {
		world.SimTick = nextInterestingTick(target)
		resolveAndTick()
	}
}

// HandleClick handles click.
func HandleClick(pt core.Pt) {
	for i := len(world.Parts) - 1; i >= 0; i-- {
		hit := world.Parts[i].HitTest(pt)
		if !hit.Hit {
			continue
		}
		if handler, ok := world.Parts[i].(part.InputHandler); ok {
			changed, _ := handler.HandleInput(true)
			if changed {
				resolveAndTick()
			}
		}
		return
	}
}

// resolveAndTick handles resolve and tick.
func resolveAndTick() {
	for iterations := 0; iterations < maxResolveIterations; iterations++ {
		resolveNets()

		ctx := part.SimContext{
			NetByPin:     netByPin,
			NetState:     netStateLookup,
			Tick:         world.SimTick,
			TickMicros:   TickMicros,
			EnableJitter: true,
			Rand:         simRand,
		}

		changed := false
		for _, p := range world.Parts {
			if sp, ok := p.(part.SimPart); ok && sp.Tick(ctx) {
				changed = true
			}
		}

		if !changed {
			return
		}
	}
}

// resolveNets handles resolve nets.
func resolveNets() {
	union := newUnionFind(len(world.Nets))
	high := map[int]bool{}
	low := map[int]bool{}

	for _, p := range world.Parts {
		if conductor, ok := p.(part.Conductor); ok {
			conductor.AddConductive(union, netByPin)
		}
	}

	for _, p := range world.Parts {
		if seeder, ok := p.(part.NetSeeder); ok {
			seeder.SeedNets(netByPin, high, low)
		}
	}

	var graph part.StateGraph
	for _, p := range world.Parts {
		if edger, ok := p.(part.StateEdger); ok {
			edger.AddStateEdges(netByPin, &graph)
		}
	}

	world.NetStates = resolveFromSeeds(union, high, low, &graph)
}

// resolveFromSeeds handles resolve from seeds.
func resolveFromSeeds(union *unionFind, high, low map[int]bool, graph *part.StateGraph) map[int]int {
	_ = graph

	out := make(map[int]int, len(world.Nets))
	for _, net := range world.Nets {
		root := union.Find(net.ID)
		switch {
		case high[root] && low[root]:
			out[net.ID] = core.NetShort
		case high[root]:
			out[net.ID] = core.NetHigh
		case low[root]:
			out[net.ID] = core.NetLow
		default:
			out[net.ID] = core.NetFloat
		}
	}
	return out
}

// nextInterestingTick handles next interesting tick.
func nextInterestingTick(target uint64) uint64 {
	return target
}

// netByPin handles net by pin.
func netByPin(pinID core.PinID) int {
	if netID, ok := world.PinNet[pinID]; ok {
		return netID
	}
	return -1
}

// netStateLookup handles net state lookup.
func netStateLookup(netID int) int {
	if state, ok := world.NetStates[netID]; ok {
		return state
	}
	return core.NetFloat
}
