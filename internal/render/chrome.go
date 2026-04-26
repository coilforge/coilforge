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
	"strconv"
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

// PanelFrameStyle defines visual-only frame styling for chrome panels.
type PanelFrameStyle struct {
	FillColor      color.Color
	OutlineColor   color.Color
	OutlineWidth   float32
	BevelLight     color.Color
	BevelDark      color.Color
	BevelInset     float32
	BevelLineWidth float32
}

// DocBrowserRow is one visible row in the save/load browser list.
type DocBrowserRow struct {
	Text     string // row label text.
	Selected bool   // selected row highlight.
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

	simplePanelWidthPx       = 520
	simplePanelHeightPx      = 300
	simplePanelMinMarginPx   = 20
	simplePanelCloseSizePx   = 28
	simplePanelCloseInsetPx  = 12
	simplePanelCloseTextOffX = 8
	simplePanelCloseTextOffY = 8

	propPanelWidthPx      = 300
	propPanelMinHeightPx  = 84
	propPanelRightPadPx   = toolbarStripWidthPx + 16
	propPanelTopPadPx     = 16
	propPanelRowHeightPx  = 34
	propPanelButtonSizePx = 22
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
	DrawPanelFrame(dst, x, y, bw, bh, toolbarPanelFrameStyle())

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

func toolbarPanelFrameStyle() PanelFrameStyle {
	return PanelFrameStyle{
		FillColor:      ToolbarPanelColor(),
		OutlineColor:   ToolbarPanelOutlineColor(),
		OutlineWidth:   2.0,
		BevelLight:     ToolbarPanelBevelTopLeftColor(),
		BevelDark:      ToolbarPanelBevelBottomRightColor(),
		BevelInset:     2.0,
		BevelLineWidth: 2.0,
	}
}

// DrawPanelFrame renders a generic chrome panel frame (fill, outline, optional bevel).
func DrawPanelFrame(dst *ebiten.Image, x, y, w, h float32, style PanelFrameStyle) {
	vector.FillRect(dst, x, y, w, h, style.FillColor, false)
	if style.OutlineWidth > 0 && style.OutlineColor != nil {
		vector.StrokeRect(dst, x, y, w, h, style.OutlineWidth, style.OutlineColor, false)
	}
	if style.BevelLineWidth <= 0 || style.BevelInset < 0 || style.BevelLight == nil || style.BevelDark == nil {
		return
	}
	inset := style.BevelInset
	sw := style.BevelLineWidth
	x1 := x + w - inset
	y1 := y + h - inset
	// Top + left edge.
	vector.StrokeLine(dst, x+inset, y+inset, x1, y+inset, sw, style.BevelLight, false)
	vector.StrokeLine(dst, x+inset, y+inset, x+inset, y1, sw, style.BevelLight, false)
	// Bottom + right edge.
	vector.StrokeLine(dst, x+inset, y1, x1, y1, sw, style.BevelDark, false)
	vector.StrokeLine(dst, x1, y+inset, x1, y1, sw, style.BevelDark, false)
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
	x, y, w, h, ok := propPanelLayout(spec, world.ScreenW, world.ScreenH)
	if !ok {
		return
	}
	DrawPanelFrame(dst, x, y, w, h, toolbarPanelFrameStyle())
	title := "Properties"
	drawAtlasText(dst, strings.ToUpper(title), snapToLogicalPixel(float64(x+12)), snapToLogicalPixel(float64(y+12)), StatusBarTextColor())

	rowY := y + 38
	for i, item := range spec.Items {
		if rowY+propPanelRowHeightPx > y+h-6 {
			break
		}
		drawAtlasText(dst, normalizeUIString(item.Label), snapToLogicalPixel(float64(x+12)), snapToLogicalPixel(float64(rowY+8)), StatusBarTextColor())
		switch item.Kind {
		case part.PropInt:
			drawPropIntControls(dst, x, rowY, w, item)
		default:
			drawAtlasText(dst, normalizeUIString(propValueText(item.Value)), snapToLogicalPixel(float64(x+w-118)), snapToLogicalPixel(float64(rowY+8)), StatusBarTextColor())
		}
		_ = i
		rowY += propPanelRowHeightPx
	}
}

func propValueText(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func drawPropIntControls(dst *ebiten.Image, x, rowY, w float32, item part.PropItem) {
	minusX, minusY, plusX, plusY, valueX, valueY := propIntControlRects(x, rowY, w)
	btnFill := ToolbarButtonFillColor(false, false, false)
	btnOutline := ToolbarPanelOutlineColor()
	vector.FillRect(dst, minusX, minusY, float32(propPanelButtonSizePx), float32(propPanelButtonSizePx), btnFill, false)
	vector.StrokeRect(dst, minusX, minusY, float32(propPanelButtonSizePx), float32(propPanelButtonSizePx), 2.0, btnOutline, false)
	vector.FillRect(dst, plusX, plusY, float32(propPanelButtonSizePx), float32(propPanelButtonSizePx), btnFill, false)
	vector.StrokeRect(dst, plusX, plusY, float32(propPanelButtonSizePx), float32(propPanelButtonSizePx), 2.0, btnOutline, false)
	drawAtlasText(dst, "-", snapToLogicalPixel(float64(minusX+7)), snapToLogicalPixel(float64(minusY+5)), StatusBarTextColor())
	drawAtlasText(dst, "+", snapToLogicalPixel(float64(plusX+6)), snapToLogicalPixel(float64(plusY+5)), StatusBarTextColor())
	drawAtlasText(dst, strconv.Itoa(propIntValue(item)), snapToLogicalPixel(float64(valueX)), snapToLogicalPixel(float64(valueY)), StatusBarTextColor())
}

func propIntValue(item part.PropItem) int {
	if v, ok := item.Value.(int); ok {
		return v
	}
	return 0
}

func propIntControlRects(panelX, rowY, panelW float32) (minusX, minusY, plusX, plusY, valueX, valueY float32) {
	btn := float32(propPanelButtonSizePx)
	plusX = panelX + panelW - 12 - btn
	minusX = plusX - 66
	minusY = rowY + 6
	plusY = rowY + 6
	valueX = minusX + btn + 16
	valueY = rowY + 8
	return minusX, minusY, plusX, plusY, valueX, valueY
}

// PropPanelIntButtonAtScreenPoint returns the prop row index and delta (-1/+1) for clicked int controls.
func PropPanelIntButtonAtScreenPoint(spec part.PropSpec, sx, sy int) (itemIdx int, delta int, ok bool) {
	x, y, w, h, okLayout := propPanelLayout(spec, world.ScreenW, world.ScreenH)
	if !okLayout {
		return -1, 0, false
	}
	if float32(sx) < x || float32(sx) > x+w || float32(sy) < y || float32(sy) > y+h {
		return -1, 0, false
	}
	rowY := y + 38
	for i, item := range spec.Items {
		if rowY+propPanelRowHeightPx > y+h-6 {
			break
		}
		if item.Kind != part.PropInt {
			rowY += propPanelRowHeightPx
			continue
		}
		minusX, minusY, plusX, plusY, _, _ := propIntControlRects(x, rowY, w)
		btn := float32(propPanelButtonSizePx)
		px := float32(sx)
		py := float32(sy)
		if px >= minusX && px <= minusX+btn && py >= minusY && py <= minusY+btn {
			return i, -1, true
		}
		if px >= plusX && px <= plusX+btn && py >= plusY && py <= plusY+btn {
			return i, 1, true
		}
		rowY += propPanelRowHeightPx
	}
	return -1, 0, false
}

func propPanelLayout(spec part.PropSpec, screenW, screenH int) (x, y, w, h float32, ok bool) {
	if screenW <= 0 || screenH <= 0 || len(spec.Items) == 0 {
		return 0, 0, 0, 0, false
	}
	w = float32(propPanelWidthPx)
	h = float32(38 + len(spec.Items)*propPanelRowHeightPx + 10)
	if h < float32(propPanelMinHeightPx) {
		h = float32(propPanelMinHeightPx)
	}
	x = float32(screenW) - w - float32(propPanelRightPadPx)
	y = float32(propPanelTopPadPx)
	if x < 0 {
		x = 0
	}
	if y+h > float32(screenH)-float32(toolbarBottomClearancePx) {
		h = float32(screenH) - float32(toolbarBottomClearancePx) - y
	}
	if h < 40 {
		return 0, 0, 0, 0, false
	}
	return x, y, w, h, true
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

// DrawSimplePanel renders a centered panel with title and simple text rows.
func DrawSimplePanel(dst *ebiten.Image, title string, rows []string) {
	x, y, pw, ph, ok := simplePanelLayout(world.ScreenW, world.ScreenH)
	if !ok {
		return
	}
	DrawPanelFrame(dst, x, y, pw, ph, toolbarPanelFrameStyle())
	drawSimplePanelCloseButton(dst, x, y, pw)

	titleText := strings.TrimSpace(normalizeUIString(title))
	if titleText != "" {
		drawAtlasText(dst, strings.ToUpper(titleText), snapToLogicalPixel(float64(x+16)), snapToLogicalPixel(float64(y+14)), StatusBarTextColor())
	}

	rowY := y + 56
	for i := range rows {
		if rowY > y+ph-24 {
			break
		}
		line := strings.TrimSpace(normalizeUIString(rows[i]))
		if line != "" {
			drawAtlasText(dst, line, snapToLogicalPixel(float64(x+16)), snapToLogicalPixel(float64(rowY)), StatusBarTextColor())
		}
		rowY += 22
	}
}

// DrawTextInputBox renders a simple single-line text input field.
func DrawTextInputBox(dst *ebiten.Image, x, y, w, h float32, text string, active bool) {
	fill := ToolbarButtonFillColor(false, false, false)
	outline := ToolbarButtonOutlineColor(false, active, false)
	vector.FillRect(dst, x, y, w, h, fill, false)
	vector.StrokeRect(dst, x, y, w, h, 2.0, outline, false)
	drawAtlasText(dst, strings.TrimSpace(normalizeUIString(text)), snapToLogicalPixel(float64(x+10)), snapToLogicalPixel(float64(y+8)), ToolbarLabelColor(false, active, false))
}

// DrawDocBrowserPanel renders a centered save/load browser panel.
func DrawDocBrowserPanel(dst *ebiten.Image, title, fileName string, rows []DocBrowserRow, footer string) {
	x, y, pw, ph, ok := simplePanelLayout(world.ScreenW, world.ScreenH)
	if !ok {
		return
	}
	DrawPanelFrame(dst, x, y, pw, ph, toolbarPanelFrameStyle())
	drawSimplePanelCloseButton(dst, x, y, pw)

	drawAtlasText(dst, strings.ToUpper(strings.TrimSpace(normalizeUIString(title))), snapToLogicalPixel(float64(x+16)), snapToLogicalPixel(float64(y+14)), StatusBarTextColor())
	drawAtlasText(dst, "Filename", snapToLogicalPixel(float64(x+16)), snapToLogicalPixel(float64(y+42)), StatusBarTextColor())
	DrawTextInputBox(dst, x+16, y+58, pw-32, 34, fileName, true)

	listX := x + 16
	listY := y + 104
	listW := pw - 32
	listH := ph - 144
	vector.FillRect(dst, listX, listY, listW, listH, ToolbarButtonDisabledFillColor(), false)
	vector.StrokeRect(dst, listX, listY, listW, listH, 2.0, ToolbarPanelOutlineColor(), false)

	rowY := listY + 8
	for i := range rows {
		if rowY+22 > listY+listH-4 {
			break
		}
		if rows[i].Selected {
			vector.FillRect(dst, listX+4, rowY-2, listW-8, 22, ToolbarButtonFillColor(false, true, false), false)
		}
		drawAtlasText(dst, normalizeUIString(rows[i].Text), snapToLogicalPixel(float64(listX+8)), snapToLogicalPixel(float64(rowY)), ToolbarLabelColor(false, rows[i].Selected, false))
		rowY += 22
	}

	if footer != "" {
		drawAtlasText(dst, normalizeUIString(footer), snapToLogicalPixel(float64(x+16)), snapToLogicalPixel(float64(y+ph-28)), StatusBarTextColor())
	}
}

// SimplePanelContainsScreenPoint reports whether a screen-space pointer is inside the centered simple panel.
func SimplePanelContainsScreenPoint(sx, sy int) bool {
	x, y, w, h, ok := simplePanelLayout(world.ScreenW, world.ScreenH)
	if !ok {
		return false
	}
	px := float32(sx)
	py := float32(sy)
	return px >= x && px <= x+w && py >= y && py <= y+h
}

// SimplePanelCloseButtonAtScreenPoint reports whether a pointer hits the panel [X] button.
func SimplePanelCloseButtonAtScreenPoint(sx, sy int) bool {
	x, y, w, _, ok := simplePanelLayout(world.ScreenW, world.ScreenH)
	if !ok {
		return false
	}
	bx, by, bw, bh := simplePanelCloseButtonRect(x, y, w)
	px := float32(sx)
	py := float32(sy)
	return px >= bx && px <= bx+bw && py >= by && py <= by+bh
}

func simplePanelLayout(screenW, screenH int) (x, y, w, h float32, ok bool) {
	if screenW <= 0 || screenH <= 0 {
		return 0, 0, 0, 0, false
	}
	w = float32(simplePanelWidthPx)
	h = float32(simplePanelHeightPx)
	margin := float32(simplePanelMinMarginPx)
	if float32(screenW) < w+2*margin {
		w = float32(screenW) - 2*margin
	}
	if float32(screenH) < h+2*margin {
		h = float32(screenH) - 2*margin
	}
	if w < 120 || h < 80 {
		return 0, 0, 0, 0, false
	}
	x = (float32(screenW) - w) * 0.5
	y = (float32(screenH) - h) * 0.5
	return x, y, w, h, true
}

func simplePanelCloseButtonRect(panelX, panelY, panelW float32) (x, y, w, h float32) {
	size := float32(simplePanelCloseSizePx)
	inset := float32(simplePanelCloseInsetPx)
	return panelX + panelW - inset - size, panelY + inset, size, size
}

func drawSimplePanelCloseButton(dst *ebiten.Image, panelX, panelY, panelW float32) {
	x, y, w, h := simplePanelCloseButtonRect(panelX, panelY, panelW)
	vector.FillRect(dst, x, y, w, h, ToolbarButtonFillColor(false, false, false), false)
	vector.StrokeRect(dst, x, y, w, h, 2.0, ToolbarButtonOutlineColor(false, false, false), false)
	drawAtlasText(
		dst,
		"X",
		snapToLogicalPixel(float64(x+float32(simplePanelCloseTextOffX))),
		snapToLogicalPixel(float64(y+float32(simplePanelCloseTextOffY))),
		ToolbarLabelColor(false, false, false),
	)
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
