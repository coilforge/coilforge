package wire

// File overview:
// part defines the wire part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (wire).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "wire"

type Wire struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	Half          core.Pt    `json:"half"`  // half value.
	PinA          core.PinID `json:"pinA"`  // pin a value.
	PinB          core.PinID `json:"pinB"`  // pin b value.
	State         int        `json:"state"` // current state.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:     newPart,
		NewWire: newSegment,
		Decode:  decodePart,
		Label:   "Wire",
		Tools:   []string{"wire"},
		Icon:    toolbarIcon,
	})
}

// New constructs its work.
func New(id int, from, to core.Pt, allocPinA, allocPinB func() core.PinID) *Wire {
	wire := &Wire{
		BasePart: core.BasePart{
			ID:     id,
			TypeID: TypeID,
			Pos: core.Pt{
				X: (from.X + to.X) / 2,
				Y: (from.Y + to.Y) / 2,
			},
		},
		Half: core.Pt{
			X: (to.X - from.X) / 2,
			Y: (to.Y - from.Y) / 2,
		},
	}
	if allocPinA != nil {
		wire.PinA = allocPinA()
	}
	if allocPinB != nil {
		wire.PinB = allocPinB()
	}
	return wire
}

// newSegment handles new segment.
func newSegment(id int, from, to core.Pt, allocPin func() core.PinID) part.Part {
	return New(id, from, to, allocPin, allocPin)
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	return &Wire{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		Half:     core.Pt{X: 16, Y: 0},
	}
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var w Wire
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, err
	}
	if w.TypeID == "" {
		w.TypeID = TypeID
	}
	return &w, nil
}

// Base handles base.
func (self *Wire) Base() *core.BasePart {
	return &self.BasePart
}

// Clone handles clone.
func (self *Wire) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.State = core.NetFloat
	return &c
}

// MarshalJSON handles marshal json.
func (self *Wire) MarshalJSON() ([]byte, error) {
	type partJSON Wire
	return json.Marshal((*partJSON)(self))
}
