package wire

import "coilforge/internal/part"

func (w *Wire) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: w.Label},
		},
	}
}

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
