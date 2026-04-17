package switches

// File overview:
// part defines the switches part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (switches).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "switch"

type Switch struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinA          core.PinID `json:"pinA"`      // pin a value.
	PinB          core.PinID `json:"pinB"`      // pin b value.
	Closed        bool       `json:"closed"`    // closed value.
	Momentary     bool       `json:"momentary"` // momentary value.
	Pressed       bool       `json:"pressed"`   // pressed value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newPart,
		Decode: decodePart,
		Label:  "Switch",
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Switch{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var s Switch
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.TypeID == "" {
		s.TypeID = TypeID
	}
	return &s, nil
}

// Base handles base.
func (self *Switch) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Switch) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Switch) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.Pressed = false
	return &c
}

// MarshalJSON handles marshal json.
func (self *Switch) MarshalJSON() ([]byte, error) {
	type partJSON Switch
	return json.Marshal((*partJSON)(self))
}
