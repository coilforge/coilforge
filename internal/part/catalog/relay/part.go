package relay

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "relay"

type Relay struct {
	core.BasePart
	RelayPinIDs
	PoleCount    int   `json:"poleCount"`    // Number of contact poles (1..8).
	MidFlipMask  uint8 `json:"midFlipMask"`  // Per-pole flip bits (bit0=pole1 ... bit7=pole8).
	Energized    bool  `json:"-"`            // Driven by simulation in run mode; not persisted.
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "Relay",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

func newPart(id int, pos core.Pt) part.Part {
	return &Relay{
		BasePart:  core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		PoleCount: 1,
	}
}

func decodePart(data json.RawMessage) (part.Part, error) {
	var relay Relay
	if err := json.Unmarshal(data, &relay); err != nil {
		return nil, err
	}
	// Backward-compat: migrate legacy relay-wide midFlipped bool to all active poles.
	var legacyFlip struct {
		MidFlipped bool `json:"midFlipped"`
	}
	_ = json.Unmarshal(data, &legacyFlip)
	// Backward-compat: migrate legacy COM/NC/NO into pole 1 when present.
	var legacy struct {
		COM core.PinID `json:"cOM"`
		NC  core.PinID `json:"nC"`
		NO  core.PinID `json:"nO"`
	}
	_ = json.Unmarshal(data, &legacy)
	if relay.COM1 == 0 && legacy.COM != 0 {
		relay.COM1 = legacy.COM
	}
	if relay.NC1 == 0 && legacy.NC != 0 {
		relay.NC1 = legacy.NC
	}
	if relay.NO1 == 0 && legacy.NO != 0 {
		relay.NO1 = legacy.NO
	}
	if relay.TypeID == "" {
		relay.TypeID = TypeID
	}
	relay.PoleCount = clampRelayPoleCount(relay.PoleCount)
	if relay.MidFlipMask == 0 && legacyFlip.MidFlipped {
		relay.MidFlipMask = relayPoleMask(relay.PoleCount)
	}
	relay.Rotation = normalizeRelayRotation(relay.Rotation)
	return &relay, nil
}

func clampRelayPoleCount(v int) int {
	if v < 1 {
		return 1
	}
	if v > 8 {
		return 8
	}
	return v
}

func relayPoleMask(poleCount int) uint8 {
	n := clampRelayPoleCount(poleCount)
	return uint8((1 << n) - 1)
}

func normalizeRelayRotation(r int) int {
	return (r%RotationSlots + RotationSlots) % RotationSlots
}

// ClearSchematicRuntime resets draw/sim runtime fields after load or before sim runs.
func (self *Relay) ClearSchematicRuntime() {
	self.Energized = false
}

func (self *Relay) Base() *core.BasePart {
	return &self.BasePart
}

func (self *Relay) Segments() []core.Seg {
	return nil
}

func (self *Relay) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PoleCount = clampRelayPoleCount(c.PoleCount)
	assignNewRelayPins(&c, allocPin)
	c.Energized = false
	return &c
}

// ToggleMidFlip toggles relay contact-side orientation at rest.
func (self *Relay) TogglePoleFlip(pole int) bool {
	if pole < 1 || pole > self.poleCountClamped() {
		return false
	}
	bit := uint8(1 << (pole - 1))
	self.MidFlipMask ^= bit
	return true
}

func (self *Relay) MarshalJSON() ([]byte, error) {
	type partJSON Relay
	return json.Marshal((*partJSON)(self))
}
