package app

import (
	"strings"

	"coilforge/internal/appsettings"
	"coilforge/internal/render"
	"coilforge/internal/storage"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (a *App) syncRenderThemeFromSettings() {
	render.DarkMode = appsettings.Current.DarkMode
}

func (a *App) settingsPanelRows() []render.SettingsRow {
	if strings.TrimSpace(a.settingsPath) == "" {
		a.settingsPath = appsettings.Current.DefaultSaveDir
	}
	return []render.SettingsRow{
		{
			Label:     "Dark mode",
			Kind:      render.SettingsRowCheckbox,
			BoolValue: appsettings.Current.DarkMode,
		},
		{
			Label:     "Default save directory",
			Kind:      render.SettingsRowTextInput,
			TextValue: a.settingsPath,
			Active:    a.settingsPathActive,
		},
	}
}

func (a *App) settingsPanelFooter() []string {
	return []string{
		"Enter: apply path",
		"F4: toggle dark mode    F3: close settings",
	}
}

func (a *App) handleSettingsMousePress(mouseX, mouseY int) {
	rows := a.settingsPanelRows()
	if idx, ok := render.SettingsPanelCheckboxAtScreenPoint(rows, mouseX, mouseY); ok {
		if idx == 0 {
			changed := appsettings.Apply(appsettings.Action{
				Index:    0,
				NewValue: !appsettings.Current.DarkMode,
			})
			if changed {
				_ = appsettings.SaveLocalCurrent()
			}
			a.syncRenderThemeFromSettings()
		}
		return
	}
	if _, ok := render.SettingsPanelTextInputAtScreenPoint(rows, mouseX, mouseY); ok {
		a.settingsPathActive = true
		return
	}
	a.settingsPathActive = false
}

func (a *App) handleSettingsTyping() {
	if !a.settingsOpen {
		return
	}
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		if a.handleSettingsKey(key) {
			return
		}
	}
	if !a.settingsPathActive {
		return
	}
	for _, ch := range ebiten.AppendInputChars(nil) {
		if ch < 32 || ch == 127 {
			continue
		}
		a.settingsPath += string(ch)
	}
}

func (a *App) handleSettingsKey(key ebiten.Key) bool {
	switch key {
	case ebiten.KeyF3:
		a.settingsOpen = false
		a.settingsPathActive = false
		return true
	case ebiten.KeyF4:
		a.toggleDarkModeSetting()
	case ebiten.KeyBackspace:
		a.backspaceSettingsPath()
	case ebiten.KeyEnter:
		a.applySettingsPath()
	}
	return false
}

func (a *App) toggleDarkModeSetting() {
	changed := appsettings.Apply(appsettings.Action{
		Index:    0,
		NewValue: !appsettings.Current.DarkMode,
	})
	if changed {
		_ = appsettings.SaveLocalCurrent()
	}
	a.syncRenderThemeFromSettings()
}

func (a *App) backspaceSettingsPath() {
	if !a.settingsPathActive {
		return
	}
	r := []rune(a.settingsPath)
	if len(r) == 0 {
		return
	}
	a.settingsPath = string(r[:len(r)-1])
}

func (a *App) applySettingsPath() {
	if !a.settingsPathActive {
		return
	}
	path := strings.TrimSpace(a.settingsPath)
	if path == "" {
		return
	}
	changed := appsettings.Apply(appsettings.Action{
		Index:    1,
		NewValue: path,
	})
	if !changed {
		return
	}
	_ = appsettings.SaveLocalCurrent()
	a.docDialog.store = storage.NewLocalFSStore(appsettings.Current.DefaultSaveDir)
	if a.docDialog.mode != docDialogClosed {
		a.refreshDocDialogList()
	}
}

// closeTopmostOverlay closes the topmost closable app-level UI and reports whether it handled Esc.
func (a *App) closeTopmostOverlay() bool {
	if a.docDialog.mode != docDialogClosed {
		a.closeDocDialog("")
		return true
	}
	if a.settingsOpen {
		a.settingsOpen = false
		a.settingsPathActive = false
		return true
	}
	return false
}
