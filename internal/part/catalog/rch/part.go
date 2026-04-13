package rch

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "rch"

type RCH struct {
	core.BasePart
	PinIn   core.PinID `json:"pinIn"`
	PinOut  core.PinID `json:"pinOut"`
	DelayMs int        `json:"delayMs"`
	Active  bool       `json:"active"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newRCH,
		Decode: decodeRCH,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newRCH(id int, pos core.Pt) part.Part {
	return &RCH{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		DelayMs:  10,
	}
}

func decodeRCH(data json.RawMessage) (part.Part, error) {
	var r RCH
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if r.TypeID == "" {
		r.TypeID = TypeID
	}
	if r.DelayMs <= 0 {
		r.DelayMs = 10
	}
	return &r, nil
}

func (r *RCH) Base() *core.BasePart {
	return &r.BasePart
}

func (r *RCH) Segments() []core.Seg {
	return nil
}

func (r *RCH) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *r
	c.ID = newID
	c.PinIn = allocPin()
	c.PinOut = allocPin()
	c.Active = false
	return &c
}

func (r *RCH) MarshalJSON() ([]byte, error) {
	type rchJSON RCH
	return json.Marshal((*rchJSON)(r))
}
