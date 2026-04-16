package diode

// File overview:
// part defines the diode part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (diode).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "diode"

type Diode struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinAnode      core.PinID `json:"pinAnode"`   // pin anode value.
	PinCathode    core.PinID `json:"pinCathode"` // pin cathode value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newDiode,
		Decode: decodeDiode,
		Label:  "Diode",
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

// newDiode handles new diode.
func newDiode(id int, pos core.Pt) part.Part {
	return &Diode{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
	}
}

// decodeDiode handles decode diode.
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

// Base handles base.
func (d *Diode) Base() *core.BasePart {
	return &d.BasePart
}

// Segments handles segments.
func (d *Diode) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (d *Diode) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *d
	c.ID = newID
	c.PinAnode = allocPin()
	c.PinCathode = allocPin()
	return &c
}

// MarshalJSON handles marshal json.
func (d *Diode) MarshalJSON() ([]byte, error) {
	type diodeJSON Diode
	return json.Marshal((*diodeJSON)(d))
}
