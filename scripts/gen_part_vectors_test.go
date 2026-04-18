//go:build genpartvectors

// Run (scripts package has multiple mains): go test -tags genpartvectors scripts/gen_part_vectors.go scripts/gen_part_vectors_test.go

package main

import (
	"math"
	"testing"
)

func TestSnapLineSegmentEndpoints_nearlyHorizontal(t *testing.T) {
	// Sloppy "horizontal" from real Boxy export: -1.088 to 0 over ~71 units
	x1, y1, x2, y2 := snapLineSegmentEndpoints(-62.94, -1.088, -134.098, 0)
	wantY := (-1.088 + 0) * 0.5
	if math.Abs(y1-wantY) > 1e-9 || math.Abs(y2-wantY) > 1e-9 {
		t.Fatalf("got (%v,%v)-(%v,%v) want y=%v both ends", x1, y1, x2, y2, wantY)
	}
	if math.Abs(y1-y2) > 1e-9 {
		t.Fatalf("y not aligned: %v vs %v", y1, y2)
	}
}

func TestSnapLineSegmentEndpoints_longShallowDiagonalUntouched(t *testing.T) {
	x1, y1, x2, y2 := snapLineSegmentEndpoints(0, 0, 200, 5)
	if x1 != 0 || y1 != 0 || x2 != 200 || y2 != 5 {
		t.Fatalf("long shallow diagonal was modified: (%v,%v)-(%v,%v)", x1, y1, x2, y2)
	}
}

func TestParseSVG_circle_fillFromStyle_ORANGE(t *testing.T) {
	const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="-1 -1 2 2"><circle style="fill: rgb(255, 172, 0);" r="1"/></svg>`
	a, err := parseSVG([]byte(svg))
	if err != nil {
		t.Fatal(err)
	}
	var fills []string
	for _, prim := range a.Prims {
		if prim.Kind == primCircle {
			fills = append(fills, prim.Fill)
		}
	}
	if len(fills) != 1 || fills[0] != "#FFAC00" {
		t.Fatalf("circle fill from style=rgb(255,172,0): got fills %v", fills)
	}
}

func TestNormalizePaintToHex_namedAndRGBA(t *testing.T) {
	if g := normalizePaintToHex("  gold  "); g != "#FFD700" {
		t.Fatalf("named color: got %q", g)
	}
	if g := normalizePaintToHex("rgba(10, 20, 30, 0.5)"); g != "#0A141E80" {
		t.Fatalf("rgba: got %q", g)
	}
	if g := normalizePaintToHex("rgb(100% 0% 0% / 0.25)"); g != "#FF000040" {
		t.Fatalf("rgb slash: got %q", g)
	}
}

func TestSnapLineSegmentEndpoints_verticalLineUnchanged(t *testing.T) {
	x1, y1, x2, y2 := snapLineSegmentEndpoints(0, -64, 0, 64)
	if x1 != 0 || x2 != 0 || y1 != -64 || y2 != 64 {
		t.Fatalf("pure vertical modified: (%v,%v)-(%v,%v)", x1, y1, x2, y2)
	}
}

func TestRotSVG90K_fourStepsIdentity(t *testing.T) {
	x, y := 12.3, -45.6
	tx, ty := x, y
	for k := 0; k < 4; k++ {
		tx, ty = rotSVG90K(tx, ty, 1)
	}
	if math.Abs(tx-x) > 1e-9 || math.Abs(ty-y) > 1e-9 {
		t.Fatalf("after 4 quarter turns want (%v,%v) got (%v,%v)", x, y, tx, ty)
	}
}
