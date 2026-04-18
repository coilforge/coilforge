package indicator

// File overview:
// part defines the indicator part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (indicator).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "indicator"

type Indicator struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinA          core.PinID `json:"pinA"` // pin a value.
	PinB          core.PinID `json:"pinB"` // pin b value.
	Lit           bool       `json:"lit"`  // lit value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "Indicator",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Indicator{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var ind Indicator
	if err := json.Unmarshal(data, &ind); err != nil {
		return nil, err
	}
	if ind.TypeID == "" {
		ind.TypeID = TypeID
	}
	return &ind, nil
}

// Base handles base.
func (self *Indicator) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Indicator) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Indicator) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.Lit = false
	return &c
}

// MarshalJSON handles marshal json.
func (self *Indicator) MarshalJSON() ([]byte, error) {
	type partJSON Indicator
	return json.Marshal((*partJSON)(self))
}
