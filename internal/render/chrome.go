package render

import (
	"coilforge/internal/core"
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

type ToolButton struct {
	TypeID string
	Hotkey rune
	Label  string
}

func DrawToolbar(dst *ebiten.Image, tools []ToolButton, activeTool int) {
	_, _, _ = dst, tools, activeTool
}

func DrawPropPanel(dst *ebiten.Image, spec part.PropSpec) {
	_, _ = dst, spec
}

func DrawStatusBar(dst *ebiten.Image, text string) {
	_ = SelectionColor()
	_, _ = dst, text
}

func DrawSelectionOutline(dst *ebiten.Image, bounds core.Rect) {
	_ = core.RotateRect(bounds, 0, false)
	_ = SelectionColor()
	_ = dst
}

func DrawWireDraft(dst *ebiten.Image, points []core.Pt) {
	_ = WireColor(core.NetFloat)
	_ = GhostTint()
	_, _ = dst, points
}

func DrawBoxSelect(dst *ebiten.Image, box core.Rect) {
	_ = GhostTint()
	_, _ = dst, box
}
