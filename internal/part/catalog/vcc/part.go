package vcc

// File overview:
// part defines the vcc part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (vcc).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "vcc"

type Vcc struct {
	core.BasePart // BasePart carries shared part identity and transform state.
	VccPinIDs
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "VCC",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Vcc{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var v Vcc
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	if v.TypeID == "" {
		v.TypeID = TypeID
	}
	return &v, nil
}

// Base handles base.
func (self *Vcc) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Vcc) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Vcc) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	assignNewVccPins(&c, allocPin)
	return &c
}

// MarshalJSON handles marshal json.
func (self *Vcc) MarshalJSON() ([]byte, error) {
	type partJSON Vcc
	return json.Marshal((*partJSON)(self))
}
