package app

import (
	"fmt"

	"coilforge/internal/appsettings"
	"coilforge/internal/render"
)

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

// closeTopmostOverlay closes the topmost closable app-level UI and reports whether it handled Esc.
func (a *App) closeTopmostOverlay() bool {
	if a.docDialog.mode != docDialogClosed {
		a.closeDocDialog("")
		return true
	}
	if a.settingsOpen {
		a.settingsOpen = false
		return true
	}
	return false
}
