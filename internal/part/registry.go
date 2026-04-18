package part

// File overview:
// registry stores constructors and metadata for all registered part types.
// Subsystem: part registration.
// Catalog packages register in init; app/editor lookup types through this shared table.
// Flow position: runtime type factory between startup registration and part creation.

import (
	"coilforge/internal/core"
	"encoding/json"

	"github.com/hajimehoshi/ebiten/v2"
)

type TypeInfo struct {
	New     func(id int, pos core.Pt) Part                                  // new value.
	NewWire func(id int, from, to core.Pt, allocPin func() core.PinID) Part // new wire value.
	Decode  func(data json.RawMessage) (Part, error)                        // decode value.
	Label   string                                                          // display name used in editor chrome.
	Tools   []string                                                        // tools value.
	Icon    func() *ebiten.Image                                            // icon value.
	// RotationSlots is how many discrete rotations R cycles through (typically 4, or 8 when _45 SVGs exist). Zero means 4.
	RotationSlots int
}

// Registry stores constructors and metadata for each part type.
var Registry = map[core.PartTypeID]TypeInfo{}

// Register handles register.
func Register(id core.PartTypeID, info TypeInfo) {
	Registry[id] = info
}
