package core

// Pt is a 2D point in world space.
type Pt struct {
	X float64
	Y float64
}

// Seg is a line segment between two points.
type Seg struct {
	A Pt
	B Pt
}

// Rect is an axis-aligned rectangle.
type Rect struct {
	Min Pt
	Max Pt
}

// Identity types.
type PartTypeID string
type PinID int
type NetID int

// BasePart is the common authored state embedded in every part.
type BasePart struct {
	ID       int        `json:"id"`
	TypeID   PartTypeID `json:"type"`
	Pos      Pt         `json:"pos"`
	Rotation int        `json:"rotation"`
	Mirror   bool       `json:"mirror"`
	Label    string     `json:"label"`
}

// PinAnchor is a pin's world position plus its ID.
type PinAnchor struct {
	Pt    Pt
	PinID PinID
}

// Net is a group of connected pins.
type Net struct {
	ID   int     `json:"id"`
	Pins []PinID `json:"pins"`
	Segs []Seg   `json:"segs"`
}

const (
	NetFloat = 0
	NetLow   = 1
	NetHigh  = 2
	NetShort = 3
)

func RectFromPoints(a, b Pt) Rect {
	r := Rect{Min: a, Max: b}
	return NormalizeRect(r)
}

func NormalizeRect(r Rect) Rect {
	if r.Min.X > r.Max.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Min.Y > r.Max.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r Rect) Intersects(other Rect) bool {
	return r.Min.X <= other.Max.X &&
		r.Max.X >= other.Min.X &&
		r.Min.Y <= other.Max.Y &&
		r.Max.Y >= other.Min.Y
}
