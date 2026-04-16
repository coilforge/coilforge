package app

// File overview:
// fileio owns project save/load entrypoints and snapshot persistence wiring.
// Subsystem: app I/O orchestration.
// It bridges world state with part codec and avoids file handling in editor/sim.
// Flow position: app-level persistence boundary invoked by hotkeys and menus.

import (
	"coilforge/internal/core"
	"coilforge/internal/editor"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"encoding/json"
	"os"
)

type FileFormat struct {
	NextPartID int           `json:"nextPartID"` // next part id value.
	NextPinID  core.PinID    `json:"nextPinID"`  // next pin id value.
	Parts      []part.Record `json:"parts"`      // part list.
}

// SaveProject saves project.
func SaveProject(path string) error {
	records := make([]part.Record, 0, len(world.Parts))
	for _, p := range world.Parts {
		record, err := part.EncodeRecord(p)
		if err != nil {
			return err
		}
		records = append(records, record)
	}

	file := FileFormat{
		NextPartID: world.NextPartID,
		NextPinID:  world.NextPinID,
		Parts:      records,
	}

	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// LoadProject loads project.
func LoadProject(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var file FileFormat
	if err := json.Unmarshal(data, &file); err != nil {
		return err
	}

	parts := make([]part.Part, 0, len(file.Parts))
	for _, record := range file.Parts {
		p, err := part.DecodeRecord(record)
		if err != nil {
			return err
		}
		parts = append(parts, p)
	}

	world.Parts = parts
	world.NextPartID = file.NextPartID
	world.NextPinID = file.NextPinID
	editor.Reset()
	return nil
}
