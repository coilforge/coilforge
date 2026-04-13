package indicator

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "indicator"

type Indicator struct {
	core.BasePart
	PinA core.PinID `json:"pinA"`
	Lit  bool       `json:"lit"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newIndicator,
		Decode: decodeIndicator,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newIndicator(id int, pos core.Pt) part.Part {
	return &Indicator{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

func decodeIndicator(data json.RawMessage) (part.Part, error) {
	var ind Indicator
	if err := json.Unmarshal(data, &ind); err != nil {
		return nil, err
	}
	if ind.TypeID == "" {
		ind.TypeID = TypeID
	}
	return &ind, nil
}

func (ind *Indicator) Base() *core.BasePart {
	return &ind.BasePart
}

func (ind *Indicator) Segments() []core.Seg {
	return nil
}

func (ind *Indicator) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *ind
	c.ID = newID
	c.PinA = allocPin()
	c.Lit = false
	return &c
}

func (ind *Indicator) MarshalJSON() ([]byte, error) {
	type indicatorJSON Indicator
	return json.Marshal((*indicatorJSON)(ind))
}
