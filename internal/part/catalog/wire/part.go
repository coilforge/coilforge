package wire

// File overview:
// part defines the wire polyline part type (orthogonal routing, two endpoint pins).
// Subsystem: part catalog (wire).

import (
	"encoding/json"
	"fmt"
	"math"

	"coilforge/internal/core"
	"coilforge/internal/part"
)

// TypeID defines a package-level constant.
const TypeID core.PartTypeID = "wire"

// Wire is a conducting polyline in world space (minor-grid snapped). Pins sit on the first and last points.
type Wire struct {
	core.BasePart
	Points []core.Pt   `json:"points"`
	PinA   core.PinID  `json:"pinA"`
	PinB   core.PinID  `json:"pinB"`
}

func init() {
	part.Register(TypeID, part.TypeInfo{
		NewWire:       newWireEndpoints,
		Decode:        decodePart,
		Label:         "Wire",
		Tools:         []string{"main"},
		Icon:          toolbarIcon,
		RotationSlots: 0,
	})
}

func newWireEndpoints(id int, from, to core.Pt, allocPin func() core.PinID) part.Part {
	pts := OrthogonalRoute(from, to)
	if len(pts) < 2 {
		return nil
	}
	// Editor places multi-leg routes as multiple Wire parts (one straight segment each); callers that
	// still pass a diagonal pair get one polyline with an elbow here.
	return newWirePolyline(id, pts, allocPin)
}

// NewStraightWire creates one axis-aligned schematic segment using exactly two waypoints — no OrthogonalRoute.
//
// Routing each L-leg through [OrthogonalRoute] uses float equality on X/Y; long grid coordinates can differ
// in the least bits so a single leg is misclassified as diagonal and collapses into one Wire with three points.
// Splitting routes in the editor must call this instead of registry NewWire per leg.
func NewStraightWire(id int, from, to core.Pt, allocPin func() core.PinID) part.Part {
	if from.X == to.X && from.Y == to.Y {
		return nil
	}
	return newWirePolyline(id, []core.Pt{from, to}, allocPin)
}

// NewPolylineWire builds a wire through an explicit vertex path (≥2 snapped points).
// Used when splitting an existing polyline into two nets at a tee/junction on a segment interior.
func NewPolylineWire(id int, pts []core.Pt, allocPin func() core.PinID) part.Part {
	if len(pts) < 2 {
		return nil
	}
	return newWirePolyline(id, pts, allocPin)
}

func newWirePolyline(id int, pts []core.Pt, allocPin func() core.PinID) *Wire {
	w := &Wire{
		BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pts[0]},
		Points:   append([]core.Pt(nil), pts...),
		PinA:     allocPin(),
		PinB:     allocPin(),
	}
	return w
}

func decodePart(data json.RawMessage) (part.Part, error) {
	var w Wire
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, err
	}
	if w.TypeID == "" {
		w.TypeID = TypeID
	}
	if len(w.Points) < 2 {
		return nil, fmt.Errorf("wire: need at least two points")
	}
	w.BasePart.Pos = w.Points[0]
	return &w, nil
}

// Base handles base.
func (self *Wire) Base() *core.BasePart {
	return &self.BasePart
}

// Segments handles segments.
func (self *Wire) Segments() []core.Seg {
	out := make([]core.Seg, 0, len(self.Points)-1)
	for i := 0; i < len(self.Points)-1; i++ {
		out = append(out, core.Seg{A: self.Points[i], B: self.Points[i+1]})
	}
	return out
}

// Anchors handles anchors.
func (self *Wire) Anchors() []core.PinAnchor {
	if len(self.Points) < 2 {
		return nil
	}
	return []core.PinAnchor{
		{Pt: self.Points[0], PinID: self.PinA},
		{Pt: self.Points[len(self.Points)-1], PinID: self.PinB},
	}
}

// Bounds handles bounds.
func (self *Wire) Bounds() core.Rect {
	if len(self.Points) == 0 {
		return core.Rect{}
	}
	minX, maxX := self.Points[0].X, self.Points[0].X
	minY, maxY := self.Points[0].Y, self.Points[0].Y
	for _, p := range self.Points[1:] {
		minX = math.Min(minX, p.X)
		maxX = math.Max(maxX, p.X)
		minY = math.Min(minY, p.Y)
		maxY = math.Max(maxY, p.Y)
	}
	pad := 2.0
	return core.Rect{
		Min: core.Pt{X: minX - pad, Y: minY - pad},
		Max: core.Pt{X: maxX + pad, Y: maxY + pad},
	}
}

// Clone handles clone.
func (self *Wire) Clone(newID int, allocPin func() core.PinID) part.Part {
	c := *self
	c.ID = newID
	c.PinA = allocPin()
	c.PinB = allocPin()
	c.Points = append([]core.Pt(nil), self.Points...)
	c.BasePart.Pos = c.Points[0]
	return &c
}

// MarshalJSON handles marshal json.
func (self *Wire) MarshalJSON() ([]byte, error) {
	type partJSON Wire
	self.BasePart.Pos = self.Points[0]
	return json.Marshal((*partJSON)(self))
}

// ApplyWorldOffset shifts every waypoint (drag/paste).
func (self *Wire) ApplyWorldOffset(delta core.Pt) {
	for i := range self.Points {
		self.Points[i].X += delta.X
		self.Points[i].Y += delta.Y
	}
	self.BasePart.Pos = self.Points[0]
}

const snapEps = 1e-6

// SnapWaypointsToMajorGrid rounds every vertex to the schematic major grid and drops duplicate consecutive
// vertices (can happen after a free drag that was not aligned to the grid).
func (self *Wire) SnapWaypointsToMajorGrid(grid float64) {
	if grid <= 0 || len(self.Points) < 2 {
		return
	}
	tmp := append([]core.Pt(nil), self.Points...)
	for i := range tmp {
		tmp[i].X = math.Round(tmp[i].X/grid) * grid
		tmp[i].Y = math.Round(tmp[i].Y/grid) * grid
	}
	var deduped []core.Pt
	for _, q := range tmp {
		if len(deduped) > 0 && ptNearEq(deduped[len(deduped)-1], q) {
			continue
		}
		deduped = append(deduped, q)
	}
	if len(deduped) < 2 {
		return
	}
	self.Points = deduped
	self.BasePart.Pos = self.Points[0]
}

func ptNearEq(a, b core.Pt) bool {
	return math.Abs(a.X-b.X) < snapEps && math.Abs(a.Y-b.Y) < snapEps
}
