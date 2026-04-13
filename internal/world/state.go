package world

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

var Parts []part.Part
var NextPartID int
var NextPinID core.PinID

var Cam core.Pt
var Zoom = 1.0
var ScreenW int
var ScreenH int

var RunMode bool

var Nets []core.Net
var NetStates map[int]int
var PinNet map[core.PinID]int
var SimTick uint64

func AllocPartID() int {
	id := NextPartID
	NextPartID++
	return id
}

func AllocPinID() core.PinID {
	id := NextPinID
	NextPinID++
	return id
}

func ScreenToWorld(sx, sy int) core.Pt {
	return core.Pt{
		X: (float64(sx)-float64(ScreenW)/2)/Zoom + Cam.X,
		Y: (float64(sy)-float64(ScreenH)/2)/Zoom + Cam.Y,
	}
}

func WorldToScreen(pt core.Pt) (float64, float64) {
	return (pt.X-Cam.X)*Zoom + float64(ScreenW)/2,
		(pt.Y-Cam.Y)*Zoom + float64(ScreenH)/2
}

func Reset() {
	Parts = nil
	NextPartID = 0
	NextPinID = 0
	Cam = core.Pt{}
	Zoom = 1.0
	RunMode = false
	Nets = nil
	NetStates = nil
	PinNet = nil
	SimTick = 0
}
