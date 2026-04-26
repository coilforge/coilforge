package app

import (
	"coilforge/internal/appsettings"
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/partmanifest"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if a.closeTopmostOverlay() {
			return
		}
		if !world.RunMode {
			editor.HandleKey(ebiten.KeyEscape)
		}
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
		a.openDocDialog(docDialogSave)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF7) {
		a.openDocDialog(docDialogLoad)
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
		if a.docDialog.mode != docDialogClosed {
			if render.SimplePanelCloseButtonAtScreenPoint(mouseX, mouseY) {
				a.closeDocDialog("")
			}
			a.toolbarCapture = true
			break
		}
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
				a.openDocDialog(docDialogSave)
			case "_load":
				a.openDocDialog(docDialogLoad)
			case "_settings":
				a.settingsOpen = !a.settingsOpen
			}
		}
		return true
	}
	return false
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
