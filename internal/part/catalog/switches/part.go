package switches

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "switch"

type Switch struct {
	core.BasePart
	PinA      core.PinID `json:"pinA"`
	PinB      core.PinID `json:"pinB"`
	Closed    bool       `json:"closed"`
	Momentary bool       `json:"momentary"`
	Pressed   bool       `json:"pressed"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newSwitch,
		Decode: decodeSwitch,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newSwitch(id int, pos core.Pt) part.Part {
	return &Switch{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

func decodeSwitch(data json.RawMessage) (part.Part, error) {
	var s Switch
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.TypeID == "" {
		s.TypeID = TypeID
	}
	return &s, nil
}

func (s *Switch) Base() *core.BasePart {
	return &s.BasePart
}

func (s *Switch) Segments() []core.Seg {
	return nil
}

func (s *Switch) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *s
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.Pressed = false
	return &c
}

func (s *Switch) MarshalJSON() ([]byte, error) {
	type switchJSON Switch
	return json.Marshal((*switchJSON)(s))
}
