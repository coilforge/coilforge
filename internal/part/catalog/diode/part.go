package diode

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "diode"

type Diode struct {
	core.BasePart
	DiodePinIDs
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "Diode",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

func newPart(id int, pos core.Pt) part.Part {
	return &Diode{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

func decodePart(data json.RawMessage) (part.Part, error) {
	var d Diode
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	if d.TypeID == "" {
		d.TypeID = TypeID
	}
	d.Rotation = normalizeRotation(d.Rotation)
	return &d, nil
}

func normalizeRotation(r int) int {
	return (r%RotationSlots + RotationSlots) % RotationSlots
}

func (self *Diode) Base() *core.BasePart {
	return &self.BasePart
}

func (self *Diode) Segments() []core.Seg {
	return nil
}

func (self *Diode) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	assignNewDiodePins(&c, allocPin)
	return &c
}

func (self *Diode) MarshalJSON() ([]byte, error) {
	type partJSON Diode
	return json.Marshal((*partJSON)(self))
}
