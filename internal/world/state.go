package world

// File overview:
// state stores shared schematic runtime data for parts, camera, nets, and mode flags.
// Subsystem: world shared state.
// It is read and mutated by app/editor/sim/render through package-level access.
// Flow position: common state hub linking otherwise independent subsystems.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"sync"
)

// SimMu guards simulation mutations (nets, net states, part tick state) against concurrent reads during Draw.
// Writers: background sim loop and run-mode click handling. Readers: Draw while RunMode is true.
var SimMu sync.RWMutex

// Parts stores package-level state.
var Parts []part.Part

// NextPartID stores package-level state.
var NextPartID int

// NextPinID stores package-level state.
var NextPinID core.PinID

// Cam stores package-level state.
var Cam core.Pt

// defaultZoom is pixels per world unit at startup; with SVGUserUnitToWorld 1/16, a full 512-wide
// symbol spans ~128 screen pixels so strokes stay visible without wheel zoom.
const defaultZoom = 4.0

// Zoom stores package-level state.
var Zoom = defaultZoom

// ZoomMin and ZoomMax clamp keyboard/wheel zoom so coordinates stay sane.
const (
	ZoomMin = 0.25
	ZoomMax = 128.0
)

// ClampZoom clamps [Zoom] to [ZoomMin, ZoomMax].
func ClampZoom() {
	if Zoom < ZoomMin {
		Zoom = ZoomMin
	}
	if Zoom > ZoomMax {
		Zoom = ZoomMax
	}
}

// ScreenW stores package-level state.
var ScreenW int

// ScreenH stores package-level state.
var ScreenH int

// RunMode stores package-level state.
var RunMode bool

// Nets stores package-level state.
var Nets []core.Net

// NetStates stores package-level state.
var NetStates map[int]int

// PinNet stores package-level state.
var PinNet map[core.PinID]int

// SimTimeMicros is monotonic simulated time since run-mode sim start, in microseconds.
var SimTimeMicros uint64

// AllocPartID handles alloc part id.
func AllocPartID() int {
	id := NextPartID
	NextPartID++
	return id
}

// AllocPinID handles alloc pin id.
func AllocPinID() core.PinID {
	id := NextPinID
	NextPinID++
	return id
}

// ScreenToWorld handles screen to world.
func ScreenToWorld(sx, sy int) core.Pt {
	return core.Pt{
		X: (float64(sx)-float64(ScreenW)/2)/Zoom + Cam.X,
		Y: (float64(sy)-float64(ScreenH)/2)/Zoom + Cam.Y,
	}
}

// WorldToScreen handles world to screen.
func WorldToScreen(pt core.Pt) (float64, float64) {
	return (pt.X-Cam.X)*Zoom + float64(ScreenW)/2,
		(pt.Y-Cam.Y)*Zoom + float64(ScreenH)/2
}

// Reset resets its work.
func Reset() {
	Parts = nil
	NextPartID = 0
	NextPinID = 0
	Cam = core.Pt{}
	Zoom = defaultZoom
	RunMode = false
	Nets = nil
	NetStates = nil
	PinNet = nil
	SimTimeMicros = 0
}
