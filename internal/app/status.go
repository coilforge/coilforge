package app

import (
	"fmt"
	"os"

	"coilforge/internal/editor"
	"coilforge/internal/world"
)

// statusText reports the current top-level operating mode.
func (a *App) statusText() string {
	base := "Edit mode active"
	if world.RunMode {
		base = "Run mode active"
	}
	if os.Getenv("COILFORGE_UNDO_MEM") == "" {
		return base
	}
	u, r := editor.UndoRedoStacksApproxBytes()
	return fmt.Sprintf("%s  undo~%s redo~%s", base, formatApproxBytes(u), formatApproxBytes(r))
}

func formatApproxBytes(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%dB", n)
	}
	if n < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(n)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(n)/(1024*1024))
}
