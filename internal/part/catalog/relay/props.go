package relay

import "coilforge/internal/part"

func (self *Relay) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: self.Label},
		},
	}
}

func (self *Relay) ApplyProp(action part.PropAction) bool {
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
