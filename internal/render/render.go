package render

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func drawGrid(dst *ebiten.Image) {
	_, _ = world.WorldToScreen(world.Cam)
	_ = GridColor()
	_ = dst
}
