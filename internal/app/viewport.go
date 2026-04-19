package app

// File overview:
// viewport applies wheel/trackpad scrolling to zoom (vertical) and pan (horizontal) in world space.

import (
	"coilforge/internal/world"
)

const (
	wheelZoomPerUnit = 0.06 // vertical wheel → multiplicative zoom per frame unit (trackpad-friendly).
	wheelPanPerUnit  = 0.85 // horizontal wheel → world pan scale before dividing by Zoom.
)

func handleViewportWheel(screenX, screenY int, wheelX, wheelY float64) {
	if wheelX != 0 {
		world.Cam.X -= wheelX * wheelPanPerUnit / world.Zoom
	}
	if wheelY != 0 {
		zoomWorldAtScreen(screenX, screenY, wheelY)
	}
}

func zoomWorldAtScreen(screenX, screenY int, wheelY float64) {
	if wheelY == 0 {
		return
	}
	before := world.ScreenToWorld(screenX, screenY)
	world.Zoom *= 1.0 + wheelY*wheelZoomPerUnit
	world.ClampZoom()
	after := world.ScreenToWorld(screenX, screenY)
	world.Cam.X += before.X - after.X
	world.Cam.Y += before.Y - after.Y
}
