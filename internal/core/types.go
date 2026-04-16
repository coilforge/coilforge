package core

// File overview:
// types defines shared schematic primitives, IDs, and base part data structures.
// Subsystem: core leaf types.
// It is consumed by part, world, editor, sim, render, and app without importing internals.
// Flow position: lowest-level data contract for the full architecture.

// Pt is a 2D point in world space.
type Pt struct {
	X float64 // x coordinate.
	Y float64 // y coordinate.
}

// Seg is a line segment between two points.
type Seg struct {
	A Pt // a value.
	B Pt // b value.
}

// Rect is an axis-aligned rectangle.
type Rect struct {
	Min Pt // min value.
	Max Pt // max value.
}

// Identity types.
type PartTypeID string
type PinID int
type NetID int

// BasePart is the common authored state embedded in every part.
type BasePart struct {
	ID       int        `json:"id"`       // identifier.
	TypeID   PartTypeID `json:"type"`     // type id value.
	Pos      Pt         `json:"pos"`      // world position.
	Rotation int        `json:"rotation"` // rotation value.
	Mirror   bool       `json:"mirror"`   // mirror value.
	Label    string     `json:"label"`    // label text.
}

// PinAnchor is a pin's world position plus its ID.
type PinAnchor struct {
	Pt    Pt    // pt value.
	PinID PinID // pin id value.
}

// Net is a group of connected pins.
type Net struct {
	ID   int     `json:"id"`   // identifier.
	Pins []PinID `json:"pins"` // pin list.
	Segs []Seg   `json:"segs"` // segs value.
}

const (
	NetFloat = 0 // NetFloat marks an undriven net.
	NetLow   = 1 // NetLow marks a driven-low net.
	NetHigh  = 2 // NetHigh marks a driven-high net.
	NetShort = 3 // NetShort marks a conflicting net.
)

// RectFromPoints handles rect from points.
func RectFromPoints(a, b Pt) Rect {
	r := Rect{Min: a, Max: b}
	return NormalizeRect(r)
}

// NormalizeRect handles normalize rect.
func NormalizeRect(r Rect) Rect {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// Intersects handles intersects.
func (r Rect) Intersects(other Rect) bool {
	return r.Min.X <= other.Max.X &&
		r.Max.X >= other.Min.X &&
		r.Min.Y <= other.Max.Y &&
		r.Max.Y >= other.Min.Y
}
