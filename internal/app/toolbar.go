package app

// File overview:
// toolbar assembles app-side toolbar items from placement order and part metadata.
// Subsystem: app UI orchestration.
// It ties render chrome widgets to editor placement actions through part type IDs.
// Flow position: app mediation layer between input intent and renderer controls.

import (
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/flatten"
	"coilforge/internal/part"
	"coilforge/internal/partmanifest"
	"coilforge/internal/render"
	"coilforge/internal/sim"
	"coilforge/internal/world"
)

// toolbarButtons builds placement toolbar entries from partmanifest order + part metadata.
func toolbarButtons() []render.ToolButton {
	tools := make([]render.ToolButton, 0, len(partmanifest.PlacementTools))
	for _, item := range partmanifest.PlacementTools {
		label := string(item.TypeID)
		if info, ok := part.Registry[item.TypeID]; ok && info.Label != "" {
			label = info.Label
		}
		tools = append(tools, render.ToolButton{
			TypeID: string(item.TypeID),
			Hotkey: item.Hotkey,
			Label:  label,
		})
	}
	return tools
}

// rightToolbarButtons lists command-strip placeholders until real actions are wired.
func rightToolbarButtons() []render.ToolButton {
	return []render.ToolButton{
		{TypeID: "_run", Label: "Run"},
		{TypeID: "_step", Label: "Step"},
		{TypeID: "_pause", Label: "Pause"},
		{TypeID: "_save", Label: "Save"},
		{TypeID: "_load", Label: "Load"},
	}
}

// activeToolIndex returns the toolbar index for the currently selected place tool.
func activeToolIndex() int {
	for idx, item := range toolbarButtons() {
		if part.Registry[partType(item.TypeID)].New != nil && string(editor.PlaceTool) == item.TypeID {
			return idx
		}
	}
	return -1
}

// selectedPart returns the first selected part when available.
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

// ToggleRunMode switches between edit and simulation modes.
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

// partType converts a raw string ID into a strongly typed part ID.
func partType(raw string) core.PartTypeID {
	return core.PartTypeID(raw)
}
