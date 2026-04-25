package app

// File overview:
// app coordinates frame-level input, mode routing, and top-level draw calls.
// Subsystem: app orchestration.
// It calls editor and sim independently while sharing state through world and render.
// Flow position: primary runtime controller between Ebiten loop and subsystems.

import (
	"fmt"
	"math"
	"os"
	"time"

	"coilforge/internal/appsettings"
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/partmanifest"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/uidebug"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	leftDown       bool // Tracks previous-frame left mouse state for edge detection.
	hoverLeftTool  int  // Hovered placement tool index for left toolbar.
	hoverRightTool int  // Hovered command tool index for right toolbar.
	toolbarCapture bool // True while current mouse press started on toolbar chrome.
	settingsOpen   bool // True while the app settings panel is visible.

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
	return &App{hoverLeftTool: -1, hoverRightTool: -1}
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

	a.handleToolHotkeys()
	a.handleEditorHotkeys()
	a.handleProjectHotkeys()
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
		render.DrawSimplePanel(screen, "Settings", a.settingsPanelRows())
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

// handleToolHotkeys maps numeric keys to placement tool selection.
func (a *App) handleToolHotkeys() {
	for _, item := range partmanifest.PlacementTools {
		key, ok := hotkeyToEbitenKey(item.Hotkey)
		if ok && inpututil.IsKeyJustPressed(key) {
			editor.StartPlacement(item.TypeID)
		}
	}
}

// hotkeyToEbitenKey maps toolbar hotkey runes to Ebiten key constants.
func hotkeyToEbitenKey(hotkey rune) (ebiten.Key, bool) {
	switch hotkey {
	case '1':
		return ebiten.Key1, true
	case '2':
		return ebiten.Key2, true
	case '3':
		return ebiten.Key3, true
	case '4':
		return ebiten.Key4, true
	case '5':
		return ebiten.Key5, true
	case '6':
		return ebiten.Key6, true
	case '7':
		return ebiten.Key7, true
	case '8':
		return ebiten.Key8, true
	case '9':
		return ebiten.Key9, true
	case '0':
		return ebiten.Key0, true
	case 'w', 'W':
		return ebiten.KeyW, true
	default:
		return 0, false
	}
}

// handleEditorHotkeys runs editor-specific keyboard handlers.
func (a *App) handleEditorHotkeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && a.closeTopmostOverlay() {
		return
	}
	a.handleTransformHotkeys()
	a.handleHistoryHotkeys()
	a.handleLabelHotkeys()
}

// handleTransformHotkeys processes rotate, mirror, delete, and mode keys.
func (a *App) handleTransformHotkeys() {
	for _, key := range []ebiten.Key{
		ebiten.KeyR,
		ebiten.KeyM,
		ebiten.KeyDelete,
		ebiten.KeyBackspace,
	} {
		if inpututil.IsKeyJustPressed(key) && !world.RunMode {
			editor.HandleKey(key)
		}
	}
}

// closeTopmostOverlay closes the topmost closable app-level UI and reports whether it handled Esc.
func (a *App) closeTopmostOverlay() bool {
	if a.settingsOpen {
		a.settingsOpen = false
		return true
	}
	return false
}

// handleHistoryHotkeys processes undo/redo and clipboard actions.
func (a *App) handleHistoryHotkeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) && !world.RunMode {
		editor.Undo()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyY) && !world.RunMode {
		editor.Redo()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) && !world.RunMode {
		editor.CopySelected()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) && !world.RunMode {
		editor.Paste(core.Pt{X: 16, Y: 16})
	}
}

// handleLabelHotkeys starts and commits label editing when available.
func (a *App) handleLabelHotkeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyL) && !world.RunMode && len(editor.Selection) > 0 {
		editor.StartLabelEdit(editor.Selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !world.RunMode && editor.LabelEditing {
		editor.CommitLabelEdit()
	}
}

// handleProjectHotkeys processes run-mode and project load/save shortcuts.
func (a *App) handleProjectHotkeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		a.settingsOpen = !a.settingsOpen
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		changed := appsettings.Apply(appsettings.Action{
			Index:    0,
			NewValue: !appsettings.Current.DarkMode,
		})
		if changed {
			_ = appsettings.SaveLocalCurrent()
		}
		a.syncRenderThemeFromSettings()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		ToggleRunMode()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		_ = SaveProject(DefaultProjectPath)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		_ = LoadProject(DefaultProjectPath)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF8) && world.RunMode {
		world.SimFullSpeed = !world.SimFullSpeed
	}
}

// handleZoomHotkeys adjusts camera zoom (+ / −).
func (a *App) handleZoomHotkeys() {
	if editor.LabelEditing {
		return
	}

	const factor = 1.125

	zoomOut := inpututil.IsKeyJustPressed(ebiten.KeyMinus) ||
		inpututil.IsKeyJustPressed(ebiten.KeyNumpadSubtract)

	shift := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
	zoomIn := inpututil.IsKeyJustPressed(ebiten.KeyNumpadAdd) ||
		inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) ||
		(inpututil.IsKeyJustPressed(ebiten.KeyEqual) && shift)

	switch {
	case zoomOut:
		world.Zoom /= factor
		world.ClampZoom()
	case zoomIn:
		world.Zoom *= factor
		world.ClampZoom()
	}
}

