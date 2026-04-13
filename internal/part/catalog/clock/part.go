package clock

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "clock"

type Clock struct {
	core.BasePart
	PinOut     core.PinID `json:"pinOut"`
	PeriodTick int        `json:"periodTick"`
	HighTick   int        `json:"highTick"`
	OutputHigh bool       `json:"outputHigh"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newClock,
		Decode: decodeClock,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newClock(id int, pos core.Pt) part.Part {
	return &Clock{
		BasePart:   core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		PeriodTick: 1000,
		HighTick:   500,
	}
}

func decodeClock(data json.RawMessage) (part.Part, error) {
	var c Clock
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.TypeID == "" {
		c.TypeID = TypeID
	}
	if c.PeriodTick <= 0 {
		c.PeriodTick = 1000
	}
	if c.HighTick <= 0 {
		c.HighTick = c.PeriodTick / 2
	}
	return &c, nil
}

func (c *Clock) Base() *core.BasePart {
	return &c.BasePart
}

func (c *Clock) Segments() []core.Seg {
	return nil
}

func (c *Clock) Clone(newID int, allocPin func() core.PinID) part.Part {
	clone := *c
	clone.ID = newID
	clone.PinOut = allocPin()
	clone.OutputHigh = false
	return &clone
}

func (c *Clock) MarshalJSON() ([]byte, error) {
	type clockJSON Clock
	return json.Marshal((*clockJSON)(c))
}
