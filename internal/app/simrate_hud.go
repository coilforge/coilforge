package app

// File overview:
// simrate_hud tracks simulated time vs wall time and formats a bottom-right HUD string.
// Subsystem: app (orchestration). Uses sim.StepMicros and world.SimTimeMicros.

import (
	"fmt"
	"time"

	"coilforge/internal/sim"
	"coilforge/internal/world"
)

// updateSimRealtimeHUD samples world.SimTimeMicros vs wall clock and EMA-smoothes the ratio
// (simulated μs advanced per real μs). 1.0 = realtime; 2.0 = 200% of realtime, etc.
func (a *App) updateSimRealtimeHUD() {
	if !world.RunMode {
		a.simRTHasSample = false
		a.simRTSmoothInit = false
		a.simRTPaceKnown = false
		return
	}
	if !a.simRTPaceKnown {
		a.simRTPaceLast = world.SimFullSpeed
		a.simRTPaceKnown = true
	} else if world.SimFullSpeed != a.simRTPaceLast {
		a.simRTHasSample = false
		a.simRTSmoothInit = false
		a.simRTPaceLast = world.SimFullSpeed
	}
	world.SimMu.RLock()
	simNow := world.SimTimeMicros
	world.SimMu.RUnlock()
	wallNow := time.Now()
	if !a.simRTHasSample {
		a.simRTLastWall = wallNow
		a.simRTLastSim = simNow
		a.simRTHasSample = true
		return
	}
	dtWall := wallNow.Sub(a.simRTLastWall).Microseconds()
	dtSim := int64(simNow) - int64(a.simRTLastSim)
	a.simRTLastWall = wallNow
	a.simRTLastSim = simNow
	if dtSim < 0 {
		a.simRTSmoothInit = false
		return
	}
	if dtWall < 5000 {
		return
	}
	if dtWall <= 0 {
		return
	}
	instant := float64(dtSim) / float64(dtWall)
	const emaAlpha = 0.55
	if !a.simRTSmoothInit {
		a.simRTRatio = instant
		a.simRTSmoothInit = true
	} else {
		a.simRTRatio = (1-emaAlpha)*a.simRTRatio + emaAlpha*instant
	}
}

// simRealtimeHUDText is a short line for the lower-right corner: rate vs realtime and step quantum.
func (a *App) simRealtimeHUDText() string {
	step := sim.StepMicros
	// TEMP: show world.Zoom while tuning zoom limits / grid visibility; remove when done.
	z := fmt.Sprintf("  zoom=%.4g", world.Zoom)
	if !world.RunMode {
		return fmt.Sprintf("Sim --  %dus step%s", step, z)
	}
	pace := simPaceHUDLabel()
	if !a.simRTSmoothInit {
		return fmt.Sprintf("Sim ...%s  %dus step%s", pace, step, z)
	}
	return fmt.Sprintf("Sim %s%s  %dus step%s", formatSimRTPercent(a.simRTRatio), pace, step, z)
}

// simPaceHUDLabel reports full-speed vs paced target for the HUD (F8 toggles while running).
func simPaceHUDLabel() string {
	if world.SimFullSpeed || world.SimTargetTicksPerSec <= 0 {
		return " full"
	}
	return fmt.Sprintf(" %dt/s", world.SimTargetTicksPerSec)
}

// formatSimRTPercent turns ratio (sim/wall, 1=realtime) into a readable percent string.
func formatSimRTPercent(ratio float64) string {
	p := ratio * 100
	switch {
	case p >= 1000:
		return fmt.Sprintf("%.0f%% RT", p)
	case p >= 10:
		return fmt.Sprintf("%.1f%% RT", p)
	case p >= 0.01:
		return fmt.Sprintf("%.3f%% RT", p)
	default:
		return fmt.Sprintf("%.4g%% RT", p)
	}
}
