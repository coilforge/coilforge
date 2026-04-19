package editor

// File overview:
// editor implements interactive edit-mode behaviors for selection, placement, and transforms.
// Subsystem: editor interaction logic.
// It operates on world and part abstractions and is called by app input orchestration.
// Flow position: edit pipeline branch, independent from simulation execution.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/part/catalog/wire"
	"coilforge/internal/render"
	"coilforge/internal/world"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// wireToolID matches [wire.TypeID]; editor avoids importing catalog packages for routing helpers only.
const wireToolID core.PartTypeID = "wire"

// HandleMouseDown handles mouse down (edit mode): part pick, placement commit, or marquee start.
func HandleMouseDown(pt core.Pt, button int) {
	_ = button
	PressWorld = pt
	MouseDownOnEmpty = false
	PointerDownPart = -1
	Dragging = false
	DragMoved = false
	DragUndoRecorded = false
	BoxSelecting = false

	if PlaceMode && PlaceTool == wireToolID {
		handleWirePlaceClick(pt)
		return
	}

	if PlaceMode && PlacePreview != nil {
		commitPlacement(pt)
		return
	}

	idx := partAt(pt)
	PointerDownPart = idx
	MouseDownOnEmpty = idx < 0
	if idx >= 0 {
		if !selectionContains(idx) {
			Selection = []int{idx}
		}
		HoverIndex = idx
	}
}

func selectionContains(idx int) bool {
	for _, i := range Selection {
		if i == idx {
			return true
		}
	}
	return false
}

// HandleMouseUp handles mouse up: finalize marquee selection or apply empty-click clearing.
func HandleMouseUp(pt core.Pt, button int) {
	_ = button
	if ViewportPanDrag {
		endViewportPan()
		return
	}
	if BoxSelecting {
		sx0, _ := world.WorldToScreen(PressWorld)
		sx1, _ := world.WorldToScreen(pt)
		crossing := sx1 < sx0
		r := core.NormalizeRect(core.RectFromPoints(PressWorld, pt))
		Selection = collectBoxSelection(r, crossing)
		if len(Selection) == 0 {
			Selection = nil
			HoverIndex = -1
		} else {
			HoverIndex = Selection[len(Selection)-1]
		}
	} else if MouseDownOnEmpty {
		Selection = nil
		HoverIndex = -1
	}
	if !BoxSelecting && DragMoved {
		snapSelectedToMajorGrid()
	}
	Dragging = false
	DragMoved = false
	DragUndoRecorded = false
	BoxSelecting = false
	PointerDownPart = -1
	MouseDownOnEmpty = false
}

// HandleDrag handles drag: move selection, or update marquee when the press began on empty canvas.
func HandleDrag(pt core.Pt) {
	if PlaceMode && PlacePreview != nil {
		return
	}
	if PointerDownPart >= 0 {
		if !Dragging {
			Dragging = true
			DragStart = pt
			return
		}

		delta := core.Pt{X: pt.X - DragStart.X, Y: pt.Y - DragStart.Y}
		if delta.X != 0 || delta.Y != 0 {
			DragMoved = true
		}
		MoveSelected(delta)
		DragStart = pt
		return
	}
	if PointerDownPart < 0 {
		if !dragExceedsMarqueeThreshold(PressWorld, pt) {
			return
		}
		BoxSelecting = true
		BoxRect = core.NormalizeRect(core.RectFromPoints(PressWorld, pt))
		sx0, _ := world.WorldToScreen(PressWorld)
		sx1, _ := world.WorldToScreen(pt)
		BoxSelectCrossing = sx1 < sx0
	}
}

// HandleKey handles key.
func HandleKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyEscape:
		ClearTransient()
	case ebiten.KeyR:
		if PlaceMode && PlacePreview != nil {
			RotatePlacementPreview()
			return
		}
		RotateSelected()
	case ebiten.KeyM:
		MirrorSelected()
	case ebiten.KeyDelete, ebiten.KeyBackspace:
		DeleteSelected()
	}
}

// StartPlacement starts placement.
func StartPlacement(typeID core.PartTypeID) {
	info, ok := part.Registry[typeID]
	if !ok {
		return
	}

	PlaceTool = typeID
	if typeID == wireToolID {
		PlacePreview = nil
		WireAnchorSet = false
		WireAnchor = core.Pt{}
		WireHoverWorld = core.Pt{}
		PlaceMode = true
		return
	}
	if info.New == nil {
		return
	}
	PlacePreview = info.New(world.AllocPartID(), core.Pt{})
	PlaceMode = true
}

