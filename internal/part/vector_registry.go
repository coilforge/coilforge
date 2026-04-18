package part

import "coilforge/internal/core"

// VectorDrawFunc draws a named [VectorAsset] with symbol-centred SVG mapping in [DrawVGLine] and friends.
type VectorDrawFunc func(DrawContext, core.BasePart) bool

type vectorPinLayout struct {
	centers   []core.Pt
	markerIDs []string // parallel to centers; empty string when the SVG circle had no id=""
}

// vectorHitBounds is a tight axis-aligned box in symbol-centred SVG user units (from generated bounds of drawn geometry).
type vectorHitBounds struct {
	minX, minY, maxX, maxY float64
}

var (
	vectorDrawByName  = make(map[string]VectorDrawFunc)
	vectorPinByName   = make(map[string]vectorPinLayout)
	vectorHitByName   = make(map[string]vectorHitBounds)
)

// RegisterVectorDraw records build-time vector art for the given [VectorAsset].Name.
// Each name is registered at most once; duplicate registration panics.
func RegisterVectorDraw(name string, draw VectorDrawFunc) {
	if name == "" || draw == nil {
		panic("part.RegisterVectorDraw: invalid args")
	}
	if _, ok := vectorDrawByName[name]; ok {
		panic("part.RegisterVectorDraw: duplicate " + name)
	}
	vectorDrawByName[name] = draw
}

// RegisterPinLayout records red pin marker centers in symbol-centred SVG user units for [AnchorsFromVectorPrefix]
// and [AnchorsFromVectorMarkerIDs].
func RegisterPinLayout(name string, centers []core.Pt, markerIDs []string) {
	if name == "" {
		panic("part.RegisterPinLayout: empty name")
	}
	if _, ok := vectorPinByName[name]; ok {
		panic("part.RegisterPinLayout: duplicate " + name)
	}
	switch {
	case len(markerIDs) == 0:
		markerIDs = make([]string, len(centers))
	case len(markerIDs) != len(centers):
		panic("part.RegisterPinLayout: markerIDs length mismatch")
	}
	seenNonEmpty := map[string]bool{}
	for _, id := range markerIDs {
		if id == "" {
			continue
		}
		if seenNonEmpty[id] {
			panic("part.RegisterPinLayout: duplicate marker id " + id + " for " + name)
		}
		seenNonEmpty[id] = true
	}
	c := make([]core.Pt, len(centers))
	copy(c, centers)
	ids := make([]string, len(markerIDs))
	copy(ids, markerIDs)
	vectorPinByName[name] = vectorPinLayout{centers: c, markerIDs: ids}
}

// RegisterVectorHitBounds records body hit-testing extents in symbol-centred SVG user units (from codegen geometry bounds).
func RegisterVectorHitBounds(name string, minX, minY, maxX, maxY float64) {
	if name == "" {
		panic("part.RegisterVectorHitBounds: empty name")
	}
	if _, ok := vectorHitByName[name]; ok {
		panic("part.RegisterVectorHitBounds: duplicate " + name)
	}
	if maxX < minX || maxY < minY {
		panic("part.RegisterVectorHitBounds: invalid box for " + name)
	}
	vectorHitByName[name] = vectorHitBounds{minX: minX, minY: minY, maxX: maxX, maxY: maxY}
}

// HitBoundsFromVectorLayout returns selection bounds for a registered vector name (generated geometry bounds).
func HitBoundsFromVectorLayout(layoutName string, base core.BasePart) (core.Rect, bool) {
	h, ok := vectorHitByName[layoutName]
	if !ok {
		return core.Rect{}, false
	}
	return HitBoundsFromSVGExtents(base, h.minX, h.minY, h.maxX, h.maxY), true
}

func drawGeneratedVectorAsset(name string, ctx DrawContext, base core.BasePart) bool {
	if fn, ok := vectorDrawByName[name]; ok {
		return fn(ctx, base)
	}
	return drawManualVectorAsset(name, ctx, base)
}

func generatedVectorPinLayout(name string) (centers []core.Pt, markerIDs []string, ok bool) {
	lay, ok := vectorPinByName[name]
	if !ok {
		return nil, nil, false
	}
	return lay.centers, lay.markerIDs, true
}
