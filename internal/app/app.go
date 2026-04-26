package app

// File overview:
// app coordinates frame-level input, mode routing, and top-level draw calls.
// Subsystem: app orchestration.
// It calls editor and sim independently while sharing state through world and render.
// Flow position: primary runtime controller between Ebiten loop and subsystems.

import (
	"math"
	"time"

	"coilforge/internal/appsettings"
	"coilforge/internal/editor"
	"coilforge/internal/render"
	"coilforge/internal/storage"
	"coilforge/internal/uidebug"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type App struct {
	leftDown       bool // Tracks previous-frame left mouse state for edge detection.
	hoverLeftTool  int  // Hovered placement tool index for left toolbar.
	hoverRightTool int  // Hovered command tool index for right toolbar.
	toolbarCapture bool // True while current mouse press started on toolbar chrome.
	settingsOpen   bool // True while the app settings panel is visible.
	settingsPath   string
	settingsPathActive bool
	docDialog      docDialogState

	simRTLastWall   time.Time // Wall clock sample for sim vs realtime HUD.
	simRTLastSim    uint64    // SimTimeMicros sample paired with simRTLastWall.
	simRTRatio      float64   // Smoothed (sim μs / wall μs); 1.0 = realtime.
	simRTHasSample  bool      // True after first paired sample in this run session.
	simRTSmoothInit bool      // True after first EMA update (valid simRTRatio for display).
	simRTPaceKnown  bool      // True after recording SimFullSpeed for HUD reset on F8.
	simRTPaceLast   bool      // Last seen world.SimFullSpeed while in run mode.
}

// windowTPS is how many times per second Ebiten calls Update (and typically Draw).
// Keep this well above ~20 so short double-clicks and other press/release edges are not missed
// between samples (10 TPS was too low: a second click could land in the same 100ms window with no extra frame).
// Simulation runs in a separate goroutine; this only controls input polling and redraw cadence.
const windowTPS = 60

// New constructs a fresh application instance.
func New() *App {
	store := storage.NewLocalFSStore(appsettings.Current.DefaultSaveDir)
	return &App{
		hoverLeftTool:  -1,
		hoverRightTool: -1,
		settingsPath:   appsettings.Current.DefaultSaveDir,
		settingsPathActive: true,
		docDialog: docDialogState{
			mode:     docDialogClosed,
			store:    store,
			selected: -1,
		},
	}
}

// Run resets world state, configures the window, and starts the Ebiten loop.
func Run() error {
	world.Reset()
	if loaded, err := appsettings.LoadLocal(); err == nil {
		appsettings.Current = loaded
	}
	render.DarkMode = appsettings.Current.DarkMode
	ebiten.SetWindowSize(1440, 900)
	ebiten.SetWindowTitle("CoilForge")
	ebiten.SetTPS(windowTPS)
	return ebiten.RunGame(New())
}

// Update runs one frame of app orchestration and input dispatch.
// It refreshes viewport size, processes keyboard/mouse input, and routes input
// into editor or simulation behavior depending on run mode.
// Run-mode simulation advances in sim.LoopBegin's background goroutine, not here.
func (a *App) Update() error {
	winW, winH := ebiten.WindowSize()
	sw, sh := a.LayoutF(float64(winW), float64(winH))
	world.ScreenW = int(math.Ceil(sw))
	world.ScreenH = int(math.Ceil(sh))
	mx, my := ebiten.CursorPosition()
	a.updateToolbarHover(mx, my)
	if !world.RunMode {
		pointerX, pointerY, ok := toolbarPointerPosition(mx, my)
		if ok {
			editor.UpdatePlacementPreview(world.ScreenToWorld(pointerX, pointerY))
		}
	}

	// Modal doc dialog: consume all keyboard/mouse input while open.
	if a.docDialog.mode != docDialogClosed {
		a.handleDocDialogTyping()
		a.handleMouse(mx, my)
		a.updateSimRealtimeHUD()
		uidebug.LogUpdateFrame()
		return nil
	}
	// Settings is keyboard-editable; keep input focused here while open.
	if a.settingsOpen {
		a.handleSettingsTyping()
		a.handleMouse(mx, my)
		a.updateSimRealtimeHUD()
		uidebug.LogUpdateFrame()
		return nil
	}

	a.handleToolHotkeys()
	a.handleEditorHotkeys()
	a.handleProjectHotkeys()
	a.handleDocDialogTyping()
	a.handleZoomHotkeys()
	a.handleMouse(mx, my)

	if wx, wy := ebiten.Wheel(); wx != 0 || wy != 0 {
		handleViewportWheel(mx, my, wx, wy)
	}

	a.updateSimRealtimeHUD()
	uidebug.LogUpdateFrame()
	return nil
}

// Draw composes scene rendering, editor overlays, and screen-space chrome.
func (a *App) Draw(screen *ebiten.Image) {
	render.DrawScene(screen)

	if !world.RunMode {
		editor.DrawOverlays(screen)
		render.DrawToolbar(screen, render.ToolbarLeft, toolbarButtons(), activeToolIndex(), a.hoverLeftTool)
	}
	// Command strip: visible in edit and run mode.
	render.DrawToolbar(screen, render.ToolbarRight, rightToolbarButtons(), -1, a.hoverRightTool)
	if selectedPart := selectedPart(); selectedPart != nil {
		render.DrawPropPanel(screen, selectedPart.PropSpec())
	}
	if a.settingsOpen {
		render.DrawSettingsPanel(screen, "Settings", a.settingsPanelRows(), a.settingsPanelFooter())
	}
	if a.docDialog.mode != docDialogClosed {
		render.DrawDocBrowserPanel(screen, a.docDialogTitle(), a.docDialog.input, a.docDialogRows(), a.docDialogFooter())
	}
	render.DrawStatusBar(screen, a.statusText())
	render.DrawSimRealtimeHUD(screen, a.simRealtimeHUDText())
}

// Layout satisfies [ebiten.Game] (used when LayoutF is unavailable).
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	wf, hf := a.LayoutF(float64(outsideWidth), float64(outsideHeight))
	return int(math.Ceil(wf)), int(math.Ceil(hf))
}

// LayoutF scales the framebuffer by the monitor device scale factor so HiDPI / Retina draws at native
// resolution instead of upscaling a low-res buffer (fixes soft UI and vector strokes).
func (a *App) LayoutF(outsideWidth, outsideHeight float64) (float64, float64) {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		return outsideWidth, outsideHeight
	}
	s := 1.0
	if m := ebiten.Monitor(); m != nil {
		if f := m.DeviceScaleFactor(); f > 0 {
			s = f
		}
	}
	return outsideWidth * s, outsideHeight * s
}
