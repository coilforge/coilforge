package editor

// File overview:
// snapshot captures and restores full editor/world schematic states for undo and redo.
// Subsystem: editor history.
// It works with editor state and world/part data without crossing into sim internals.
// Flow position: safety layer around mutating edit operations.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
)

type Snapshot struct {
	Parts      []part.Record // part list.
	NextPartID int           // next part id value.
	NextPinID  core.PinID    // next pin id value.
}

// pushUndo handles push undo.
func pushUndo() {
	UndoStack = append(UndoStack, captureSnapshot())
	RedoStack = nil
}

// Undo undoes its work.
func Undo() {
	if len(UndoStack) == 0 {
		return
	}
	RedoStack = append(RedoStack, captureSnapshot())
	restoreSnapshot(UndoStack[len(UndoStack)-1])
	UndoStack = UndoStack[:len(UndoStack)-1]
}

// Redo redoes its work.
func Redo() {
	if len(RedoStack) == 0 {
		return
	}
	UndoStack = append(UndoStack, captureSnapshot())
	restoreSnapshot(RedoStack[len(RedoStack)-1])
	RedoStack = RedoStack[:len(RedoStack)-1]
}

// captureSnapshot handles capture snapshot.
func captureSnapshot() Snapshot {
	records := make([]part.Record, 0, len(world.Parts))
	for _, p := range world.Parts {
		record, err := part.EncodeRecord(p)
		if err != nil {
			continue
		}
		records = append(records, record)
	}

	return Snapshot{
		Parts:      records,
		NextPartID: world.NextPartID,
		NextPinID:  world.NextPinID,
	}
}

// restoreSnapshot handles restore snapshot.
func restoreSnapshot(snapshot Snapshot) {
	parts := make([]part.Part, 0, len(snapshot.Parts))
	for _, record := range snapshot.Parts {
		p, err := part.DecodeRecord(record)
		if err != nil {
			continue
		}
		parts = append(parts, p)
	}

	world.Parts = parts
	world.NextPartID = snapshot.NextPartID
	world.NextPinID = snapshot.NextPinID
	Selection = nil
}
