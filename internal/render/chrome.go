package render

// File overview:
// chrome draws screen-space UI elements such as toolbar, status, and property panels.
// Subsystem: render chrome.
// It complements scene rendering and is invoked by app through render entrypoints.
// Flow position: final UI overlay layer on top of world-space part drawing.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ToolButton struct {
	TypeID   string // Part type identifier associated with this button.
	Hotkey   rune   // Keyboard shortcut shown for quick tool selection.
	Label    string // Human-readable button label for chrome rendering.
	Disabled bool   // Disabled buttons are rendered but not interactive.
}

// Toolbar dock side (plain ints). Submenus can use the same values to pick a direction.
const (
	ToolbarLeft  = 0
	ToolbarRight = 1
)

// Chrome layout for vertical toolbar strips (screen pixels).
// Kept ~2× the original design to match doubled schematic scale (SVGUserUnitToWorld) + larger UI font atlas.
const (
	chromeEdgeMargin = 16 // Used by status bar, sim HUD, schematic chrome; not for docking the toolbar strips.

	toolbarStripWidthPx      = 112
	toolbarPanelInnerPadPx   = 8
	toolbarButtonHitPx       = 96 // Touch-style hit target; fits inside strip with inner pad.
	toolbarButtonGapPx       = 12
	toolbarIconSlotPx        = 56 // Icon drawn scaled inside this square.
	toolbarHitStrokeWidth    = 2.0
	toolbarActiveStrokeWidth = 4.0

	statusBarBottomMarginPx = 20
	// Toolbar strips are flush to top and side edges; only the bottom is inset so status/HUD text stays readable.
	toolbarBottomClearancePx = statusBarBottomMarginPx + 30

	// simRealtimeRightPad is space from the right window edge reserved for the right toolbar (flush to edge).
	simRealtimeRightPad = toolbarStripWidthPx
)

// toolbarStripLayout returns the toolbar panel rectangle for [DrawToolbar] and hit-testing (same math).
func toolbarStripLayout(side int, w, h int) (x, y, bw, bh float32) {
	bw = float32(toolbarStripWidthPx)
	bh = float32(h) - float32(toolbarBottomClearancePx)
	if bh < 1 {
		bh = 1
	}
	y = 0
	switch side {
	case ToolbarLeft:
		x = 0
	case ToolbarRight:
		x = float32(w) - bw
	default:
		return 0, 0, 0, 0
	}
	return x, y, bw, bh
}

// ToolbarButtonAtScreenPoint returns the button index under a screen-space pointer,
// or -1 when the pointer is outside the visible button stack.
func ToolbarButtonAtScreenPoint(side int, tools []ToolButton, sx, sy int) int {
	w, h := world.ScreenW, world.ScreenH
	if w <= 0 || h <= 0 {
		return -1
	}
	x, y, bw, bh := toolbarStripLayout(side, w, h)
	if bw < 1 || bh < 1 {
		return -1
	}
	if len(tools) == 0 {
		return -1
	}

	inner := float32(toolbarPanelInnerPadPx)
	hit := float32(toolbarButtonHitPx)
	gap := float32(toolbarButtonGapPx)
	contentLeft := x + inner + (bw-2*inner-hit)*0.5
	contentTop := y + inner
	maxY := y + bh - inner

	for i := range tools {
		y := contentTop + float32(i)*(hit+gap)
		if y+hit > maxY+0.01 {
			break
		}
		if tools[i].Disabled {
			continue
		}
		if float32(sx) >= contentLeft && float32(sx) <= contentLeft+hit &&
			float32(sy) >= y && float32(sy) <= y+hit {
			return i
		}
	}
	return -1
}

// DrawToolbar draws the toolbar panel and stacked tool buttons with optional
// active/hover styling and centered icon rendering.
func DrawToolbar(dst *ebiten.Image, side int, tools []ToolButton, activeTool int, hoverTool int) {
	w, h := world.ScreenW, world.ScreenH
	if w <= 0 || h <= 0 {
		return
	}
	x, y, bw, bh := toolbarStripLayout(side, w, h)
	if bw < 1 || bh < 1 {
		return
	}
	vector.FillRect(dst, x, y, bw, bh, ToolbarPanelColor(), false)

	if len(tools) == 0 {
		return
	}

	inner := float32(toolbarPanelInnerPadPx)
	hit := float32(toolbarButtonHitPx)
	gap := float32(toolbarButtonGapPx)
	iconSz := float32(toolbarIconSlotPx)
	// Center the square hit target in the strip.
	contentLeft := x + inner + (bw-2*inner-hit)*0.5
	contentTop := y + inner
	maxY := y + bh - inner

	for i := range tools {
		y := contentTop + float32(i)*(hit+gap)
		if y+hit > maxY+0.01 {
			break
		}
		drawToolbarButton(dst, tools[i], i, activeTool, hoverTool, contentLeft, y, hit, iconSz)
	}
}

