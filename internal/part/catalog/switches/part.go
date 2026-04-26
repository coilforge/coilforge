package switches

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "switches"

const (
	ModeToggle = "toggle"
	ModeMomentary = "momentary"
)

type Switches struct {
	core.BasePart
	SwitchesPinIDs
	On   bool   `json:"on"`
	Mode string `json:"mode"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:           newPart,
		Decode:        decodePart,
		Label:         "Switch",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: RotationSlots,
	})
}

func newPart(id int, pos core.Pt) part.Part {
	return &Switches{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		Mode:     ModeToggle,
	}
}

func decodePart(data json.RawMessage) (part.Part, error) {
	var s Switches
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.TypeID == "" {
		s.TypeID = TypeID
	}
	if s.Mode != ModeMomentary {
		s.Mode = ModeToggle
	}
	s.Rotation = normalizeRotation(s.Rotation)
	return &s, nil
}

func normalizeRotation(r int) int {
	return (r%RotationSlots + RotationSlots) % RotationSlots
}

func (self *Switches) Base() *core.BasePart {
	return &self.BasePart
}

func (self *Switches) Segments() []core.Seg {
	return nil
}

func (self *Switches) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	assignNewSwitchesPins(&c, allocPin)
	return &c
}

func (self *Switches) MarshalJSON() ([]byte, error) {
	type partJSON Switches
	return json.Marshal((*partJSON)(self))
}
