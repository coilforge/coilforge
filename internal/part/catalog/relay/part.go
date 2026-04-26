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
	COM1 core.PinID `json:"cOM1"`
	COM2 core.PinID `json:"cOM2"`
	COM3 core.PinID `json:"cOM3"`
	COM4 core.PinID `json:"cOM4"`
	COM5 core.PinID `json:"cOM5"`
	COM6 core.PinID `json:"cOM6"`
	COM7 core.PinID `json:"cOM7"`
	COM8 core.PinID `json:"cOM8"`
	NC1  core.PinID `json:"nC1"`
	NC2  core.PinID `json:"nC2"`
	NC3  core.PinID `json:"nC3"`
	NC4  core.PinID `json:"nC4"`
	NC5  core.PinID `json:"nC5"`
	NC6  core.PinID `json:"nC6"`
	NC7  core.PinID `json:"nC7"`
	NC8  core.PinID `json:"nC8"`
	NO1  core.PinID `json:"nO1"`
	NO2  core.PinID `json:"nO2"`
	NO3  core.PinID `json:"nO3"`
	NO4  core.PinID `json:"nO4"`
	NO5  core.PinID `json:"nO5"`
	NO6  core.PinID `json:"nO6"`
	NO7  core.PinID `json:"nO7"`
	NO8  core.PinID `json:"nO8"`
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
	// Migrate generated single-contact pins (COM/NC/NO) to pole 1 if needed.
	if relay.COM1 == 0 && relay.COM != 0 {
		relay.COM1 = relay.COM
	}
	if relay.NC1 == 0 && relay.NC != 0 {
		relay.NC1 = relay.NC
	}
	if relay.NO1 == 0 && relay.NO != 0 {
		relay.NO1 = relay.NO
	}
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
	assignNewRelayContactPins(&c, allocPin)
	c.Energized = false
	return &c
}

func assignNewRelayContactPins(self *Relay, alloc func() core.PinID) {
	self.COM1 = alloc()
	self.COM2 = alloc()
	self.COM3 = alloc()
	self.COM4 = alloc()
	self.COM5 = alloc()
	self.COM6 = alloc()
	self.COM7 = alloc()
	self.COM8 = alloc()
	self.NC1 = alloc()
	self.NC2 = alloc()
	self.NC3 = alloc()
	self.NC4 = alloc()
	self.NC5 = alloc()
	self.NC6 = alloc()
	self.NC7 = alloc()
	self.NC8 = alloc()
	self.NO1 = alloc()
	self.NO2 = alloc()
	self.NO3 = alloc()
	self.NO4 = alloc()
	self.NO5 = alloc()
	self.NO6 = alloc()
	self.NO7 = alloc()
	self.NO8 = alloc()
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
