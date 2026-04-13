package skeleton

import "coilforge/internal/part"

func (t *Template) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: t.Label},
		},
	}
}

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
