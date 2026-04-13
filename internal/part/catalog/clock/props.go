package clock

import "coilforge/internal/part"

func (c *Clock) PropSpec() part.PropSpec {
	return part.PropSpec{
		Items: []part.PropItem{
			{Label: "Label", Kind: part.PropText, Value: c.Label},
			{Label: "Period", Kind: part.PropInt, Value: c.PeriodTick, Min: 1, Max: 100000},
			{Label: "High", Kind: part.PropInt, Value: c.HighTick, Min: 1, Max: 100000},
		},
	}
}

func (c *Clock) ApplyProp(action part.PropAction) bool {
	switch action.Index {
	case 0:
		value, ok := action.NewValue.(string)
		if !ok {
			return false
		}
		c.Label = value
		return true
	case 1:
		value, ok := action.NewValue.(int)
		if !ok || value <= 0 {
			return false
		}
		c.PeriodTick = value
		return true
	case 2:
		value, ok := action.NewValue.(int)
		if !ok || value <= 0 {
			return false
		}
		c.HighTick = value
		return true
	default:
		return false
	}
}
