package editor

// File overview:
// viewport handles Space+drag camera pan over the schematic (screen-space deltas → world.Cam).

import (
	"coilforge/internal/world"
)

// BeginViewportPan starts a Space+drag pan gesture (primary button held).
func BeginViewportPan(screenX, screenY int) {
	ViewportPanDrag = true
	PanLastScreenX = screenX
	PanLastScreenY = screenY
}

// ViewportPanActive is true while Space+drag pan is in progress.
func ViewportPanActive() bool {
	return ViewportPanDrag
}

// HandleViewportPanDrag applies incremental pan from screen-pixel movement.
func HandleViewportPanDrag(screenX, screenY int) {
	if !ViewportPanDrag {
		return
	}
	dsx := screenX - PanLastScreenX
	dsy := screenY - PanLastScreenY
	world.Cam.X -= float64(dsx) / world.Zoom
	world.Cam.Y -= float64(dsy) / world.Zoom
	PanLastScreenX = screenX
	PanLastScreenY = screenY
}

func endViewportPan() {
	ViewportPanDrag = false
}
