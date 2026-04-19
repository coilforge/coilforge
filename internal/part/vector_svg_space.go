package part

import "coilforge/internal/core"

// SVGUserUnitToWorld is how many schematic world units one SVG user unit covers.
// A 512-unit-wide path is 512*this world units long; a 64-unit shape is 1/8 of that.
// Adjust this single constant to set global schematic scale for vector parts.
// 1/8 makes symbols and stroke widths twice as large in world space vs the earlier 1/16 mapping.
const SVGUserUnitToWorld = 1.0 / 8.0

// SVGLocalToWorld maps symbol-centred SVG user coordinates through [core.LocalToWorld]
// (position, rotation, mirror).
func SVGLocalToWorld(base core.BasePart, svgX, svgY float64) core.Pt {
	return core.LocalToWorld(base, core.Pt{X: svgX * SVGUserUnitToWorld, Y: svgY * SVGUserUnitToWorld})
}

// SVGPointToWorld maps SVG user space to world with identity rotation/mirror (pos only).
func SVGPointToWorld(pos core.Pt, x, y float64) core.Pt {
	return SVGLocalToWorld(core.BasePart{Pos: pos}, x, y)
}

// HitBoundsFromSVGExtents maps an axis-aligned box in symbol-centred SVG user units to a world AABB
// after rotation/mirror (uses the four corners).
func HitBoundsFromSVGExtents(base core.BasePart, minX, minY, maxX, maxY float64) core.Rect {
	corners := []core.Pt{
		SVGLocalToWorld(base, minX, minY),
		SVGLocalToWorld(base, maxX, minY),
		SVGLocalToWorld(base, minX, maxY),
		SVGLocalToWorld(base, maxX, maxY),
	}
	r := core.RectFromPoints(corners[0], corners[1])
	for _, p := range corners[2:] {
		if p.X < r.Min.X {
			r.Min.X = p.X
		}
		if p.Y < r.Min.Y {
			r.Min.Y = p.Y
		}
		if p.X > r.Max.X {
			r.Max.X = p.X
		}
		if p.Y > r.Max.Y {
			r.Max.Y = p.Y
		}
	}
	return r
}
