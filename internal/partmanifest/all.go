// Package partmanifest links concrete catalog packages into the binary and defines placement order.
package partmanifest

// File overview:
// all imports catalog packages into the binary via blank imports.
// Subsystem: partmanifest — catalog inclusion and app-level placement ordering.
// It depends on part/catalog/* init hooks and also defines placement order/hotkeys.
// Flow position: startup registration layer before app runtime orchestration.

import "coilforge/internal/core"

import (
	_ "coilforge/internal/part/catalog/clock"
	_ "coilforge/internal/part/catalog/diode"
	_ "coilforge/internal/part/catalog/indicator"
	_ "coilforge/internal/part/catalog/power"
	_ "coilforge/internal/part/catalog/rch"
	_ "coilforge/internal/part/catalog/relay"
	_ "coilforge/internal/part/catalog/switches"
	_ "coilforge/internal/part/catalog/wire"
)

type PlacementTool struct {
	TypeID core.PartTypeID // Part type selected by this toolbar/hotkey entry.
	Hotkey rune            // Keyboard shortcut displayed in this placement order slot.
}

// PlacementTools is the single app-level placement order and hotkey list.
var PlacementTools = []PlacementTool{
	{TypeID: "relay", Hotkey: '1'},
	{TypeID: "vcc", Hotkey: '2'},
	{TypeID: "gnd", Hotkey: '3'},
	{TypeID: "switch", Hotkey: '4'},
	{TypeID: "indicator", Hotkey: '5'},
	{TypeID: "diode", Hotkey: '6'},
	{TypeID: "rch", Hotkey: '7'},
	{TypeID: "clock", Hotkey: '8'},
	{TypeID: "wire", Hotkey: 'W'},
}