// handleMouse routes click, drag, and release events by active mode.
func (a *App) handleMouse(mouseX, mouseY int) {
	leftNow := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	switch {
	case leftNow && !a.leftDown:
		if a.settingsOpen {
			if render.SimplePanelCloseButtonAtScreenPoint(mouseX, mouseY) {
				a.settingsOpen = false
			}
			// While settings is open, consume pointer presses so schematic/editor
			// interactions do not happen behind the panel.
			a.toolbarCapture = true
			break
		}
		if a.handleToolbarPress(mouseX, mouseY) {
			a.toolbarCapture = true
			break
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !editor.LabelEditing {
			editor.BeginViewportPan(mouseX, mouseY)
			break
		}
		pt := world.ScreenToWorld(mouseX, mouseY)
		if world.RunMode {
			world.SimMu.Lock()
			sim.HandleClick(pt)
			world.SimMu.Unlock()
		} else {
			editor.HandleMouseDown(pt, int(ebiten.MouseButtonLeft))
		}
	case leftNow && a.leftDown:
		if a.toolbarCapture {
			break
		}
		if editor.ViewportPanActive() {
			editor.HandleViewportPanDrag(mouseX, mouseY)
			break
		}
		if world.RunMode {
			break
		}
		pt := world.ScreenToWorld(mouseX, mouseY)
		editor.HandleDrag(pt)
	case !leftNow && a.leftDown:
		if a.toolbarCapture {
			a.toolbarCapture = false
			break
		}
		pt := world.ScreenToWorld(mouseX, mouseY)
		editor.HandleMouseUp(pt, int(ebiten.MouseButtonLeft))
	}

	a.leftDown = leftNow
}

// handleToolbarPress applies toolbar click behavior and reports whether a press hit toolbar chrome.
func (a *App) handleToolbarPress(mouseX, mouseY int) bool {
	if !world.RunMode {
		if idx := render.ToolbarButtonAtScreenPoint(render.ToolbarLeft, toolbarButtons(), mouseX, mouseY); idx >= 0 {
			tools := toolbarButtons()
			if idx < len(tools) {
				editor.StartPlacement(core.PartTypeID(tools[idx].TypeID))
			}
			return true
		}
	}
	if idx := render.ToolbarButtonAtScreenPoint(render.ToolbarRight, rightToolbarButtons(), mouseX, mouseY); idx >= 0 {
		cmds := rightToolbarButtons()
		if idx < len(cmds) {
			switch cmds[idx].TypeID {
			case "_save":
				_ = SaveProject(DefaultProjectPath)
			case "_load":
				_ = LoadProject(DefaultProjectPath)
			case "_settings":
				a.settingsOpen = !a.settingsOpen
			}
		}
		return true
	}
	return false
}

func (a *App) syncRenderThemeFromSettings() {
	render.DarkMode = appsettings.Current.DarkMode
}

func (a *App) settingsPanelRows() []string {
	spec := appsettings.BuildSpec()
	rows := make([]string, 0, len(spec.Items)+2)
	for _, item := range spec.Items {
		switch item.Kind {
		case appsettings.ItemBool:
			v, _ := item.Value.(bool)
			state := "OFF"
			if v {
				state = "ON"
			}
			rows = append(rows, fmt.Sprintf("%s: %s", item.Label, state))
		default:
			rows = append(rows, fmt.Sprintf("%s", item.Label))
		}
	}
	rows = append(rows, "")
	rows = append(rows, "F4: toggle dark mode")
	rows = append(rows, "F3: close settings")
	return rows
}

// statusText reports the current top-level operating mode.
func (a *App) statusText() string {
	base := "Edit mode active"
	if world.RunMode {
		base = "Run mode active"
	}
	if os.Getenv("COILFORGE_UNDO_MEM") == "" {
		return base
	}
	u, r := editor.UndoRedoStacksApproxBytes()
	return fmt.Sprintf("%s  undo~%s redo~%s", base, formatApproxBytes(u), formatApproxBytes(r))
}

func formatApproxBytes(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%dB", n)
	}
	if n < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(n)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(n)/(1024*1024))
}

// updateToolbarHover computes hovered toolbar button indices from mouse/touch pointer.
func (a *App) updateToolbarHover(mouseX, mouseY int) {
	pointerX, pointerY, ok := toolbarPointerPosition(mouseX, mouseY)
	if !ok {
		a.hoverLeftTool = -1
		a.hoverRightTool = -1
		return
	}
	if world.RunMode {
		a.hoverLeftTool = -1
	} else {
		a.hoverLeftTool = render.ToolbarButtonAtScreenPoint(render.ToolbarLeft, toolbarButtons(), pointerX, pointerY)
	}
	a.hoverRightTool = render.ToolbarButtonAtScreenPoint(render.ToolbarRight, rightToolbarButtons(), pointerX, pointerY)
}

// toolbarPointerPosition picks the active pointer for hover (touch preferred, then mouse).
func toolbarPointerPosition(mouseX, mouseY int) (int, int, bool) {
	touchIDs := ebiten.AppendTouchIDs(nil)
	if len(touchIDs) > 0 {
		tx, ty := ebiten.TouchPosition(touchIDs[0])
		return tx, ty, true
	}
	return mouseX, mouseY, true
}
