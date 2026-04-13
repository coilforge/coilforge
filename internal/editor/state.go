package editor

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
)

var (
	Selection    []int
	HoverIndex   = -1
	PlaceMode    bool
	PlaceTool    core.PartTypeID
	PlacePreview part.Part
	WireMode     bool
	WireDraft    []core.Pt
	Dragging     bool
	DragStart    core.Pt
	BoxSelecting bool
	BoxRect      core.Rect
	UndoStack    []Snapshot
	RedoStack    []Snapshot
	Clipboard    []part.Part
	LabelEditing bool
	LabelIndex   = -1
	LabelBuffer  []rune
)

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

func ClearTransient() {
	PlaceMode = false
	PlacePreview = nil
	WireMode = false
	WireDraft = nil
	Dragging = false
	BoxSelecting = false
	LabelEditing = false
}
