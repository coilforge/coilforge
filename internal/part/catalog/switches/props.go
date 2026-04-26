package switches

import "coilforge/internal/part"

func (self *Switches) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: self.Label},
			{Label: "Momentary", Kind: part.PropBool, Value: self.Mode == ModeMomentary},
		},
	}
}

func (self *Switches) ApplyProp(action part.PropAction) bool {
	switch action.Index {
	case 0:
		v, ok := action.NewValue.(string)
		if !ok {
			return false
		}
		self.Label = v
		return true
	case 1:
		v, ok := action.NewValue.(bool)
		if !ok {
			return false
		}
		if v {
			self.Mode = ModeMomentary
		} else {
			self.Mode = ModeToggle
		}
		return true
	default:
		return false
	}
}
