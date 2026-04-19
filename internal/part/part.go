package part

// File overview:
// part defines the core part interfaces used uniformly across editor, sim, and render.
// Subsystem: part contracts.
// It depends only on core and is implemented by each catalog part package.
// Flow position: abstraction boundary between generic tooling and concrete parts.

import (
	"coilforge/internal/core"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Part is the contract every catalog part type implements.
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
	Dst      *ebiten.Image        // dst value.
	Cam      core.Pt              // camera position.
	Zoom     float64              // zoom factor.
	ScreenW  int                  // screen width.
	ScreenH  int                  // screen height.
	DarkMode bool                 // schematic theme: vector ink maps to light strokes when true.
	Ghost    bool                 // ghost value.
	Selected bool                 // selected value.
	NetState func(core.PinID) int // net state value.
}

// WorldToScreen handles world to screen.
func (ctx DrawContext) WorldToScreen(pt core.Pt) (float64, float64) {
	return (pt.X-ctx.Cam.X)*ctx.Zoom + float64(ctx.ScreenW)/2,
		(pt.Y-ctx.Cam.Y)*ctx.Zoom + float64(ctx.ScreenH)/2
}

type HitResult struct {
	Hit   bool       // hit value.
	Kind  int        // kind value.
	PinID core.PinID // pin id value.
}

const (
	HitBody  = iota // HitBody marks a body hit target.
	HitLabel        // HitLabel marks a label hit target.
	HitPin          // HitPin marks a pin hit target.
)

type SimPart interface {
	Tick(ctx SimContext) bool
}

type SimContext struct {
	NetByPin     func(core.PinID) int // net by pin value.
	NetState     func(int) int        // net state value.
	NowMicros    uint64               // monotonic simulated time since sim start, in microseconds.
	EnableJitter bool                 // enable jitter value.
	Rand         *rand.Rand           // rand value.
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
	Edges []StateEdge // edges value.
}

type StateEdge struct {
	FromNet int // from net value.
	ToNet   int // to net value.
	Drive   int // drive value.
}

type VectorAsset struct {
	Name string // display name.
}

// Draw draws vector art with symbol-centred SVG coordinates mapped via [SVGLocalToWorld].
func (asset VectorAsset) Draw(ctx DrawContext, base core.BasePart) {
	_ = drawGeneratedVectorAsset(asset.Name, ctx, base)
}
