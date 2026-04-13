package skeleton

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// Template is a non-registered reference component for new catalog entries.
type Template struct {
	core.BasePart
	PinA core.PinID `json:"pinA"`
	PinB core.PinID `json:"pinB"`
}

func New(id int, pos core.Pt) part.Part {
	return &Template{
		BasePart: core.BasePart{ID: id, TypeID: "template", Pos: pos},
	}
}

func Decode(data json.RawMessage) (part.Part, error) {
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *Template) Base() *core.BasePart {
	return &t.BasePart
}

func (t *Template) Segments() []core.Seg {
	return nil
}

func (t *Template) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *t
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	return &c
}

func (t *Template) MarshalJSON() ([]byte, error) {
	type templateJSON Template
	return json.Marshal((*templateJSON)(t))
}
