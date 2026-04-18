package part

import (
	"coilforge/internal/core"
	"sort"
)

// VectorPinMarker binds a schematic pin ID to one red authoring marker index. Indices are
// 0-based into the slice from [RegisterPinLayout] — same document order as red
// circles in the merged SVG after build-time snapPinMarkersToGrid64 only.
type VectorPinMarker struct {
	PinID       core.PinID
	MarkerIndex int
}

// AnchorsFromVectorMarkers maps selected markers to pins. The generator emits every red
// dot from the SVG; parts may expose one pin, several, or all by listing the bindings
// they care about.
func AnchorsFromVectorMarkers(layoutName string, base core.BasePart, markers []VectorPinMarker) []core.PinAnchor {
	centers, _, ok := generatedVectorPinLayout(layoutName)
	if !ok || len(markers) == 0 {
		return nil
	}
	out := make([]core.PinAnchor, len(markers))
	for i := range markers {
		idx := markers[i].MarkerIndex
		if idx < 0 || idx >= len(centers) {
			return nil
		}
		w := SVGLocalToWorld(base, centers[idx].X, centers[idx].Y)
		out[i] = core.PinAnchor{Pt: w, PinID: markers[i].PinID}
	}
	return out
}

// AnchorsFromVectorPrefix binds pinIDs to markers 0, 1, … len(pinIDs)-1 in order (common
// case: use the first K markers when K equals the part’s exported pin count).
func AnchorsFromVectorPrefix(layoutName string, base core.BasePart, pinIDs []core.PinID) []core.PinAnchor {
	if len(pinIDs) == 0 {
		return nil
	}
	markers := make([]VectorPinMarker, len(pinIDs))
	for i := range pinIDs {
		markers[i] = VectorPinMarker{PinID: pinIDs[i], MarkerIndex: i}
	}
	return AnchorsFromVectorMarkers(layoutName, base, markers)
}

// AnchorsFromVectorMarkerIDs places pins using the SVG red circle id attributes as keys.
func AnchorsFromVectorMarkerIDs(layoutName string, base core.BasePart, pinByMarkerID map[string]core.PinID) []core.PinAnchor {
	centers, markerIDs, ok := generatedVectorPinLayout(layoutName)
	if !ok || len(pinByMarkerID) == 0 {
		return nil
	}
	indexByID := make(map[string]int)
	for i, id := range markerIDs {
		if id == "" {
			continue
		}
		if _, dup := indexByID[id]; dup {
			return nil
		}
		indexByID[id] = i
	}
	keys := make([]string, 0, len(pinByMarkerID))
	for k := range pinByMarkerID {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]core.PinAnchor, 0, len(keys))
	for _, k := range keys {
		idx, found := indexByID[k]
		if !found {
			return nil
		}
		pid := pinByMarkerID[k]
		w := SVGLocalToWorld(base, centers[idx].X, centers[idx].Y)
		out = append(out, core.PinAnchor{Pt: w, PinID: pid})
	}
	return out
}
