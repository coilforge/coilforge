package clock

// File overview:
// props declares editable clock properties and applies updates to part state.
// Subsystem: part catalog (clock) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (self *Clock) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: self.Label},
			{Label: "On (ms)", Kind: part.PropInt, Value: clampPhaseMs(self.OnMs), Min: minPhaseMs, Max: maxPhaseMs},
			{Label: "Off (ms)", Kind: part.PropInt, Value: clampPhaseMs(self.OffMs), Min: minPhaseMs, Max: maxPhaseMs},
		},
	}
}

// ApplyProp handles apply prop.
func (self *Clock) ApplyProp(action part.PropAction) bool {
	switch action.Index {
	case 0:
		value, ok := action.NewValue.(string)
		if !ok {
			return false
		}
		self.Label = value
		return true
	case 1:
		value, ok := action.NewValue.(int)
		if !ok {
			return false
		}
		self.OnMs = clampPhaseMs(value)
		return true
	case 2:
		value, ok := action.NewValue.(int)
		if !ok {
			return false
		}
		self.OffMs = clampPhaseMs(value)
		return true
	default:
		return false
	}
}
