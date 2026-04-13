package switches

import "coilforge/internal/part"

func (s *Switch) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: s.Label},
			{Label: "Momentary", Kind: part.PropBool, Value: s.Momentary},
		},
	}
}

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
