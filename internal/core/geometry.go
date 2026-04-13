package core

import "math"

func LocalToWorld(base BasePart, local Pt) Pt {
	x := local.X
	y := local.Y
	if base.Mirror {
		x = -x
	}

	switch ((base.Rotation % 4) + 4) % 4 {
	case 1:
		x, y = -y, x
	case 2:
		x, y = -x, -y
	case 3:
		x, y = y, -x
	}

	return Pt{X: base.Pos.X + x, Y: base.Pos.Y + y}
}

func WorldToLocal(base BasePart, world Pt) Pt {
	x := world.X - base.Pos.X
	y := world.Y - base.Pos.Y

	switch ((base.Rotation % 4) + 4) % 4 {
	case 1:
		x, y = y, -x
	case 2:
		x, y = -x, -y
	case 3:
		x, y = -y, x
	}

	if base.Mirror {
		x = -x
	}

	return Pt{X: x, Y: y}
}

func PointInRect(pt Pt, r Rect) bool {
	r = NormalizeRect(r)
	return pt.X >= r.Min.X && pt.X <= r.Max.X && pt.Y >= r.Min.Y && pt.Y <= r.Max.Y
}

func PointNearSeg(pt Pt, seg Seg, tolerance float64) bool {
	dx := seg.B.X - seg.A.X
	dy := seg.B.Y - seg.A.Y
	if dx == 0 && dy == 0 {
		return math.Hypot(pt.X-seg.A.X, pt.Y-seg.A.Y) <= tolerance
	}

	t := ((pt.X-seg.A.X)*dx + (pt.Y-seg.A.Y)*dy) / (dx*dx + dy*dy)
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	nearest := Pt{
		X: seg.A.X + t*dx,
		Y: seg.A.Y + t*dy,
	}
	return math.Hypot(pt.X-nearest.X, pt.Y-nearest.Y) <= tolerance
}

func RotateRect(r Rect, rotation int, mirror bool) Rect {
	corners := []Pt{
		r.Min,
		{X: r.Min.X, Y: r.Max.Y},
		{X: r.Max.X, Y: r.Min.Y},
		r.Max,
	}

	base := BasePart{Rotation: rotation, Mirror: mirror}
	out := Rect{
		Min: Pt{X: math.MaxFloat64, Y: math.MaxFloat64},
		Max: Pt{X: -math.MaxFloat64, Y: -math.MaxFloat64},
	}

	for _, corner := range corners {
		p := LocalToWorld(base, corner)
		if p.X < out.Min.X {
			out.Min.X = p.X
		}
		if p.Y < out.Min.Y {
			out.Min.Y = p.Y
		}
		if p.X > out.Max.X {
			out.Max.X = p.X
		}
		if p.Y > out.Max.Y {
			out.Max.Y = p.Y
		}
	}

	return out
}
