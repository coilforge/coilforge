package part

import (
	"coilforge/internal/core"
	"encoding/json"

	"github.com/hajimehoshi/ebiten/v2"
)

type TypeInfo struct {
	New     func(id int, pos core.Pt) Part
	NewWire func(id int, from, to core.Pt, allocPin func() core.PinID) Part
	Decode  func(data json.RawMessage) (Part, error)
	Tools   []string
	Icon    func() *ebiten.Image
}

var Registry = map[core.PartTypeID]TypeInfo{}

func Register(id core.PartTypeID, info TypeInfo) {
	Registry[id] = info
}
