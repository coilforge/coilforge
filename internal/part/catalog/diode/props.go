package diode

import "coilforge/internal/part"

func (d *Diode) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: d.Label},
		},
	}
}

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
