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
	Dragging     bool            // Dragging stores package-level state.
	DragMoved    bool            // DragMoved is true after a non-zero move delta while dragging a part.
	DragUndoRecorded bool        // DragUndoRecorded is true after pushUndo for the current drag gesture (one undo per drag, including snap).
	DragStart    core.Pt         // DragStart stores package-level state.
	PressWorld   core.Pt         // PressWorld is mouse-down origin in world space (marquee corner / move baseline).
	PointerDownPart = -1        // PointerDownPart is index under press, or -1 when starting on empty canvas.
	MouseDownOnEmpty bool       // MouseDownOnEmpty is true when the latest press began on empty canvas (not placement).
	BoxSelecting bool            // BoxSelecting stores package-level state.
	BoxRect      core.Rect       // BoxRect stores package-level state.
	BoxSelectCrossing bool       // BoxSelectCrossing: R→L marquee on screen uses crossing (intersect) vs window (fully inside).
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
	Dragging = false
	DragMoved = false
	DragUndoRecorded = false
	DragStart = core.Pt{}
	PressWorld = core.Pt{}
	PointerDownPart = -1
	MouseDownOnEmpty = false
	BoxSelecting = false
	BoxRect = core.Rect{}
	BoxSelectCrossing = false
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
	Dragging = false
	DragMoved = false
	DragUndoRecorded = false
	BoxSelecting = false
	PointerDownPart = -1
	MouseDownOnEmpty = false
	BoxSelectCrossing = false
	LabelEditing = false
}