func drawToolbarButton(dst *ebiten.Image, btn ToolButton, index, activeTool, hoverTool int, contentLeft, y, hit, iconSz float32) {
	active := index == activeTool
	disabled := btn.Disabled
	if disabled {
		active = false
	}
	hovered := index == hoverTool && !disabled
	sw := toolbarButtonStrokeWidth(active, hovered)
	vector.FillRect(dst, contentLeft, y, hit, hit, ToolbarButtonFillColor(active, hovered, disabled), false)
	vector.StrokeRect(dst, contentLeft, y, hit, hit, sw, ToolbarButtonOutlineColor(active, hovered, disabled), false)
	drawButtonBevel(dst, contentLeft, y, hit, active, disabled)
	if drawToolbarButtonIcon(dst, btn, contentLeft, y, hit, iconSz, active, hovered, disabled) {
		return
	}
	drawToolbarLabel(dst, btn.Label, contentLeft, y, hit, ToolbarLabelColor(active, hovered, disabled), active, hovered)
}

func toolbarButtonStrokeWidth(active, hovered bool) float32 {
	if active {
		return toolbarActiveStrokeWidth
	}
	if hovered {
		return 3.0
	}
	return float32(toolbarHitStrokeWidth)
}

func drawToolbarButtonIcon(dst *ebiten.Image, btn ToolButton, x, y, hit, iconSz float32, active, hovered, disabled bool) bool {
	info, ok := part.Registry[core.PartTypeID(btn.TypeID)]
	if !ok || info.Icon == nil {
		return false
	}
	img := info.Icon()
	if img == nil {
		return false
	}
	b := img.Bounds()
	iw, ih := b.Dx(), b.Dy()
	if iw <= 0 || ih <= 0 {
		return false
	}
	off := (hit - iconSz) * 0.5
	ix := x + off
	iy := y + off
	scale := float64(iconSz) / float64(max(iw, ih))
	drawW := float64(iw) * scale
	drawH := float64(ih) * scale
	nudgeX := 0.0
	nudgeY := 0.0
	if active || hovered {
		// Positional nudge reads as press/hover better than icon scaling.
		nudgeX = 2
		nudgeY = 2
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(ix)+(float64(iconSz)-drawW)*0.5+nudgeX,
		float64(iy)+(float64(iconSz)-drawH)*0.5+nudgeY,
	)
	tint := ToolbarIconTintColor(active, hovered, disabled)
	op.ColorScale.Scale(
		float32(tint.R)/255.0,
		float32(tint.G)/255.0,
		float32(tint.B)/255.0,
		1.0,
	)
	op.ColorScale.ScaleAlpha(float32(tint.A) / 255.0)
	dst.DrawImage(img, op)
	return true
}

func drawToolbarLabel(dst *ebiten.Image, label string, x, y, size float32, clr color.Color, active, hovered bool) {
	trimmed := normalizeUIString(label)
	if trimmed == "" {
		return
	}
	textLabel := strings.ToUpper(trimmed)
	if len(textLabel) > 6 {
		textLabel = textLabel[:6]
	}
	atlas := uiLabelAtlas()
	aw, ah := atlasMeasure(textLabel, atlas)
	nudgeX := 0.0
	nudgeY := 0.0
	if active || hovered {
		// Match icon nudge in drawToolbarButtonIcon — reads as hover/press affordance.
		nudgeX = 2
		nudgeY = 2
	}
	targetX := snapToLogicalPixel(float64(x+(size-float32(aw))*0.5) + nudgeX)
	targetY := snapToLogicalPixel(float64(y+(size-float32(ah))*0.5) + nudgeY)
	drawAtlasText(dst, textLabel, targetX, targetY, clr)
}

func snapToLogicalPixel(v float64) float64 {
	return math.Round(v)
}

func drawButtonBevel(dst *ebiten.Image, x, y, size float32, active bool, disabled bool) {
	light := ToolbarButtonBevelTopLeftColor(active, disabled)
	dark := ToolbarButtonBevelBottomRightColor(active, disabled)
	if !DarkMode && !active {
		// In light mode, non-active buttons should read raised rather than inset.
		light, dark = dark, light
	}
	inset := float32(2)
	sw := float32(2)
	// Top + left edge.
	vector.StrokeLine(dst, x+inset, y+inset, x+size-inset, y+inset, sw, light, false)
	vector.StrokeLine(dst, x+inset, y+inset, x+inset, y+size-inset, sw, light, false)
	// Bottom + right edge.
	vector.StrokeLine(dst, x+inset, y+size-inset, x+size-inset, y+size-inset, sw, dark, false)
	vector.StrokeLine(dst, x+size-inset, y+inset, x+size-inset, y+size-inset, sw, dark, false)
}

