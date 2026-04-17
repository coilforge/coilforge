package skeleton

// File overview:
// part defines the skeleton part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (skeleton).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// Template is a non-registered reference part for new catalog entries.
type Template struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinA          core.PinID `json:"pinA"` // pin a value.
	PinB          core.PinID `json:"pinB"` // pin b value.
}

// New constructs its work.
func New(id int, pos core.Pt) part.Part {
	return &Template{
		BasePart: core.BasePart{ID: id, TypeID: "template", Pos: pos},
	}
}

// Decode handles decode.
func Decode(data json.RawMessage) (part.Part, error) {
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Base handles base.
func (self *Template) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Template) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Template) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	return &c
}

// MarshalJSON handles marshal json.
func (self *Template) MarshalJSON() ([]byte, error) {
	type partJSON Template
	return json.Marshal((*partJSON)(self))
}