// UpdatePlacementPreview moves the ghost preview to the current pointer position.
func UpdatePlacementPreview(pt core.Pt) {
	if !PlaceMode {
		return
	}
	if PlaceTool == wireToolID {
		WireHoverWorld = snapToGrid(pt)
		return
	}
	if PlacePreview == nil {
		return
	}
	PlacePreview.Base().Pos = snapToGrid(pt)
}

// MoveSelected handles move selected.
func MoveSelected(delta core.Pt) {
	if len(Selection) == 0 {
		return
	}
	if delta.X != 0 || delta.Y != 0 {
		if !DragUndoRecorded {
			pushUndo()
			DragUndoRecorded = true
		}
	}
	for _, idx := range Selection {
		p := world.Parts[idx]
		if wo, ok := p.(part.WorldOffsettable); ok {
			wo.ApplyWorldOffset(delta)
			continue
		}
		base := p.Base()
		base.Pos.X += delta.X
		base.Pos.Y += delta.Y
	}
}

// RotateSelected rotates selected.
func RotateSelected() {
	if len(Selection) == 0 {
		return
	}
	if !selectionCanRotate() {
		return
	}
	if len(Selection) == 1 {
		idx := Selection[0]
		if idx < 0 || idx >= len(world.Parts) {
			return
		}
		pushUndo()
		base := world.Parts[idx].Base()
		slots := registryRotationSlots(world.Parts[idx])
		base.Rotation = rotateIndexBackward(slots, base.Rotation)
		return
	}

	u, ok := selectionUnionBounds()
	if !ok {
		return
	}
	pushUndo()
	g := rectCenter(u)
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			continue
		}
		p := world.Parts[idx]
		base := p.Base()
		slots := registryRotationSlots(p)
		steps := quarterTurnSlotSteps(slots)
		base.Rotation = rotateIndexBackwardN(slots, base.Rotation, steps)
		d := rotateVectorQuarterWorld(core.Pt{X: base.Pos.X - g.X, Y: base.Pos.Y - g.Y})
		base.Pos.X = g.X + d.X
		base.Pos.Y = g.Y + d.Y
	}
}

// RotatePlacementPreview rotates the ghost part clockwise (same convention as [RotateSelected]).
func RotatePlacementPreview() {
	if PlacePreview == nil {
		return
	}
	slots := registryRotationSlots(PlacePreview)
	if slots <= 0 {
		return
	}
	b := PlacePreview.Base()
	b.Rotation = rotateIndexBackward(slots, b.Rotation)
}

// rotateIndexBackward steps the rotation slot backward (CW through baked variants; +1 was CCW for our SVG basis).
func rotateIndexBackward(slots, idx int) int {
	if slots <= 0 {
		return idx
	}
	return ((idx-1)%slots + slots) % slots
}

func rotateIndexBackwardN(slots, idx, steps int) int {
	r := idx
	for i := 0; i < steps; i++ {
		r = rotateIndexBackward(slots, r)
	}
	return r
}

// quarterTurnSlotSteps is how many discrete baked-orientation steps equal one 90° world turn.
func quarterTurnSlotSteps(slots int) int {
	if slots <= 0 {
		return 1
	}
	st := int(math.Round(float64(slots) / 4.0))
	if st < 1 {
		st = 1
	}
	return st
}

// rotateVectorQuarterWorld rotates offset vector by 90° (same quarter as a 4-slot single-step rotate).
func rotateVectorQuarterWorld(v core.Pt) core.Pt {
	return core.Pt{X: v.Y, Y: -v.X}
}

func selectionUnionBounds() (core.Rect, bool) {
	var u core.Rect
	ok := false
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			continue
		}
		b := core.NormalizeRect(world.Parts[idx].Bounds())
		if !ok {
			u = b
			ok = true
			continue
		}
		u = unionRects(u, b)
	}
	return u, ok
}

