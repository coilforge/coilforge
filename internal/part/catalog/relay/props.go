package relay

// File overview:
// props declares editable relay properties and applies updates to part state.
// Subsystem: part catalog (relay) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (r *Relay) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: r.Label},
			{Label: "Poles", Kind: part.PropInt, Value: len(r.Poles), Min: 1, Max: 8},
			{Label: "Pickup", Kind: part.PropInt, Value: r.PickupMs, Min: 0, Max: 1000},
			{Label: "Release", Kind: part.PropInt, Value: r.ReleaseMs, Min: 0, Max: 1000},
		},
	}
}

// ApplyProp handles apply prop.
func (r *Relay) ApplyProp(action part.PropAction) bool {
	switch action.Index {
	case 0:
		value, ok := action.NewValue.(string)
		if !ok {
			return false
		}
		r.Label = value
		return true
	case 1:
		value, ok := action.NewValue.(int)
		if !ok || value < 1 {
			return false
		}
		poles := make([]Pole, value)
		for i := range poles {
			if i < len(r.Poles) {
				poles[i] = r.Poles[i]
			}
		}
		r.Poles = poles
		r.ensureContactSlices()
		return true
	case 2:
		value, ok := action.NewValue.(int)
		if !ok || value < 0 {
			return false
		}
		r.PickupMs = value
		return true
	case 3:
		value, ok := action.NewValue.(int)
		if !ok || value < 0 {
			return false
		}
		r.ReleaseMs = value
		return true
	default:
		return false
	}
}
