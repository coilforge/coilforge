package clock

// File overview:
// part defines the clock part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (clock).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "clock"

type Clock struct {
	core.BasePart // BasePart carries shared part identity and transform state.
	ClockPinIDs
	OnMs  int `json:"onMs"`
	OffMs int `json:"offMs"`
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "Clock",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Clock{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		OnMs:     defaultOnMs,
		OffMs:    defaultOffMs,
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var c Clock
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.TypeID == "" {
		c.TypeID = TypeID
	}
	c.OnMs = clampPhaseMs(c.OnMs)
	c.OffMs = clampPhaseMs(c.OffMs)
	return &c, nil
}

// Base handles base.
func (self *Clock) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Clock) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Clock) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	assignNewClockPins(&c, allocPin)
	return &c
}

// MarshalJSON handles marshal json.
func (self *Clock) MarshalJSON() ([]byte, error) {
	type partJSON Clock
	return json.Marshal((*partJSON)(self))
}
