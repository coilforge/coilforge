package indicator

// File overview:
// props declares editable indicator properties and applies updates to part state.
// Subsystem: part catalog (indicator) properties.
// It implements part property contracts used by app/editor property panels.
// Flow position: part-specific metadata and mutation rules in edit flow.

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (self *Indicator) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: self.Label},
		},
	}
}

// ApplyProp handles apply prop.
func (self *Indicator) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	self.Label = value
	return true
}
