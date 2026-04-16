package editor

// File overview:
// state defines editor package state, tool modes, and selection bookkeeping.
// Subsystem: editor state management.
// It is shared by editor handlers in this package and persisted via snapshot helpers.
// Flow position: editor-internal backing store behind app-dispatched actions.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

var (
	Selection    []int           // Selection stores package-level state.
	HoverIndex   = -1            // HoverIndex stores package-level state.
	PlaceMode    bool            // PlaceMode stores package-level state.
	PlaceTool    core.PartTypeID // PlaceTool stores package-level state.
	PlacePreview part.Part       // PlacePreview stores package-level state.
	WireMode     bool            // WireMode stores package-level state.
	WireDraft    []core.Pt       // WireDraft stores package-level state.
	Dragging     bool            // Dragging stores package-level state.
	DragStart    core.Pt         // DragStart stores package-level state.
	BoxSelecting bool            // BoxSelecting stores package-level state.
	BoxRect      core.Rect       // BoxRect stores package-level state.
	UndoStack    []Snapshot      // UndoStack stores package-level state.
	RedoStack    []Snapshot      // RedoStack stores package-level state.
	Clipboard    []part.Part     // Clipboard stores package-level state.
	LabelEditing bool            // LabelEditing stores package-level state.
	LabelIndex   = -1            // LabelIndex stores package-level state.
	LabelBuffer  []rune          // LabelBuffer stores package-level state.
)

// Reset resets its work.
func Reset() {
	Selection = nil
	HoverIndex = -1
	PlaceMode = false
	PlaceTool = ""
	PlacePreview = nil
	WireMode = false
	WireDraft = nil
	Dragging = false
	DragStart = core.Pt{}
	BoxSelecting = false
	BoxRect = core.Rect{}
	UndoStack = nil
	RedoStack = nil
	Clipboard = nil
	LabelEditing = false
	LabelIndex = -1
	LabelBuffer = nil
}

// ClearTransient handles clear transient.
func ClearTransient() {
	PlaceMode = false
	PlacePreview = nil
	WireMode = false
	WireDraft = nil
	Dragging = false
	BoxSelecting = false
	LabelEditing = false
}
