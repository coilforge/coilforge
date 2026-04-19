package gnd

// File overview:
// part defines the gnd part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (gnd).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "gnd"

type Gnd struct {
	core.BasePart // BasePart carries shared part identity and transform state.
	GndPinIDs
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "GND",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Gnd{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var g Gnd
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}
	if g.TypeID == "" {
		g.TypeID = TypeID
	}
	return &g, nil
}

// Base handles base.
func (self *Gnd) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Gnd) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Gnd) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	assignNewGndPins(&c, allocPin)
	return &c
}

// MarshalJSON handles marshal json.
func (self *Gnd) MarshalJSON() ([]byte, error) {
	type partJSON Gnd
	return json.Marshal((*partJSON)(self))
}
