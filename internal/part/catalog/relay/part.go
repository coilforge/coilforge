package relay

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"encoding/json"
)

const TypeID core.PartTypeID = "relay"

type ContactState int

const (
	ContactNC ContactState = iota
	ContactNO
)

type Pole struct {
	PinCommon core.PinID `json:"pinCommon"`
	PinNC     core.PinID `json:"pinNC"`
	PinNO     core.PinID `json:"pinNO"`
}

type Relay struct {
	core.BasePart
	PinCoilA core.PinID `json:"pinCoilA"`
	PinCoilB core.PinID `json:"pinCoilB"`

	Poles     []Pole `json:"poles"`
	PickupMs  int    `json:"pickupMs"`
	ReleaseMs int    `json:"releaseMs"`
	FlightMs  int    `json:"flightMs"`
	JitterMs  int    `json:"jitterMs"`

	CoilActive          bool           `json:"coilActive"`
	Contacts            []ContactState `json:"contacts"`
	PendingContacts     []ContactState `json:"pendingContacts"`
	TransitionDueTick   uint64         `json:"transitionDueTick"`
	TransitionScheduled bool           `json:"transitionScheduled"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		New:    newRelay,
		Decode: decodeRelay,
		Tools:  []string{"main"},
		Icon:   toolbarIcon,
	})
}

func newRelay(id int, pos core.Pt) part.Part {
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

func decodeRelay(data json.RawMessage) (part.Part, error) {
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

func (r *Relay) Base() *core.BasePart {
	return &r.BasePart
}

func (r *Relay) Segments() []core.Seg {
	return nil
}

func (r *Relay) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *r
	c.ID = newID
	c.PinCoilA = allocPin()
	c.PinCoilB = allocPin()
	c.Poles = make([]Pole, len(r.Poles))
	for i := range r.Poles {
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

func (r *Relay) MarshalJSON() ([]byte, error) {
	type relayJSON Relay
	return json.Marshal((*relayJSON)(r))
}

func (r *Relay) ensureContactSlices() {
	if len(r.Poles) == 0 {
		r.Poles = []Pole{{}}
	}
	if len(r.Contacts) != len(r.Poles) {
		r.Contacts = make([]ContactState, len(r.Poles))
	}
	if len(r.PendingContacts) != len(r.Poles) {
		r.PendingContacts = make([]ContactState, len(r.Poles))
	}
}
