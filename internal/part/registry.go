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
	// RotationSlots: 0 = not rotatable (R does nothing). 4 = four quarter-turn layouts.
	// 8 = eight baked orientations; single selection / preview steps all 8, while group rotate advances
	// by [quarterTurnSlotSteps] (2) per 90° so 8-way parts stay aligned with the group.
	// Any other value is treated as non-rotatable.
	RotationSlots int
}

// AllowsDiscreteRotation reports whether the part type supports R-key rotation (4- or 8-step only).
func AllowsDiscreteRotation(registrySlots int) bool {
	return registrySlots == 4 || registrySlots == 8
}

// Registry stores constructors and metadata for each part type.
var Registry = map[core.PartTypeID]TypeInfo{}

// Register handles register.
func Register(id core.PartTypeID, info TypeInfo) {
	Registry[id] = info
}
