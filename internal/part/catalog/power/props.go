package power

// File overview:
// props declares editable power properties and applies updates to part state.
// Subsystem: part catalog (power) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (p *Power) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: p.Label},
		},
	}
}

// ApplyProp handles apply prop.
func (p *Power) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	p.Label = value
	return true
}
