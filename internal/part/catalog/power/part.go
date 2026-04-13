package power

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const (
	VCCTypeID core.PartTypeID = "vcc"
	GNDTypeID core.PartTypeID = "gnd"
)

type Power struct {
	core.BasePart
	Kind string     `json:"kind"`
	Pin  core.PinID `json:"pin"`
}

func init() {
	part.Register(VCCTypeID, part.TypeInfo{
		New:    newVCC,
		Decode: decodePower,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
	part.Register(GNDTypeID, part.TypeInfo{
		New:    newGND,
		Decode: decodePower,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newVCC(id int, pos core.Pt) part.Part {
	return &Power{
		BasePart: core.BasePart{ID: id, TypeID: VCCTypeID, Pos: pos},
		Kind:     "vcc",
	}
}

func newGND(id int, pos core.Pt) part.Part {
	return &Power{
		BasePart: core.BasePart{ID: id, TypeID: GNDTypeID, Pos: pos},
		Kind:     "gnd",
	}
}

func decodePower(data json.RawMessage) (part.Part, error) {
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

func (p *Power) Base() *core.BasePart {
	return &p.BasePart
}

func (p *Power) Segments() []core.Seg {
	return nil
}

func (p *Power) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *p
	c.ID = newID
	c.Pin = allocPin()
	return &c
}

func (p *Power) MarshalJSON() ([]byte, error) {
	type powerJSON Power
	return json.Marshal((*powerJSON)(p))
}
