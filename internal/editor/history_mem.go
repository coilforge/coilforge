package editor

// File overview:
// history_mem estimates heap retained by undo/redo snapshot stacks for debugging only.
// Subsystem: editor history.
// It sums slice capacities and json.RawMessage backing allocations from encoded parts.

import (
	"coilforge/internal/part"
	"unsafe"
)

// UndoRedoStacksApproxBytes returns a rough byte estimate for the undo and redo snapshot stacks.
// It sums cap([]Snapshot), cap([]part.Record) per snapshot, and cap(json.RawMessage) per record,
// plus the logical byte length of each part type id string. Intended for debugging only; not exact heap.
func UndoRedoStacksApproxBytes() (undoBytes, redoBytes int64) {
	return stackApproxBytes(UndoStack), stackApproxBytes(RedoStack)
}

func stackApproxBytes(stack []Snapshot) int64 {
	var n int64
	n += int64(uintptr(cap(stack)) * unsafe.Sizeof(Snapshot{}))
	for _, snap := range stack {
		n += snapshotApproxBytes(snap)
	}
	return n
}

func snapshotApproxBytes(s Snapshot) int64 {
	var n int64
	n += int64(uintptr(cap(s.Parts)) * unsafe.Sizeof(part.Record{}))
	for _, r := range s.Parts {
		n += int64(len(r.Type))
		n += int64(cap(r.Data))
	}
	return n
}
