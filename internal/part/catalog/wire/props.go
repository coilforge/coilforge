package wire

import "coilforge/internal/part"

// PropSpec handles prop spec.
func (self *Wire) PropSpec() part.PropSpec {
	return part.PropSpec{}
}

// ApplyProp handles apply prop.
func (self *Wire) ApplyProp(action part.PropAction) bool {
	_ = action
	return false
}
