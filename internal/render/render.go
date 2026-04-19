package render

// File overview:
// render composes scene draw order for schematic content and UI chrome.
// Subsystem: render orchestration.
// It calls part draw methods, theme definitions, and chrome helpers from app.Draw.
// Flow position: visual output stage after world/editor/sim state updates.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawScene draws scene.
// In run mode, SimMu is held only while parts are drawn (NetState reads PinNet/NetStates).
// Chrome and other UI draw outside this path so the sim goroutine can advance SimTimeMicros during rasterization.
func DrawScene(dst *ebiten.Image) {
	fillSchematicBackground(dst)
	drawGrid(dst)

	ctx := part.DrawContext{
		Dst:      dst,
		Cam:      world.Cam,
		Zoom:     world.Zoom,
		ScreenW:  world.ScreenW,
		ScreenH:  world.ScreenH,
		DarkMode: DarkMode,
	}

	if world.RunMode {
		world.SimMu.RLock()
		defer world.SimMu.RUnlock()
		ctx.NetState = func(pinID core.PinID) int {
			netID, ok := world.PinNet[pinID]
			if !ok {
				return core.NetFloat
			}
			if state, ok := world.NetStates[netID]; ok {
				return state
			}
			return core.NetFloat
		}
	}

	for _, p := range world.Parts {
		p.Draw(ctx)
	}
}

func fillSchematicBackground(dst *ebiten.Image) {
	b := dst.Bounds()
	vector.FillRect(
		dst,
		float32(b.Min.X),
		float32(b.Min.Y),
		float32(b.Dx()),
		float32(b.Dy()),
		SchematicBackgroundColor(),
		false,
	)
}

// drawGrid draws world-space minor (wire) and major (part pitch) grid lines in screen space.
func drawGrid(dst *ebiten.Image) {
	w, h := world.ScreenW, world.ScreenH
	if w <= 0 || h <= 0 {
		return
	}

	minor := world.MinorGridWorld
	major := world.MajorGridWorld
	ratio := int(math.Round(major / minor))
	if ratio < 1 {
		ratio = 1
	}

	minX, maxX, minY, maxY := visibleWorldExtent(w, h)

	i0 := int(math.Floor(minX / minor))
	i1 := int(math.Ceil(maxX / minor))
	j0 := int(math.Floor(minY / minor))
	j1 := int(math.Ceil(maxY / minor))

	runSim := world.RunMode
	minorCol := GridMinorColor()
	majorCol := GridMajorColor()
	if runSim {
		majorCol = GridMajorColorRunMode()
	}

	const swMinor float32 = 1
	const swMajor float32 = 1.85

	for i := i0; i <= i1; i++ {
		x := float64(i) * minor
		modI := i % ratio
		if modI < 0 {
			modI += ratio
		}
		isMajor := modI == 0
		if runSim && !isMajor {
			continue
		}
		col := minorCol
		sw := swMinor
		if isMajor {
			col = majorCol
			sw = swMajor
		}
		x0, y0 := world.WorldToScreen(core.Pt{X: x, Y: minY})
		x1, y1 := world.WorldToScreen(core.Pt{X: x, Y: maxY})
		vector.StrokeLine(dst, float32(x0), float32(y0), float32(x1), float32(y1), sw, col, false)
	}

	for j := j0; j <= j1; j++ {
		y := float64(j) * minor
		modJ := j % ratio
		if modJ < 0 {
			modJ += ratio
		}
		isMajor := modJ == 0
		if runSim && !isMajor {
			continue
		}
		col := minorCol
		sw := swMinor
		if isMajor {
			col = majorCol
			sw = swMajor
		}
		x0, y0 := world.WorldToScreen(core.Pt{X: minX, Y: y})
		x1, y1 := world.WorldToScreen(core.Pt{X: maxX, Y: y})
		vector.StrokeLine(dst, float32(x0), float32(y0), float32(x1), float32(y1), sw, col, false)
	}
}

func visibleWorldExtent(screenW, screenH int) (minX, maxX, minY, maxY float64) {
	corners := []core.Pt{
		world.ScreenToWorld(0, 0),
		world.ScreenToWorld(screenW, 0),
		world.ScreenToWorld(0, screenH),
		world.ScreenToWorld(screenW, screenH),
	}
	minX = corners[0].X
	maxX = corners[0].X
	minY = corners[0].Y
	maxY = corners[0].Y
	for k := 1; k < len(corners); k++ {
		p := corners[k]
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return minX, maxX, minY, maxY
}
