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
	Energized bool `json:"-"` // Driven by simulation in run mode; not persisted.
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
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

func decodePart(data json.RawMessage) (part.Part, error) {
	var relay Relay
	if err := json.Unmarshal(data, &relay); err != nil {
		return nil, err
	}
	if relay.TypeID == "" {
		relay.TypeID = TypeID
	}
	relay.Rotation = normalizeRelayRotation(relay.Rotation)
	return &relay, nil
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
	assignNewRelayPins(&c, allocPin)
	c.Energized = false
	return &c
}

func (self *Relay) MarshalJSON() ([]byte, error) {
	type partJSON Relay
	return json.Marshal((*partJSON)(self))
}
