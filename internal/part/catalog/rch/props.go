package rch

// File overview:
// props declares editable rch properties and applies updates to part state.
// Subsystem: part catalog (rch) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (r *RCH) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: r.Label},
			{Label: "Delay", Kind: part.PropInt, Value: r.DelayMs, Min: 1, Max: 10000},
		},
	}
}

// ApplyProp handles apply prop.
func (r *RCH) ApplyProp(action part.PropAction) bool {
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
		if !ok || value <= 0 {
			return false
		}
		r.DelayMs = value
		return true
	default:
		return false
	}
}
