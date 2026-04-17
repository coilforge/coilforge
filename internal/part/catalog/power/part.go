package power

// File overview:
// part defines the power part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (power).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const (
	VCCTypeID core.PartTypeID = "vcc" // vcc type id constant.
	GNDTypeID core.PartTypeID = "gnd" // gnd type id constant.
)

type Power struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	Kind          string     `json:"kind"` // kind value.
	Pin           core.PinID `json:"pin"`  // pin value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(VCCTypeID, part.TypeInfo{
		New:    newPartVCC,
		Decode: decodePart,
		Label:  "VCC",
		Tools:  []string{"main"},
		Icon:   toolbarIconVCC,
	})
	part.Register(GNDTypeID, part.TypeInfo{
		New:    newPartGND,
		Decode: decodePart,
		Label:  "GND",
		Tools:  []string{"main"},
		Icon:   toolbarIconGND,
	})
}

// newPartVCC handles new vcc part.
func newPartVCC(id int, pos core.Pt) part.Part {
	return &Power{
		BasePart: core.BasePart{ID: id, TypeID: VCCTypeID, Pos: pos},
		Kind:     "vcc",
	}
}

// newPartGND handles new gnd part.
func newPartGND(id int, pos core.Pt) part.Part {
	return &Power{
		BasePart: core.BasePart{ID: id, TypeID: GNDTypeID, Pos: pos},
		Kind:     "gnd",
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var p Power
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.TypeID == "" {
		if p.Kind == "gnd" {
			p.TypeID = GNDTypeID
		} else {
			p.TypeID = VCCTypeID
			p.Kind = "vcc"
		}
	}
	if p.Kind == "" {
		if p.TypeID == GNDTypeID {
			p.Kind = "gnd"
		} else {
			p.Kind = "vcc"
		}
	}
	return &p, nil
}

// Base handles base.
func (self *Power) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Power) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Power) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.Pin = allocPin()
	return &c
}

// MarshalJSON handles marshal json.
func (self *Power) MarshalJSON() ([]byte, error) {
	type partJSON Power
	return json.Marshal((*partJSON)(self))
}
