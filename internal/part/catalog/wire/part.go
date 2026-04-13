package wire

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "wire"

type Wire struct {
	core.BasePart
	Half  core.Pt    `json:"half"`
	PinA  core.PinID `json:"pinA"`
	PinB  core.PinID `json:"pinB"`
	State int        `json:"state"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newWire,
		Decode: decodeWire,
		Tools:  []string{"wire"},
		Icon:   toolbarIcon,
	})
}

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

func newWire(id int, pos core.Pt) part.Part {
	return &Wire{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		Half:     core.Pt{X: 16, Y: 0},
	}
}

func decodeWire(data json.RawMessage) (part.Part, error) {
	var w Wire
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, err
	}
	if w.TypeID == "" {
		w.TypeID = TypeID
	}
	return &w, nil
}

func (w *Wire) Base() *core.BasePart {
	return &w.BasePart
}

func (w *Wire) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *w
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.State = core.NetFloat
	return &c
}

func (w *Wire) MarshalJSON() ([]byte, error) {
	type wireJSON Wire
	return json.Marshal((*wireJSON)(w))
}
