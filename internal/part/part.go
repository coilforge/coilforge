package part

import (
	"coilforge/internal/core"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Part is the contract every component implements.
type Part interface {
	Base() *core.BasePart
	Bounds() core.Rect
	Anchors() []core.PinAnchor
	HitTest(pt core.Pt) HitResult
	Segments() []core.Seg
	Draw(ctx DrawContext)
	Clone(newID int, allocPin func() core.PinID) Part
	PropSpec() PropSpec
	ApplyProp(action PropAction) bool
	MarshalJSON() ([]byte, error)
}

type DrawContext struct {
	Dst      *ebiten.Image
	Cam      core.Pt
	Zoom     float64
	ScreenW  int
	ScreenH  int
	Ghost    bool
	Selected bool
	NetState func(core.PinID) int
}

func (ctx DrawContext) WorldToScreen(pt core.Pt) (float64, float64) {
	return (pt.X-ctx.Cam.X)*ctx.Zoom + float64(ctx.ScreenW)/2,
		(pt.Y-ctx.Cam.Y)*ctx.Zoom + float64(ctx.ScreenH)/2
}

type HitResult struct {
	Hit   bool
	Kind  int
	PinID core.PinID
}

const (
	HitBody = iota
	HitLabel
	HitPin
)

type SimPart interface {
	Tick(ctx SimContext) bool
}

type SimContext struct {
	NetByPin     func(core.PinID) int
	NetState     func(int) int
	Tick         uint64
	TickMicros   int
	EnableJitter bool
	Rand         *rand.Rand
}

type NetSeeder interface {
	SeedNets(netByPin func(core.PinID) int, high, low map[int]bool)
}

type Conductor interface {
	AddConductive(union NetUnion, netByPin func(core.PinID) int)
}

type StateEdger interface {
	AddStateEdges(netByPin func(core.PinID) int, graph *StateGraph)
}

type InputHandler interface {
	HandleInput(active bool) (changed, momentary bool)
	ReleaseMomentary() bool
}

type NetUnion interface {
	Union(a, b int)
	Find(a int) int
}

type StateGraph struct {
	Edges []StateEdge
}

type StateEdge struct {
	FromNet int
	ToNet   int
	Drive   int
}

type VectorAsset struct {
	Name string
}

func (asset VectorAsset) Draw(ctx DrawContext, bounds core.Rect) {
	_, _ = ctx, bounds
}
