package render

// File overview:
// theme centralizes render colors, spacing, and style constants.
// Subsystem: render theme.
// It is referenced by scene and chrome drawing helpers for consistent visuals.
// Flow position: styling dependency beneath all rendering routines.

import (
	"coilforge/internal/core"
	"image/color"
)

// DarkMode stores package-level state.
var DarkMode = true // Toggles dark-versus-light theme color selection.

// WireColor returns the display color for a wire by resolved net state.
func WireColor(state int) color.RGBA {
	switch state {
	case core.NetHigh:
		return color.RGBA{R: 255, G: 180, B: 64, A: 255}
	case core.NetLow:
		return color.RGBA{R: 96, G: 180, B: 255, A: 255}
	case core.NetShort:
		return color.RGBA{R: 255, G: 64, B: 64, A: 255}
	default:
		return color.RGBA{R: 160, G: 160, B: 160, A: 255}
	}
}

// SchematicBackgroundColor fills the schematic canvas behind grid and parts.
func SchematicBackgroundColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 22, G: 24, B: 29, A: 255}
	}
	return color.RGBA{R: 244, G: 246, B: 250, A: 255}
}

// GridColor returns the schematic grid color for the active theme (major lines).
func GridColor() color.RGBA {
	return GridMajorColor()
}

// GridMinorColor returns fine grid lines (wire routing pitch); kept faint vs [SchematicBackgroundColor].
// Light mode uses opaque RGB slightly below the canvas so lines never read brighter than the fill.
func GridMinorColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 40, G: 44, B: 52, A: 95}
	}
	return color.RGBA{R: 214, G: 220, B: 232, A: 255}
}

// GridMajorColor returns coarse grid lines (part placement pitch); stronger contrast than minor.
func GridMajorColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 88, G: 96, B: 112, A: 240}
	}
	return color.RGBA{R: 164, G: 172, B: 188, A: 255}
}

// SelectionColor returns the outline color for selected items.
func SelectionColor() color.RGBA {
	return color.RGBA{R: 255, G: 208, B: 64, A: 255}
}

// GhostTint returns the tint color used for translucent preview visuals.
func GhostTint() color.RGBA {
	return color.RGBA{R: 255, G: 255, B: 255, A: 144}
}

// ToolbarPanelColor returns the background fill for chrome toolbar panels.
func ToolbarPanelColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 36, G: 40, B: 48, A: 240}
	}
	return color.RGBA{R: 236, G: 238, B: 242, A: 240}
}

// ToolbarButtonOutlineColor returns the border color by interaction state.
func ToolbarButtonOutlineColor(active, hovered, disabled bool) color.RGBA {
	if disabled {
		if DarkMode {
			return color.RGBA{R: 76, G: 82, B: 94, A: 210}
		}
		return color.RGBA{R: 166, G: 172, B: 184, A: 210}
	}
	if active {
		if DarkMode {
			return color.RGBA{R: 164, G: 136, B: 90, A: 255}
		}
		return color.RGBA{R: 178, G: 146, B: 96, A: 255}
	}
	if hovered {
		if DarkMode {
			return color.RGBA{R: 118, G: 130, B: 152, A: 255}
		}
		return color.RGBA{R: 132, G: 142, B: 164, A: 255}
	}
	if DarkMode {
		return color.RGBA{R: 92, G: 102, B: 122, A: 255}
	}
	return color.RGBA{R: 144, G: 152, B: 170, A: 255}
}

// ToolbarButtonFillColor returns the button background by interaction state.
func ToolbarButtonFillColor(active, hovered, disabled bool) color.RGBA {
	if disabled {
		return ToolbarButtonDisabledFillColor()
	}
	if active {
		if DarkMode {
			return color.RGBA{R: 88, G: 74, B: 48, A: 235}
		}
		return color.RGBA{R: 255, G: 229, B: 176, A: 240}
	}
	if hovered {
		if DarkMode {
			return color.RGBA{R: 70, G: 78, B: 94, A: 225}
		}
		return color.RGBA{R: 232, G: 238, B: 248, A: 230}
	}
	if DarkMode {
		return color.RGBA{R: 52, G: 58, B: 70, A: 215}
	}
	return color.RGBA{R: 214, G: 221, B: 232, A: 220}
}

// ToolbarButtonDisabledFillColor is the face fill for disabled buttons.
func ToolbarButtonDisabledFillColor() color.RGBA {
	if DarkMode {
		return color.RGBA{R: 38, G: 42, B: 50, A: 188}
	}
	return color.RGBA{R: 230, G: 233, B: 239, A: 188}
}

// ToolbarButtonBevelTopLeftColor returns the top/left bevel shade.
func ToolbarButtonBevelTopLeftColor(active, disabled bool) color.RGBA {
	if disabled {
		if DarkMode {
			return color.RGBA{R: 56, G: 62, B: 74, A: 165}
		}
		return color.RGBA{R: 236, G: 239, B: 244, A: 165}
	}
	if active {
		if DarkMode {
			return color.RGBA{R: 56, G: 62, B: 74, A: 220}
		}
		return color.RGBA{R: 206, G: 184, B: 140, A: 220}
	}
	if DarkMode {
		return color.RGBA{R: 108, G: 118, B: 138, A: 220}
	}
	return color.RGBA{R: 246, G: 250, B: 255, A: 220}
}

// ToolbarButtonBevelBottomRightColor returns the bottom/right bevel shade.
func ToolbarButtonBevelBottomRightColor(active, disabled bool) color.RGBA {
	if disabled {
		if DarkMode {
			return color.RGBA{R: 24, G: 28, B: 36, A: 165}
		}
		return color.RGBA{R: 194, G: 200, B: 212, A: 165}
	}
	if active {
		if DarkMode {
			return color.RGBA{R: 116, G: 96, B: 64, A: 220}
		}
		return color.RGBA{R: 160, G: 136, B: 96, A: 220}
	}
	if DarkMode {
		return color.RGBA{R: 34, G: 38, B: 46, A: 220}
	}
	return color.RGBA{R: 162, G: 174, B: 196, A: 220}
}

// ToolbarIconTintColor returns icon tint by interaction state.
func ToolbarIconTintColor(active, hovered, disabled bool) color.RGBA {
	if disabled {
		if DarkMode {
			return color.RGBA{R: 118, G: 126, B: 142, A: 56}
		}
		return color.RGBA{R: 132, G: 138, B: 150, A: 60}
	}
	_ = active
	_ = hovered
	// Theme tint assumes light (white) source icons so dark/light themes can recolor.
	if DarkMode {
		return color.RGBA{R: 230, G: 236, B: 246, A: 255}
	}
	return color.RGBA{R: 56, G: 62, B: 74, A: 255}
}

// ToolbarLabelColor returns fallback label color when an icon is unavailable.
func ToolbarLabelColor(active, hovered, disabled bool) color.RGBA {
	if disabled {
		if DarkMode {
			return color.RGBA{R: 150, G: 156, B: 170, A: 130}
		}
		return color.RGBA{R: 112, G: 118, B: 130, A: 130}
	}
	if active {
		if DarkMode {
			return color.RGBA{R: 248, G: 230, B: 188, A: 255}
		}
		return color.RGBA{R: 98, G: 72, B: 34, A: 255}
	}
	if hovered {
		if DarkMode {
			return color.RGBA{R: 236, G: 242, B: 252, A: 255}
		}
		return color.RGBA{R: 52, G: 58, B: 72, A: 255}
	}
	if DarkMode {
		return color.RGBA{R: 214, G: 220, B: 232, A: 245}
	}
	return color.RGBA{R: 70, G: 76, B: 88, A: 245}
}
