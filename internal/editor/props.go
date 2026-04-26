package editor

import (
	"coilforge/internal/part"
	"coilforge/internal/world"
)

// ApplySelectedProp applies a property action to the first selected part and records undo on success.
func ApplySelectedProp(action part.PropAction) bool {
	if len(Selection) == 0 {
		return false
	}
	idx := Selection[0]
	if idx < 0 || idx >= len(world.Parts) {
		return false
	}
	p := world.Parts[idx]
	snap := captureSnapshot()
	if !p.ApplyProp(action) {
		return false
	}
	UndoStack = append(UndoStack, snap)
	trimUndoStackToCap()
	RedoStack = nil
	return true
}
