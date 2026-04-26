package appsettings

import (
	"os"
	"path/filepath"
)

// File overview:
// settings defines app-level preferences and a UI-oriented spec/action contract.
// Subsystem: app settings domain.
// It is consumed by app orchestration and settings-panel rendering.
// Flow position: local user preference state, separate from project file storage.

type Spec struct {
	Items []Item // settings panel rows.
}

type Item struct {
	Label   string   // label text.
	Kind    int      // kind value.
	Value   any      // current value.
	Choices []string // choices value.
	Min     int      // min value.
	Max     int      // max value.
}

type Action struct {
	Index    int // index value.
	NewValue any // new value.
}

const (
	ItemText = iota // ItemText renders a text row.
	ItemInt         // ItemInt renders an integer row.
	ItemChoice      // ItemChoice renders a choice row.
	ItemBool        // ItemBool renders a boolean row.
	ItemButton      // ItemButton renders a button row.
)

const (
	itemIdxDarkMode = iota
	itemIdxDefaultSaveDir
)

// Values stores in-memory app preference values.
type Values struct {
	DarkMode       bool   `json:"darkMode"`       // Theme toggle preference.
	DefaultSaveDir string `json:"defaultSaveDir"` // Folder used by save/load document browser.
}

// Current stores active app preferences.
var Current = Defaults()

// Defaults returns built-in app preference defaults.
func Defaults() Values {
	return Values{
		DarkMode:       true,
		DefaultSaveDir: resolveDefaultSaveDir(),
	}
}

// BuildSpec returns a UI spec for the settings panel.
func BuildSpec() Spec {
	return Spec{
		Items: []Item{
			{Label: "Dark mode", Kind: ItemBool, Value: Current.DarkMode},
			{Label: "Default save directory", Kind: ItemText, Value: Current.DefaultSaveDir},
		},
	}
}

// Apply applies one setting action and reports whether a setting changed.
func Apply(action Action) bool {
	switch action.Index {
	case itemIdxDarkMode:
		v, ok := action.NewValue.(bool)
		if !ok || Current.DarkMode == v {
			return false
		}
		Current.DarkMode = v
		return true
	case itemIdxDefaultSaveDir:
		v, ok := action.NewValue.(string)
		if !ok || Current.DefaultSaveDir == v {
			return false
		}
		Current.DefaultSaveDir = v
		return true
	default:
		return false
	}
}

// Normalize fills unset settings with built-in defaults.
func Normalize(v Values) Values {
	d := Defaults()
	if v.DefaultSaveDir == "" {
		v.DefaultSaveDir = d.DefaultSaveDir
	}
	return v
}

func resolveDefaultSaveDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "."
	}
	desktop := filepath.Join(home, "Desktop")
	if info, err := os.Stat(desktop); err == nil && info.IsDir() {
		return desktop
	}
	return home
}
