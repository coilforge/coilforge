package editor

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
)

type Snapshot struct {
	Parts      []part.Record
	NextPartID int
	NextPinID  core.PinID
}

func pushUndo() {
	UndoStack = append(UndoStack, captureSnapshot())
	RedoStack = nil
}

func Undo() {
	if len(UndoStack) == 0 {
		return
	}
	RedoStack = append(RedoStack, captureSnapshot())
	restoreSnapshot(UndoStack[len(UndoStack)-1])
	UndoStack = UndoStack[:len(UndoStack)-1]
}

func Redo() {
	if len(RedoStack) == 0 {
		return
	}
	UndoStack = append(UndoStack, captureSnapshot())
	restoreSnapshot(RedoStack[len(RedoStack)-1])
	RedoStack = RedoStack[:len(RedoStack)-1]
}

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
