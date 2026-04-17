package app

// File overview:
// app coordinates frame-level input, mode routing, and top-level draw calls.
// Subsystem: app orchestration.
// It calls editor and sim independently while sharing state through world and render.
// Flow position: primary runtime controller between Ebiten loop and subsystems.

import (
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/partmanifest"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	leftDown       bool // Tracks previous-frame left mouse state for edge detection.
	hoverLeftTool  int  // Hovered placement tool index for left toolbar.
	hoverRightTool int  // Hovered command tool index for right toolbar.
}

// New constructs a fresh application instance.
func New() *App {
	return &App{hoverLeftTool: -1, hoverRightTool: -1}
}

// Run resets world state, configures the window, and starts the Ebiten loop.
func Run() error {
	world.Reset()
	ebiten.SetWindowSize(1440, 900)
	ebiten.SetWindowTitle("CoilForge")
	return ebiten.RunGame(New())
}

// Update runs one frame of app orchestration and input dispatch.
// It refreshes viewport size, processes keyboard/mouse input, and routes input
// into editor or simulation behavior depending on run mode.
// It also advances simulation frames when run mode is active.
func (a *App) Update() error {
	w, h := ebiten.WindowSize()
	world.ScreenW = w
	world.ScreenH = h
	mx, my := ebiten.CursorPosition()
	pt := world.ScreenToWorld(mx, my)
	a.updateToolbarHover(mx, my)

	a.handleToolHotkeys()
	a.handleEditorHotkeys()
	a.handleProjectHotkeys()
	a.handleMouse(pt)

	if _, wheelY := ebiten.Wheel(); wheelY != 0 && !world.RunMode {
		editor.HandleScroll(wheelY)
	}

	if world.RunMode {
		sim.AdvanceFrame()
	}

	return nil
}

// Draw composes scene rendering, editor overlays, and screen-space chrome.
func (a *App) Draw(screen *ebiten.Image) {
	render.DrawScene(screen)

	if !world.RunMode {
		editor.DrawOverlays(screen)
		render.DrawToolbar(screen, render.ToolbarLeft, toolbarButtons(), activeToolIndex(), a.hoverLeftTool)
	}
	// Command strip: visible in edit and run mode (actions not wired yet).
	render.DrawToolbar(screen, render.ToolbarRight, rightToolbarButtons(), -1, a.hoverRightTool)
	if selectedPart := selectedPart(); selectedPart != nil {
		render.DrawPropPanel(screen, selectedPart.PropSpec())
	}
	render.DrawStatusBar(screen, a.statusText())
}

// Layout keeps the game surface the same size as the window.
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
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
	a.handleTransformHotkeys()
	a.handleHistoryHotkeys()
	a.handleLabelHotkeys()
}

// handleTransformHotkeys processes rotate, mirror, delete, and mode keys.
func (a *App) handleTransformHotkeys() {
	for _, key := range []ebiten.Key{
		ebiten.KeyEscape,
		ebiten.KeyR,
		ebiten.KeyM,
		ebiten.KeyDelete,
		ebiten.KeyBackspace,
		ebiten.KeyW,
	} {
		if inpututil.IsKeyJustPressed(key) && !world.RunMode {
			editor.HandleKey(key)
		}
	}
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
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		render.DarkMode = !render.DarkMode
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		ToggleRunMode()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
		_ = SaveProject("coilforge.json")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		_ = LoadProject("coilforge.json")
	}
}

// handleMouse routes click, drag, and release events by active mode.
func (a *App) handleMouse(pt core.Pt) {
	leftNow := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	switch {
	case leftNow && !a.leftDown:
		if world.RunMode {
			sim.HandleClick(pt)
		} else {
			editor.HandleClick(pt, int(ebiten.MouseButtonLeft))
		}
	case leftNow && a.leftDown && !world.RunMode:
		editor.HandleDrag(pt)
	case !leftNow && a.leftDown && !world.RunMode:
		editor.HandleRelease(pt, int(ebiten.MouseButtonLeft))
	}

	a.leftDown = leftNow
}

// statusText reports the current top-level operating mode.
func (a *App) statusText() string {
	if world.RunMode {
		return "Run mode active"
	}
	return "Edit mode active"
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