func unionRects(a, b core.Rect) core.Rect {
	a = core.NormalizeRect(a)
	b = core.NormalizeRect(b)
	return core.NormalizeRect(core.Rect{
		Min: core.Pt{X: math.Min(a.Min.X, b.Min.X), Y: math.Min(a.Min.Y, b.Min.Y)},
		Max: core.Pt{X: math.Max(a.Max.X, b.Max.X), Y: math.Max(a.Max.Y, b.Max.Y)},
	})
}

func rectCenter(r core.Rect) core.Pt {
	r = core.NormalizeRect(r)
	return core.Pt{X: (r.Min.X + r.Max.X) * 0.5, Y: (r.Min.Y + r.Max.Y) * 0.5}
}

// registryRotationSlots returns 4 or 8 from the part type registry, else 0 (non-rotatable).
func registryRotationSlots(p part.Part) int {
	if p == nil {
		return 0
	}
	info, ok := part.Registry[p.Base().TypeID]
	if !ok {
		return 0
	}
	s := info.RotationSlots
	if !part.AllowsDiscreteRotation(s) {
		return 0
	}
	return s
}

func selectionCanRotate() bool {
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			return false
		}
		if registryRotationSlots(world.Parts[idx]) <= 0 {
			return false
		}
	}
	return true
}

// MirrorSelected mirrors selected.
func MirrorSelected() {
	if len(Selection) == 0 {
		return
	}
	pushUndo()
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			continue
		}
		if _, ok := world.Parts[idx].(part.WorldOffsettable); ok {
			continue
		}
		base := world.Parts[idx].Base()
		base.Mirror = !base.Mirror
	}
}

// DeleteSelected deletes selected.
func DeleteSelected() {
	if len(Selection) == 0 {
		return
	}

	pushUndo()
	sort.Sort(sort.Reverse(sort.IntSlice(Selection)))
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			continue
		}
		world.Parts = append(world.Parts[:idx], world.Parts[idx+1:]...)
	}
	Selection = nil
}

// CopySelected copies selected.
func CopySelected() {
	Clipboard = nil
	for _, idx := range Selection {
		if idx >= 0 && idx < len(world.Parts) {
			Clipboard = append(Clipboard, world.Parts[idx])
		}
	}
}

// Paste pastes its work.
func Paste(offset core.Pt) {
	if len(Clipboard) == 0 {
		return
	}

	pushUndo()
	start := len(world.Parts)
	for _, orig := range Clipboard {
		cloned := orig.Clone(world.AllocPartID(), world.AllocPinID)
		if wo, ok := cloned.(part.WorldOffsettable); ok {
			wo.ApplyWorldOffset(offset)
		} else {
			b := cloned.Base()
			b.Pos.X += offset.X
			b.Pos.Y += offset.Y
		}
		world.Parts = append(world.Parts, cloned)
	}

	Selection = Selection[:0]
	for idx := start; idx < len(world.Parts); idx++ {
		Selection = append(Selection, idx)
	}
}

// StartLabelEdit starts label edit.
func StartLabelEdit(partIdx int) {
	if partIdx < 0 || partIdx >= len(world.Parts) {
		return
	}
	LabelEditing = true
	LabelIndex = partIdx
	LabelBuffer = []rune(world.Parts[partIdx].Base().Label)
}

// CommitLabelEdit commits label edit.
func CommitLabelEdit() {
	if !LabelEditing || LabelIndex < 0 || LabelIndex >= len(world.Parts) {
		return
	}
	pushUndo()
	world.Parts[LabelIndex].Base().Label = string(LabelBuffer)
	LabelEditing = false
}

// DrawOverlays draws overlays.
func DrawOverlays(dst *ebiten.Image) {
	for _, idx := range Selection {
		if idx >= 0 && idx < len(world.Parts) {
			render.DrawSelectionOutline(dst, world.Parts[idx].Bounds())
		}
	}

	if PlacePreview != nil {
		PlacePreview.Draw(part.DrawContext{
			Dst:      dst,
			Cam:      world.Cam,
			Zoom:     world.Zoom,
			ScreenW:  world.ScreenW,
			ScreenH:  world.ScreenH,
			DarkMode: render.DarkMode,
			Ghost:    true,
		})
	}

	if BoxSelecting {
		render.DrawBoxSelect(dst, BoxRect, BoxSelectCrossing)
	}

	if PlaceMode && PlaceTool == wireToolID && WireAnchorSet {
		pts := orthRoutePreview(WireAnchor, WireHoverWorld)
		if len(pts) >= 2 {
			wire.DrawOrthogonalPolyline(part.DrawContext{
				Dst:      dst,
				Cam:      world.Cam,
				Zoom:     world.Zoom,
				ScreenW:  world.ScreenW,
				ScreenH:  world.ScreenH,
				DarkMode: render.DarkMode,
				Ghost:    true,
			}, pts, render.GhostTint())
		}
	}
}

