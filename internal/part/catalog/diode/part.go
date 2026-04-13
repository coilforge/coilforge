package diode

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "diode"

type Diode struct {
	core.BasePart
	PinAnode   core.PinID `json:"pinAnode"`
	PinCathode core.PinID `json:"pinCathode"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newDiode,
		Decode: decodeDiode,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newDiode(id int, pos core.Pt) part.Part {
	return &Diode{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

func decodeDiode(data json.RawMessage) (part.Part, error) {
	var d Diode
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if d.TypeID == "" {
		d.TypeID = TypeID
	}
	return &d, nil
}

func (d *Diode) Base() *core.BasePart {
	return &d.BasePart
}

func (d *Diode) Segments() []core.Seg {
	return nil
}

func (d *Diode) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *d
	c.ID = newID
	c.PinAnode = allocPin()
	c.PinCathode = allocPin()
	return &c
}

func (d *Diode) MarshalJSON() ([]byte, error) {
	type diodeJSON Diode
	return json.Marshal((*diodeJSON)(d))
}
