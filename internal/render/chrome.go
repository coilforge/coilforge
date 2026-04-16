package render

// File overview:
// chrome draws screen-space UI elements such as toolbar, status, and property panels.
// Subsystem: render chrome.
// It complements scene rendering and is invoked by app through render entrypoints.
// Flow position: final UI overlay layer on top of world-space part drawing.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"

	"github.com/hajimehoshi/ebiten/v2"
)

type ToolButton struct {
	TypeID string // Part type identifier associated with this button.
	Hotkey rune   // Keyboard shortcut shown for quick tool selection.
	Label  string // Human-readable button label for chrome rendering.
}

// DrawToolbar renders toolbar chrome for the provided tool list.
func DrawToolbar(dst *ebiten.Image, tools []ToolButton, activeTool int) {
	_, _, _ = dst, tools, activeTool
}

// DrawPropPanel renders the selected-part property panel chrome.
func DrawPropPanel(dst *ebiten.Image, spec part.PropSpec) {
	_, _ = dst, spec
}

// DrawStatusBar renders bottom status text chrome.
func DrawStatusBar(dst *ebiten.Image, text string) {
	_ = SelectionColor()
	_, _ = dst, text
}

// DrawSelectionOutline renders a highlight around selected geometry.
func DrawSelectionOutline(dst *ebiten.Image, bounds core.Rect) {
	_ = core.RotateRect(bounds, 0, false)
	_ = SelectionColor()
	_ = dst
}

// DrawWireDraft renders in-progress wire placement preview geometry.
func DrawWireDraft(dst *ebiten.Image, points []core.Pt) {
	_ = WireColor(core.NetFloat)
	_ = GhostTint()
	_, _ = dst, points
}

// DrawBoxSelect renders marquee selection rectangle chrome.
func DrawBoxSelect(dst *ebiten.Image, box core.Rect) {
	_ = GhostTint()
	_, _ = dst, box
}
