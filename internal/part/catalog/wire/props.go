package wire

// File overview:
// props declares editable wire properties and applies updates to part state.
// Subsystem: part catalog (wire) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (w *Wire) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: w.Label},
		},
	}
}

// ApplyProp handles apply prop.
func (w *Wire) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	w.Label = value
	return true
}