// DrawPropPanel renders the selected-part property panel chrome.
func DrawPropPanel(dst *ebiten.Image, spec part.PropSpec) {
	_, _ = dst, spec
}

// DrawStatusBar renders bottom status text chrome.
func DrawStatusBar(dst *ebiten.Image, text string) {
	text = strings.TrimSpace(normalizeUIString(text))
	if text == "" {
		return
	}
	w, h := world.ScreenW, world.ScreenH
	if w <= 0 || h <= 0 {
		return
	}
	atlas := uiLabelAtlas()
	tw, th := atlasMeasure(text, atlas)
	margin := float64(chromeEdgeMargin)
	maxW := float64(w) - 2*margin
	s := text
	runes := []rune(s)
	for tw > maxW && len(runes) > 4 {
		runes = runes[:len(runes)-1]
		s = string(runes)
		tw, th = atlasMeasure(s, atlas)
	}
	if tw > maxW {
		s = "..."
		tw, th = atlasMeasure(s, atlas)
	}
	targetX := snapToLogicalPixel(margin)
	targetY := snapToLogicalPixel(float64(h) - float64(statusBarBottomMarginPx) - th)
	drawAtlasText(dst, s, targetX, targetY, StatusBarTextColor())
}

// DrawSimRealtimeHUD draws bottom-right atlas text (simulated vs wall-clock rate).
func DrawSimRealtimeHUD(dst *ebiten.Image, text string) {
	text = strings.TrimSpace(normalizeUIString(text))
	if text == "" {
		return
	}
	w, h := world.ScreenW, world.ScreenH
	if w <= 0 || h <= 0 {
		return
	}
	atlas := uiLabelAtlas()
	tw, th := atlasMeasure(text, atlas)
	leftPad := float64(chromeEdgeMargin)
	rightPad := float64(simRealtimeRightPad)
	maxW := float64(w) - leftPad - rightPad
	if maxW < 8 {
		return
	}
	s := text
	runes := []rune(s)
	for tw > maxW && len(runes) > 4 {
		runes = runes[:len(runes)-1]
		s = string(runes)
		tw, th = atlasMeasure(s, atlas)
	}
	if tw > maxW {
		s = "..."
		tw, th = atlasMeasure(s, atlas)
	}
	targetX := snapToLogicalPixel(float64(w) - rightPad - tw)
	if targetX < leftPad {
		targetX = leftPad
	}
	targetY := snapToLogicalPixel(float64(h) - float64(statusBarBottomMarginPx) - th)
	drawAtlasText(dst, s, targetX, targetY, StatusBarTextColor())
}

// DrawSelectionOutline renders a highlight around selected geometry.
func DrawSelectionOutline(dst *ebiten.Image, bounds core.Rect) {
	x0, y0 := world.WorldToScreen(core.Pt{X: bounds.Min.X, Y: bounds.Min.Y})
	x1, y1 := world.WorldToScreen(core.Pt{X: bounds.Max.X, Y: bounds.Max.Y})
	minX := min(x0, x1)
	maxX := max(x0, x1)
	minY := min(y0, y1)
	maxY := max(y0, y1)
	sw := float32(maxX - minX)
	sh := float32(maxY - minY)
	if sw < 1 {
		sw = 1
	}
	if sh < 1 {
		sh = 1
	}
	vector.StrokeRect(dst, float32(minX), float32(minY), sw, sh, 3.0, SelectionColor(), false)
}

// DrawBoxSelect renders marquee selection rectangle chrome (world-space rect).
func DrawBoxSelect(dst *ebiten.Image, box core.Rect, crossing bool) {
	x0, y0 := world.WorldToScreen(core.Pt{X: box.Min.X, Y: box.Min.Y})
	x1, y1 := world.WorldToScreen(core.Pt{X: box.Max.X, Y: box.Max.Y})
	minX := min(x0, x1)
	maxX := max(x0, x1)
	minY := min(y0, y1)
	maxY := max(y0, y1)
	sw := float32(maxX - minX)
	sh := float32(maxY - minY)
	if sw < 1 {
		sw = 1
	}
	if sh < 1 {
		sh = 1
	}
	fill := BoxSelectFillColor(crossing)
	vector.FillRect(dst, float32(minX), float32(minY), sw, sh, fill, false)
	vector.StrokeRect(dst, float32(minX), float32(minY), sw, sh, 2.0, SelectionColor(), false)
}
