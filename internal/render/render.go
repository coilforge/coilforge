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

	"github.com/hajimehoshi/ebiten/v2"
)

// DrawScene draws scene.
func DrawScene(dst *ebiten.Image) {
	drawGrid(dst)

	ctx := part.DrawContext{
		Dst:     dst,
		Cam:     world.Cam,
		Zoom:    world.Zoom,
		ScreenW: world.ScreenW,
		ScreenH: world.ScreenH,
	}

	if world.RunMode {
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

// drawGrid handles draw grid.
func drawGrid(dst *ebiten.Image) {
	_, _ = world.WorldToScreen(world.Cam)
	_ = GridColor()
	_ = dst
}
