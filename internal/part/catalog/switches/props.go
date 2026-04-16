package switches

// File overview:
// props declares editable switches properties and applies updates to part state.
// Subsystem: part catalog (switches) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (s *Switch) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: s.Label},
			{Label: "Momentary", Kind: part.PropBool, Value: s.Momentary},
		},
	}
}

// ApplyProp handles apply prop.
func (s *Switch) ApplyProp(action part.PropAction) bool {
	switch action.Index {
	case 0:
		value, ok := action.NewValue.(string)
		if !ok {
			return false
		}
		s.Label = value
		return true
	case 1:
		value, ok := action.NewValue.(bool)
		if !ok {
			return false
		}
		s.Momentary = value
		return true
	default:
		return false
	}
}
