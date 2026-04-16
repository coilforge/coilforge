package editor

// File overview:
// editor implements interactive edit-mode behaviors for selection, placement, and transforms.
// Subsystem: editor interaction logic.
// It operates on world and part abstractions and is called by app input orchestration.
// Flow position: edit pipeline branch, independent from simulation execution.

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/render"
	"coilforge/internal/world"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// HandleClick handles click.
func HandleClick(pt core.Pt, button int) {
	_ = button

	if PlaceMode && PlacePreview != nil {
		commitPlacement(pt)
		return
	}

	if WireMode {
		handleWireClick(pt)
		return
	}

	if idx := partAt(pt); idx >= 0 {
		Selection = []int{idx}
		HoverIndex = idx
		return
	}

	Selection = nil
	HoverIndex = -1
}

// HandleRelease handles release.
func HandleRelease(pt core.Pt, button int) {
	_, _ = pt, button
	Dragging = false
	BoxSelecting = false
}

// HandleDrag handles drag.
func HandleDrag(pt core.Pt) {
	if !Dragging {
		Dragging = true
		DragStart = pt
		return
	}

	delta := core.Pt{X: pt.X - DragStart.X, Y: pt.Y - DragStart.Y}
	MoveSelected(delta)
	DragStart = pt
}

// HandleKey handles key.
func HandleKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyEscape:
		ClearTransient()
	case ebiten.KeyR:
		RotateSelected()
	case ebiten.KeyM:
		MirrorSelected()
	case ebiten.KeyDelete, ebiten.KeyBackspace:
		DeleteSelected()
	case ebiten.KeyW:
		WireMode = !WireMode
		if !WireMode {
			WireDraft = nil
		}
	}
}

// HandleScroll handles scroll.
func HandleScroll(delta float64) {
	_ = delta
}

// StartPlacement starts placement.
func StartPlacement(typeID core.PartTypeID) {
	info, ok := part.Registry[typeID]
	if !ok {
		return
	}

	PlaceTool = typeID
	PlacePreview = info.New(world.AllocPartID(), core.Pt{})
	PlaceMode = true
}

// MoveSelected handles move selected.
func MoveSelected(delta core.Pt) {
	if len(Selection) == 0 {
		return
	}
	pushUndo()
	for _, idx := range Selection {
		base := world.Parts[idx].Base()
		base.Pos.X += delta.X
		base.Pos.Y += delta.Y
	}
}

// RotateSelected rotates selected.
func RotateSelected() {
	if len(Selection) == 0 {
		return
	}
	pushUndo()
	for _, idx := range Selection {
		base := world.Parts[idx].Base()
		base.Rotation = (base.Rotation + 1) % 4
	}
}

// MirrorSelected mirrors selected.
func MirrorSelected() {
	if len(Selection) == 0 {
		return
	}
	pushUndo()
	for _, idx := range Selection {
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
		cloned.Base().Pos.X += offset.X
		cloned.Base().Pos.Y += offset.Y
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
			Dst:     dst,
			Cam:     world.Cam,
			Zoom:    world.Zoom,
			ScreenW: world.ScreenW,
			ScreenH: world.ScreenH,
			Ghost:   true,
		})
	}

	if WireMode && len(WireDraft) > 0 {
		render.DrawWireDraft(dst, WireDraft)
	}

	if BoxSelecting {
		render.DrawBoxSelect(dst, BoxRect)
	}
}

// commitPlacement handles commit placement.
func commitPlacement(pos core.Pt) {
	pushUndo()
	PlacePreview.Base().Pos = snapToGrid(pos)
	world.Parts = append(world.Parts, PlacePreview)
	PlacePreview = nil
	PlaceMode = false
}

// handleWireClick handles handle wire click.
func handleWireClick(pt core.Pt) {
	snapped := snapToGridOrPin(pt)
	WireDraft = append(WireDraft, snapped)

	if len(WireDraft) < 2 {
		return
	}

	from := WireDraft[len(WireDraft)-2]
	to := WireDraft[len(WireDraft)-1]
	info, ok := part.Registry[core.PartTypeID("wire")]
	if !ok || info.NewWire == nil {
		return
	}
	pushUndo()
	world.Parts = append(world.Parts, info.NewWire(world.AllocPartID(), from, to, world.AllocPinID))
}

// snapToGrid handles snap to grid.
func snapToGrid(pt core.Pt) core.Pt {
	const grid = 16.0
	snapped := core.Pt{
		X: math.Round(pt.X/grid) * grid,
		Y: math.Round(pt.Y/grid) * grid,
	}
	return core.LocalToWorld(core.BasePart{Pos: snapped}, core.Pt{})
}

// snapToGridOrPin handles snap to grid or pin.
func snapToGridOrPin(pt core.Pt) core.Pt {
	return snapToGrid(pt)
}

// partAt handles part at.
func partAt(pt core.Pt) int {
	_ = core.WorldToLocal(core.BasePart{}, pt)
	for i := len(world.Parts) - 1; i >= 0; i-- {
		if world.Parts[i].HitTest(pt).Hit {
			return i
		}
	}
	return -1
}