// commitPlacement handles commit placement.
func commitPlacement(pos core.Pt) {
	pushUndo()
	// New() does not allocate pin IDs; Clone does. Re-clone the preview with the same part id
	// so every placed part gets unique [core.PinID] values (required for net map and save).
	placed := PlacePreview.Clone(PlacePreview.Base().ID, world.AllocPinID)
	placed.Base().Pos = snapToGrid(pos)
	world.Parts = append(world.Parts, placed)
	PlacePreview = nil
	PlaceMode = false
}

// snapToGrid snaps to the major grid (part placement / pin pitch).
func snapToGrid(pt core.Pt) core.Pt {
	grid := world.MajorGridWorld
	snapped := core.Pt{
		X: math.Round(pt.X/grid) * grid,
		Y: math.Round(pt.Y/grid) * grid,
	}
	return core.LocalToWorld(core.BasePart{Pos: snapped}, core.Pt{})
}

// snapSelectedToMajorGrid snaps placed part origins after a drag drop (same grid as placement).
func snapSelectedToMajorGrid() {
	if len(Selection) == 0 {
		return
	}
	type move struct {
		idx int
		pos core.Pt
	}
	var moves []move
	for _, idx := range Selection {
		if idx < 0 || idx >= len(world.Parts) {
			continue
		}
		p := world.Parts[idx]
		if _, ok := p.(part.WorldOffsettable); ok {
			continue
		}
		base := p.Base()
		next := snapToGrid(base.Pos)
		if next.X != base.Pos.X || next.Y != base.Pos.Y {
			moves = append(moves, move{idx: idx, pos: next})
		}
	}
	if len(moves) == 0 {
		return
	}
	for _, m := range moves {
		world.Parts[m.idx].Base().Pos = m.pos
	}
}

// partAt handles part at.
func partAt(pt core.Pt) int {
	for i := len(world.Parts) - 1; i >= 0; i-- {
		if world.Parts[i].HitTest(pt).Hit {
			return i
		}
	}
	return -1
}

const marqueeThresholdPx = 4.0

func dragExceedsMarqueeThreshold(a, b core.Pt) bool {
	sx0, sy0 := world.WorldToScreen(a)
	sx1, sy1 := world.WorldToScreen(b)
	dx := sx1 - sx0
	dy := sy1 - sy0
	return dx*dx+dy*dy >= marqueeThresholdPx*marqueeThresholdPx
}

func collectBoxSelection(r core.Rect, crossing bool) []int {
	var out []int
	for i := range world.Parts {
		b := core.NormalizeRect(world.Parts[i].Bounds())
		ok := (crossing && r.Intersects(b)) || (!crossing && r.ContainsRect(b))
		if ok {
			out = append(out, i)
		}
	}
	return out
}

func handleWirePlaceClick(pt core.Pt) {
	a := snapToGrid(pt)
	if !WireAnchorSet {
		WireAnchor = a
		WireAnchorSet = true
		WireHoverWorld = a
		return
	}
	b := snapToGrid(pt)
	if len(orthRoutePreview(WireAnchor, b)) < 2 {
		return
	}
	info, ok := part.Registry[wireToolID]
	if !ok || info.NewWire == nil {
		return
	}
	w := info.NewWire(world.AllocPartID(), WireAnchor, b, world.AllocPinID)
	if w == nil {
		return
	}
	pushUndo()
	world.Parts = append(world.Parts, w)
	WireAnchorSet = false
}

// orthRoutePreview matches [wire.OrthogonalRoute] (kept here to avoid editor → catalog import).
func orthRoutePreview(a, b core.Pt) []core.Pt {
	if a.X == b.X && a.Y == b.Y {
		return nil
	}
	if a.X == b.X || a.Y == b.Y {
		return []core.Pt{a, b}
	}
	return []core.Pt{a, {X: b.X, Y: a.Y}, b}
}
