package editor

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

func HandleRelease(pt core.Pt, button int) {
	_, _ = pt, button
	Dragging = false
	BoxSelecting = false
}

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

func HandleScroll(delta float64) {
	_ = delta
}

func StartPlacement(typeID core.PartTypeID) {
	info, ok := part.Registry[typeID]
	if !ok {
		return
	}

	PlaceTool = typeID
	PlacePreview = info.New(world.AllocPartID(), core.Pt{})
	PlaceMode = true
}

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

func CopySelected() {
	Clipboard = nil
	for _, idx := range Selection {
		if idx >= 0 && idx < len(world.Parts) {
			Clipboard = append(Clipboard, world.Parts[idx])
		}
	}
}

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

func StartLabelEdit(partIdx int) {
	if partIdx < 0 || partIdx >= len(world.Parts) {
		return
	}
	LabelEditing = true
	LabelIndex = partIdx
	LabelBuffer = []rune(world.Parts[partIdx].Base().Label)
}

func CommitLabelEdit() {
	if !LabelEditing || LabelIndex < 0 || LabelIndex >= len(world.Parts) {
		return
	}
	pushUndo()
	world.Parts[LabelIndex].Base().Label = string(LabelBuffer)
	LabelEditing = false
}

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

func commitPlacement(pos core.Pt) {
	pushUndo()
	PlacePreview.Base().Pos = snapToGrid(pos)
	world.Parts = append(world.Parts, PlacePreview)
	PlacePreview = nil
	PlaceMode = false
}

func handleWireClick(pt core.Pt) {
	snapped := snapToGridOrPin(pt)
	WireDraft = append(WireDraft, snapped)

	if len(WireDraft) < 2 {
		return
	}

	from := WireDraft[len(WireDraft)-2]
	to := WireDraft[len(WireDraft)-1]
	pushUndo()
	world.Parts = append(world.Parts, wire.New(world.AllocPartID(), from, to, world.AllocPinID, world.AllocPinID))
}

func snapToGrid(pt core.Pt) core.Pt {
	const grid = 16.0
	snapped := core.Pt{
		X: math.Round(pt.X/grid) * grid,
		Y: math.Round(pt.Y/grid) * grid,
	}
	return core.LocalToWorld(core.BasePart{Pos: snapped}, core.Pt{})
}

func snapToGridOrPin(pt core.Pt) core.Pt {
	return snapToGrid(pt)
}

func partAt(pt core.Pt) int {
	_ = core.WorldToLocal(core.BasePart{}, pt)
	for i := len(world.Parts) - 1; i >= 0; i-- {
		if world.Parts[i].HitTest(pt).Hit {
			return i
		}
	}
	return -1
}
