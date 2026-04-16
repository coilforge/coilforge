package part

// File overview:
// props defines typed property descriptors and values exposed by part implementations.
// Subsystem: part property contract.
// It is consumed by catalog props files and app/render property editing UI.
// Flow position: metadata bridge between part internals and editor controls.

type PropSpec struct {
	Items []PropItem // items value.
}

type PropItem struct {
	Label   string   // label text.
	Kind    int      // kind value.
	Value   any      // current value.
	Choices []string // choices value.
	Min     int      // min value.
	Max     int      // max value.
}

type PropAction struct {
	Index    int // index value.
	NewValue any // new value value.
}

const (
	PropText         = iota // PropText renders a text property row.
	PropInt                 // PropInt renders an integer property row.
	PropChoice              // PropChoice renders a choice property row.
	PropBool                // PropBool renders a boolean property row.
	PropActionButton        // PropActionButton renders an action button row.
)
