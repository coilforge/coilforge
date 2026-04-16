package rch

// File overview:
// part defines the rch part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (rch).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "rch"

type RCH struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinIn         core.PinID `json:"pinIn"`   // pin in value.
	PinOut        core.PinID `json:"pinOut"`  // pin out value.
	DelayMs       int        `json:"delayMs"` // delay ms value.
	Active        bool       `json:"active"`  // active value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newRCH,
		Decode: decodeRCH,
		Label:  "RCH",
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

// newRCH handles new rch.
func newRCH(id int, pos core.Pt) part.Part {
	return &RCH{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		DelayMs:  10,
	}
}

// decodeRCH handles decode rch.
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

// Base handles base.
func (r *RCH) Base() *core.BasePart {
	return &r.BasePart
}

// Segments handles segments.
func (r *RCH) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (r *RCH) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *r
	c.ID = newID
	c.PinIn = allocPin()
	c.PinOut = allocPin()
	c.Active = false
	return &c
}

// MarshalJSON handles marshal json.
func (r *RCH) MarshalJSON() ([]byte, error) {
	type rchJSON RCH
	return json.Marshal((*rchJSON)(r))
}
