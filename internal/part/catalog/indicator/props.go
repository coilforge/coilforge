package indicator

import "coilforge/internal/part"

func (ind *Indicator) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: ind.Label},
		},
	}
}

func (ind *Indicator) ApplyProp(action part.PropAction) bool {
	if action.Index != 0 {
		return false
	}
	value, ok := action.NewValue.(string)
	if !ok {
		return false
	}
	ind.Label = value
	return true
}
