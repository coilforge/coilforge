package relay

// File overview:
// part defines the relay part type, registration hooks, and clone/decode behavior.
// Subsystem: part catalog (relay).
// It works with shared part/core contracts and is complemented by draw/props/sim/assets files.
// Flow position: concrete catalog part implementation loaded through part registry.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "relay"

type ContactState int

const (
	ContactNC ContactState = iota // contact nc constant.
	ContactNO                     // contact no constant.
)

type Pole struct {
	PinCommon core.PinID `json:"pinCommon"` // pin common value.
	PinNC     core.PinID `json:"pinNC"`     // pin nc value.
	PinNO     core.PinID `json:"pinNO"`     // pin no value.
}

type Relay struct {
	core.BasePart            // BasePart carries shared part identity and transform state.
	PinCoilA      core.PinID `json:"pinCoilA"` // pin coil a value.
	PinCoilB      core.PinID `json:"pinCoilB"` // pin coil b value.

	Poles     []Pole `json:"poles"`     // poles value.
	PickupMs  int    `json:"pickupMs"`  // pickup ms value.
	ReleaseMs int    `json:"releaseMs"` // release ms value.
	FlightMs  int    `json:"flightMs"`  // flight ms value.
	JitterMs  int    `json:"jitterMs"`  // jitter ms value.

	CoilActive          bool           `json:"coilActive"`          // coil active value.
	Contacts            []ContactState `json:"contacts"`            // contacts value.
	PendingContacts     []ContactState `json:"pendingContacts"`     // pending contacts value.
	TransitionDueTick   uint64         `json:"transitionDueTick"`   // transition due tick value.
	TransitionScheduled bool           `json:"transitionScheduled"` // transition scheduled value.
}

// init registers the part type with the global registry.
func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newPart,
		Decode: decodePart,
		Label:  "Relay",
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

// newPart handles new part.
func newPart(id int, pos core.Pt) part.Part {
	r := &Relay{
		BasePart:  core.BasePart{ID: id, TypeID: TypeID, Pos: pos},
		Poles:     []Pole{{}},
		PickupMs:  5,
		ReleaseMs: 3,
		FlightMs:  1,
		JitterMs:  0,
	}
	r.ensureContactSlices()
	return r
}

// decodePart handles decode part.
func decodePart(data json.RawMessage) (part.Part, error) {
	var r Relay
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if r.TypeID == "" {
		r.TypeID = TypeID
	}
	if len(r.Poles) == 0 {
		r.Poles = []Pole{{}}
	}
	r.ensureContactSlices()
	return &r, nil
}

// Base handles base.
func (self *Relay) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Relay) Segments() []core.Seg {
	return nil
}

// Clone handles clone.
func (self *Relay) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinCoilA = allocPin()
	c.PinCoilB = allocPin()
	c.Poles = make([]Pole, len(self.Poles))
	for i := range self.Poles {
		c.Poles[i] = Pole{
			PinCommon: allocPin(),
			PinNC:     allocPin(),
			PinNO:     allocPin(),
		}
	}
	c.CoilActive = false
	c.TransitionDueTick = 0
	c.TransitionScheduled = false
	c.ensureContactSlices()
	return &c
}

// MarshalJSON handles marshal json.
func (self *Relay) MarshalJSON() ([]byte, error) {
	type partJSON Relay
	return json.Marshal((*partJSON)(self))
}

// ensureContactSlices handles ensure contact slices.
func (self *Relay) ensureContactSlices() {
	if len(self.Poles) == 0 {
		self.Poles = []Pole{{}}
	}
	if len(self.Contacts) != len(self.Poles) {
		self.Contacts = make([]ContactState, len(self.Poles))
	}
	if len(self.PendingContacts) != len(self.Poles) {
		self.PendingContacts = make([]ContactState, len(self.Poles))
	}
}
