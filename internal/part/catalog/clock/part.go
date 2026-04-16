package clock

// File overview:
// part defines the clock part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (clock).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "clock"

type Clock struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinOut        core.PinID `json:"pinOut"`     // pin out value.
	PeriodTick    int        `json:"periodTick"` // period tick value.
	HighTick      int        `json:"highTick"`   // high tick value.
	OutputHigh    bool       `json:"outputHigh"` // output high value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newClock,
		Decode: decodeClock,
		Label:  "Clock",
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

// newClock handles new clock.
func newClock(id int, pos core.Pt) part.Part {
	return &Clock{
		BasePart:   core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		PeriodTick: 1000,
		HighTick:   500,
	}
}

// decodeClock handles decode clock.
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

// Base handles base.
func (c *Clock) Base() *core.BasePart {
	return &c.BasePart
}

// Segments handles segments.
func (c *Clock) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (c *Clock) Clone(newID int, allocPin func() core.PinID) part.Part {
	clone := *c
	clone.ID = newID
	clone.PinOut = allocPin()
	clone.OutputHigh = false
	return &clone
}

// MarshalJSON handles marshal json.
func (c *Clock) MarshalJSON() ([]byte, error) {
	type clockJSON Clock
	return json.Marshal((*clockJSON)(c))
}
