package app

import (
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/flatten"
	"coilforge/internal/part"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/world"
)

func toolbarButtons() []render.ToolButton {
	return []render.ToolButton{
		{TypeID: "relay", Hotkey: '1', Label: "Relay"},
		{TypeID: "vcc", Hotkey: '2', Label: "VCC"},
		{TypeID: "gnd", Hotkey: '3', Label: "GND"},
		{TypeID: "switch", Hotkey: '4', Label: "Switch"},
		{TypeID: "indicator", Hotkey: '5', Label: "Indicator"},
		{TypeID: "diode", Hotkey: '6', Label: "Diode"},
		{TypeID: "rch", Hotkey: '7', Label: "RCH"},
		{TypeID: "clock", Hotkey: '8', Label: "Clock"},
		{TypeID: "wire", Hotkey: 'W', Label: "Wire"},
	}
}

func activeToolIndex() int {
	for idx, item := range toolbarButtons() {
		if part.Registry[partType(item.TypeID)].New != nil && string(editor.PlaceTool) == item.TypeID {
			return idx
		}
	}
	return -1
}

func selectedPart() part.Part {
	if len(editor.Selection) == 0 {
		return nil
	}
	idx := editor.Selection[0]
	if idx < 0 || idx >= len(world.Parts) {
		return nil
	}
	return world.Parts[idx]
}

func ToggleRunMode() {
	if world.RunMode {
		sim.Stop()
		world.RunMode = false
		return
	}

	editor.ClearTransient()
	flatten.BuildNets()
	sim.Start()
	world.RunMode = true
}

func partType(raw string) core.PartTypeID {
	return core.PartTypeID(raw)
}
