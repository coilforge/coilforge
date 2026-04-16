package diode

// File overview:
// props declares editable diode properties and applies updates to part state.
// Subsystem: part catalog (diode) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (d *Diode) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: d.Label},
		},
	}
}

// ApplyProp handles apply prop.
func (d *Diode) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	d.Label = value
	return true
}
