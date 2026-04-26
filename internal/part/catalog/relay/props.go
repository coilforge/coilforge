package relay

import "coilforge/internal/part"

func (self *Relay) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: self.Label},
			{Label: "Poles", Kind: part.PropInt, Value: self.poleCountClamped(), Min: 1, Max: 8},
		},
	}
}

func (self *Relay) ApplyProp(action part.PropAction) bool {
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
		self.PoleCount = clampRelayPoleCount(value)
		return true
	default:
		return false
	}
}
