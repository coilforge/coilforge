package skeleton

// File overview:
// props declares editable skeleton properties and applies updates to part state.
// Subsystem: part catalog (skeleton) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (t *Template) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: t.Label},
		},
	}
}

// ApplyProp handles apply prop.
func (t *Template) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	t.Label = value
	return true
}
