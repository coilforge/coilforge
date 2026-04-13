package app

import (
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	leftDown bool
}

func New() *App {
	return &App{}
}

func Run() error {
	world.Reset()
	ebiten.SetWindowSize(1440, 900)
	ebiten.SetWindowTitle("CoilForge")
	return ebiten.RunGame(New())
}

func (a *App) Update() error {
	w, h := ebiten.WindowSize()
	world.ScreenW = w
	world.ScreenH = h
	mx, my := ebiten.CursorPosition()
	pt := world.ScreenToWorld(mx, my)

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

func (a *App) Draw(screen *ebiten.Image) {
	render.DrawScene(screen)

	if !world.RunMode {
		editor.DrawOverlays(screen)
	}

	render.DrawToolbar(screen, toolbarButtons(), activeToolIndex())
	if selectedPart := selectedPart(); selectedPart != nil {
		render.DrawPropPanel(screen, selectedPart.PropSpec())
	}
	render.DrawStatusBar(screen, a.statusText())
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (a *App) handleToolHotkeys() {
	toolHotkeys := []struct {
		key    ebiten.Key
		typeID core.PartTypeID
	}{
		{key: ebiten.Key1, typeID: "relay"},
		{key: ebiten.Key2, typeID: "vcc"},
		{key: ebiten.Key3, typeID: "gnd"},
		{key: ebiten.Key4, typeID: "switch"},
		{key: ebiten.Key5, typeID: "indicator"},
		{key: ebiten.Key6, typeID: "diode"},
		{key: ebiten.Key7, typeID: "rch"},
		{key: ebiten.Key8, typeID: "clock"},
	}

	for _, item := range toolHotkeys {
		if inpututil.IsKeyJustPressed(item.key) {
			editor.StartPlacement(item.typeID)
		}
	}
}

func (a *App) handleEditorHotkeys() {
	a.handleTransformHotkeys()
	a.handleHistoryHotkeys()
	a.handleLabelHotkeys()
}

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

func (a *App) handleLabelHotkeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyL) && !world.RunMode && len(editor.Selection) > 0 {
		editor.StartLabelEdit(editor.Selection[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !world.RunMode && editor.LabelEditing {
		editor.CommitLabelEdit()
	}
}

func (a *App) handleProjectHotkeys() {
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

func (a *App) statusText() string {
	if world.RunMode {
		return "Run mode active"
	}
	return "Edit mode active"
}
