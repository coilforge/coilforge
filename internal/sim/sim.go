package sim

// File overview:
// sim advances run-mode logic by deriving nets and ticking part simulation interfaces.
// Subsystem: simulation engine.
// It reads world/part contracts and flatten output while staying independent from editor.
// Run mode: a background goroutine steps SimTimeMicros by StepMicros and calls resolveAndTick; UI rate is decoupled.
// Optional pacing ([world.SimFullSpeed] off): smoothed deadline scheduling targets SimTargetTicksPerSec iterations per wall second.
// Flow position: run pipeline for world state; app starts/stops the background loop on mode change.

import (
	"coilforge/internal/core"
	"coilforge/internal/flatten"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"math/rand"
	"sync"
	"time"
)

// StepMicros is the simulation clock quantum (re-export of [core.SimStepMicros]).
const StepMicros = core.SimStepMicros

// maxResolveIterations defines a package-level constant.
const maxResolveIterations = 8

const (
	minSimTargetTicksPerSec = 1
	maxSimTargetTicksPerSec = 1_000_000
	paceLateResync          = 200 * time.Millisecond
)

// simRand stores package-level state.
var simRand = rand.New(rand.NewSource(1))

var (
	simLoopStop chan struct{}
	simLoopWG   sync.WaitGroup
	activeMomentary part.InputHandler
)

// LoopBegin starts a goroutine that advances simulated time continuously until [LoopEnd].
func LoopBegin() {
	if simLoopStop != nil {
		return
	}
	simLoopStop = make(chan struct{})
	simLoopWG.Add(1)
	go simBackgroundLoop()
}

func simBackgroundLoop() {
	defer simLoopWG.Done()
	var paceNext time.Time
	var wasPaced bool
	for {
		select {
		case <-simLoopStop:
			return
		default:
		}
		paced := !world.SimFullSpeed && world.SimTargetTicksPerSec > 0
		if paced && !wasPaced {
			paceNext = time.Time{}
		}
		wasPaced = paced

		world.SimMu.Lock()
		world.SimTimeMicros += StepMicros
		resolveAndTick()
		world.SimMu.Unlock()

		if !paced {
			continue
		}

		tps := clampSimTargetTicks(world.SimTargetTicksPerSec)
		period := time.Second / time.Duration(tps)
		now := time.Now()
		if paceNext.IsZero() {
			paceNext = now
		}
		paceNext = paceNext.Add(period)
		d := time.Until(paceNext)
		if d > 0 {
			time.Sleep(d)
		} else if -d > paceLateResync {
			paceNext = now
		}
	}
}

func clampSimTargetTicks(t int) int {
	if t < minSimTargetTicksPerSec {
		return minSimTargetTicksPerSec
	}
	if t > maxSimTargetTicksPerSec {
		return maxSimTargetTicksPerSec
	}
	return t
}

// LoopEnd stops the background loop started by [LoopBegin] and waits for it to exit.
func LoopEnd() {
	if simLoopStop == nil {
		return
	}
	close(simLoopStop)
	simLoopWG.Wait()
	simLoopStop = nil
}

// Start starts its work.
func Start() {
	part.ClearAllSchematicRuntime(world.Parts)
	flatten.BuildNets()
	world.SimTimeMicros = 0
	resolveAndTick()
}

// Stop stops its work.
func Stop() {
	world.Nets = nil
	world.NetStates = nil
	world.PinNet = nil
	world.SimTimeMicros = 0
}

// HandleClick handles click.
func HandlePointerDown(pt core.Pt) {
	for i := len(world.Parts) - 1; i >= 0; i-- {
		hit := world.Parts[i].HitTest(pt)
		if !hit.Hit {
			continue
		}
		if handler, ok := world.Parts[i].(part.InputHandler); ok {
			changed, momentary := handler.HandleInput(true)
			if momentary {
				activeMomentary = handler
			}
			if changed {
				resolveAndTick()
			}
		}
		return
	}
	activeMomentary = nil
}

// HandlePointerUp releases a pressed momentary input, if any.
func HandlePointerUp() {
	if activeMomentary == nil {
		return
	}
	if activeMomentary.ReleaseMomentary() {
		resolveAndTick()
	}
	activeMomentary = nil
}

// HandleClick keeps backward compatibility with existing callers.
func HandleClick(pt core.Pt) {
	HandlePointerDown(pt)
}

// resolveAndTick handles resolve and tick.
func resolveAndTick() {
	for iterations := 0; iterations < maxResolveIterations; iterations++ {
		resolveNets()

		ctx := part.SimContext{
			NetByPin:     netByPin,
			NetState:     netStateLookup,
			NowMicros:    world.SimTimeMicros,
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
			seeder.SeedNets(union, netByPin, high, low, world.SimTimeMicros)
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
	// Propagate directed drives (e.g. diode ANODE -> CATHODE) until stable.
	// Roots are evaluated after conductive unioning so edges and seeds compose.
	for changed := true; changed; {
		changed = false
		for _, edge := range graph.Edges {
			from := union.Find(edge.FromNet)
			to := union.Find(edge.ToNet)
			if from < 0 || to < 0 {
				continue
			}
			switch edge.Drive {
			case core.NetHigh:
				if high[from] && !high[to] {
					high[to] = true
					changed = true
				}
			case core.NetLow:
				if low[from] && !low[to] {
					low[to] = true
					changed = true
				}
			}
		}
	}

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
