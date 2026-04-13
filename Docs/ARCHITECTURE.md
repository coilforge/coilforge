# CoilForge v2 — Architecture Design

This document is the architectural blueprint for a ground-up rewrite of
CoilForge. The rewrite keeps the same outward feature set — schematic editing,
simulation, the same component catalog — but replaces the internal architecture
with something immediately understandable.

## Design Philosophy

1. **Code clarity over design patterns.** Every file, struct, and function
   should be understandable on first read by a new contributor. No closure
   dispatch tables, no generic store bridges, no layered service bundles.

2. **Concrete methods, not closures.** Behavior lives in methods on concrete
   structs. Parts implement interfaces with real methods, not function-field
   structs filled at registration time.

3. **Direct struct access.** We do not protect data from code. Public fields
   are fine. Passing a `*Relay` around and reading its `Poles` slice is
   preferred over accessor methods that add nothing.

4. **Globals for shared state.** Camera, the parts list, net states, and mode
   flags are used everywhere. They live in a shared `world` package as
   package-level variables. Don't pass them through five call layers.

5. **Hierarchical packages for humans.** Group child packages under the parent
   they belong to. A contributor should see `part/catalog/relay/` and
   immediately know that relay is a part. This is for readability, not access
   restriction.

6. **World coordinates everywhere.** Except for toolbar/chrome, all schematic
   interaction uses world coordinates. Screen-to-world conversion happens once
   at the frame boundary.

7. **Editor and simulator are independent.** They share world state but do not
   import each other. An orchestrator dispatches to one or the other.

8. **Wires are first-class parts.** Selection, move, copy, paste, delete,
   rotate, undo/redo — all go through the same Part interface as every other
   component. Only wire *routing* (continuous segment drawing) is a special
   editor mode.

9. **Parts own their state and behavior.** Each instantiated part holds its
   own state. Drawing is done by the part itself, using SVG artwork
   one-time-converted into generated Ebiten vector commands. The simulator
   does not reach into parts to mutate them — it asks each part to compute
   its next state given its current inputs (net states, time/tick). The part
   returns or applies its own update.

---

## Package Layout

```
cmd/
  coilforge/              main entrypoint

internal/
  core/                   value types, geometry, IDs
  world/                  shared mutable state: parts, camera, nets, mode
  part/                   Part interface, sim interfaces, registration
    catalog/
      relay/              concrete component
      wire/               concrete component
      clock/              concrete component
      switches/           concrete component
      indicator/          concrete component
      power/              concrete component (vcc + gnd subtypes)
      diode/              concrete component
      rch/                concrete component
      skeleton/           template for new components
  components/             blank imports for side-effect registration
  editor/                 schematic editing
  sim/                    simulation engine
  flatten/                net derivation, future module expansion
  render/                 all drawing: grid, parts, overlays, chrome
  app/                    Ebiten host, input, mode dispatch, file I/O
```

### Dependency Graph

Arrows mean "imports." Read top-down: lower packages are imported by higher
ones.

```
                        cmd/coilforge
                             │
                            app
                ┌──────┬─────┼──────┬─────────┐
             editor   sim  flatten render  components
                │      │      │      │          │
                ├──────┼──────┼──────┘     catalog/*
                │      │      │                 │
                └──────┴──┬───┘                 │
                        world                   │
                          │                     │
                         part ◄─────────────────┘
                          │
                         core
```

Rules:

- `core` is a leaf. No internal imports.
- `part` imports only `core`.
- `world` imports `core` and `part` (for the `Part` interface).
- `part/catalog/*` imports `core` and `part`. Does NOT import `world`.
- `components` blank-imports all `part/catalog/*` packages.
- `editor`, `sim`, `flatten`, `render` each import `core`, `part`, `world`.
- `editor` does NOT import `sim`. `sim` does NOT import `editor`.
- `app` imports everything except concrete catalog packages.
- `cmd/coilforge` imports only `app` and `components`.

---

## Core Types (`core/`)

Pure value types and geometry. No mutable state, no Ebiten imports.

```go
// Pt is a 2D point in world space.
type Pt struct{ X, Y float64 }

// Seg is a line segment between two points.
type Seg struct{ A, B Pt }

// Rect is an axis-aligned rectangle.
type Rect struct{ Min, Max Pt }

// Identity types.
type PartTypeID string
type PinID      int
type NetID      int

// BasePart is the common authored state embedded in every part.
type BasePart struct {
    ID       int
    TypeID   PartTypeID
    Pos      Pt
    Rotation int   // 0, 1, 2, 3 (quarter turns)
    Mirror   bool
    Label    string
}

// PinAnchor is a pin's world position plus its ID.
type PinAnchor struct {
    Pt    Pt
    PinID PinID
}

// Net is a group of connected pins.
type Net struct {
    ID   int
    Pins []PinID
    Segs []Seg    // wire geometry belonging to this net
}

// NetState describes a net's resolved electrical state.
const (
    NetFloat = 0
    NetLow   = 1
    NetHigh  = 2
    NetShort = 3
)
```

Geometry helpers in `core/`:

- `LocalToWorld(base BasePart, local Pt) Pt`
- `WorldToLocal(base BasePart, world Pt) Pt`
- `PointInRect(pt Pt, r Rect) bool`
- `PointNearSeg(pt Pt, seg Seg, tolerance float64) bool`
- `RotateRect(r Rect, rotation int, mirror bool) Rect`

---

## Shared World State (`world/`)

Package-level variables for state that many packages need. This is the single
source of truth for the schematic and the camera.

```go
package world

import (
    "coilforge/internal/core"
    "coilforge/internal/part"
)

// Schematic — the canonical parts list.
var Parts      []part.Part
var NextPartID int
var NextPinID  core.PinID

// Camera / viewport.
var Cam        core.Pt     // world position of screen center
var Zoom       float64     // scale factor (default 1.0)
var ScreenW    int         // updated by app each frame
var ScreenH    int

// Mode.
var RunMode    bool

// Simulation results — populated during sim, empty during editing.
var Nets       []core.Net
var NetStates  map[int]int    // net ID → NetFloat/Low/High/Short
var SimTick    uint64

// ID allocation helpers.
func AllocPartID() int           { NextPartID++; return NextPartID - 1 }
func AllocPinID() core.PinID     { NextPinID++; return NextPinID - 1 }

// Coordinate conversion.
func ScreenToWorld(sx, sy int) core.Pt {
    return core.Pt{
        X: (float64(sx)-float64(ScreenW)/2)/Zoom + Cam.X,
        Y: (float64(sy)-float64(ScreenH)/2)/Zoom + Cam.Y,
    }
}

func WorldToScreen(pt core.Pt) (float64, float64) {
    return (pt.X-Cam.X)*Zoom + float64(ScreenW)/2,
           (pt.Y-Cam.Y)*Zoom + float64(ScreenH)/2
}
```

### Why Globals

Camera, the parts list, and net states are needed by editor, sim, flatten,
and render. Passing them through every function call adds noise without adding
safety. This is a single-instance desktop app — one schematic, one camera, one
simulation at a time.

---

## Part System (`part/`)

Defines the contracts that all components implement. No concrete part code
lives here — only interfaces, registration infrastructure, and shared helpers.

### Part Interface

```go
// Part is the contract every component implements.
type Part interface {
    Base() *core.BasePart

    // Geometry — all in world coordinates, accounting for pos/rot/mirror.
    Bounds() core.Rect
    Anchors() []core.PinAnchor
    HitTest(pt core.Pt) HitResult
    Segments() []core.Seg          // wire geometry; nil for non-wire parts

    // Drawing.
    Draw(ctx DrawContext)

    // Editing.
    Clone(newID int, allocPin func() core.PinID) Part

    // Properties.
    PropSpec() PropSpec
    ApplyProp(action PropAction) bool

    // Serialization.
    MarshalJSON() ([]byte, error)
}
```

Every operation the editor or renderer needs is a method on Part. No closure
tables, no bound registries.

### DrawContext

```go
// DrawContext carries everything a Part.Draw call needs.
// Components do NOT import world — they receive world state here.
type DrawContext struct {
    Dst      *ebiten.Image
    Cam      core.Pt
    Zoom     float64
    ScreenW  int
    ScreenH  int
    Ghost    bool           // true for placement preview
    Selected bool           // true if part is selected

    // Non-nil during simulation. Returns the resolved state for a pin's net.
    NetState func(core.PinID) int
}

// WorldToScreen converts a world point to screen pixels using this context.
func (ctx DrawContext) WorldToScreen(pt core.Pt) (float64, float64) { ... }
```

Components call `ctx.WorldToScreen(pt)` for drawing. During simulation,
`ctx.NetState` lets wire and indicator parts color themselves by net state
without importing the sim or world packages.

### HitResult

```go
type HitResult struct {
    Hit   bool
    Kind  int          // HitBody, HitLabel, HitPin
    PinID core.PinID   // valid only when Kind == HitPin
}

const (
    HitBody  = 0
    HitLabel = 1
    HitPin   = 2
)
```

### Simulation Interfaces

Parts that participate in simulation implement one or more small interfaces.
The sim checks each part with type assertions. No stubs, no base-type
embedding — a part only has the methods it actually needs.

The key principle: **the part owns its state and computes its own updates.**
The simulator is a clock. Each tick, it presents the current world state
(net states, tick number) to each part and the part decides what its next
state should be. The sim never reaches into a part's fields to mutate them
directly.

```go
// SimPart is implemented by any part that participates in simulation.
// It is the primary sim interface. The sim calls Tick() each step,
// passing the part's current inputs. The part updates its own state
// and returns whether anything changed.
type SimPart interface {
    // Tick receives the current sim context and lets the part update
    // its own internal state. Returns true if the part's state changed
    // in a way that could affect other parts (e.g. contact flipped,
    // output toggled). Time (tick) is one of the inputs — a relay
    // uses it to decide whether a scheduled transition is due.
    Tick(ctx SimContext) bool
}

// SimContext carries everything a part needs to compute its next state.
type SimContext struct {
    NetByPin     func(core.PinID) int
    NetState     func(int) int
    Tick         uint64
    TickMicros   int
    EnableJitter bool
    Rand         *rand.Rand
}

// NetSeeder marks nets as driven high or low (power sources, clocks).
// Called during net resolution before state propagation.
type NetSeeder interface {
    SeedNets(netByPin func(core.PinID) int, high, low map[int]bool)
}

// Conductor unions nets through closed contacts (relays, switches).
// The part reads its own current contact state to decide which pins
// are connected. The sim does not inspect the part's state directly.
type Conductor interface {
    AddConductive(union NetUnion, netByPin func(core.PinID) int)
}

// StateEdger adds directed state edges (diodes).
type StateEdger interface {
    AddStateEdges(netByPin func(core.PinID) int, graph *StateGraph)
}

// InputHandler handles press/release from the user during sim (switches).
type InputHandler interface {
    HandleInput(active bool) (changed, momentary bool)
    ReleaseMomentary() bool
}
```

Note: there is no `StateRefresher` or `TargetUpdater` or `Transitioner` as
separate interfaces. All of that is folded into `SimPart.Tick()`. An
indicator reads its pin net state and sets `Lit`. A relay checks whether a
scheduled transition is due at the current tick, flips its contacts if so,
and senses its coil to schedule future transitions. A wire reads its net
state and sets its display color. Each part does whatever it needs in one
call.

The sim's main loop per tick:

```go
// 1. Build conductive unions (parts declare which pins they connect).
for _, p := range world.Parts {
    if c, ok := p.(part.Conductor); ok {
        c.AddConductive(union, netByPin)
    }
}

// 2. Seed net states from power sources and clocks.
for _, p := range world.Parts {
    if s, ok := p.(part.NetSeeder); ok {
        s.SeedNets(netByPin, high, low)
    }
}

// 3. Resolve net states.
world.NetStates = resolve(union, high, low, edges)

// 4. Let each sim-capable part update itself.
anyChanged := false
for _, p := range world.Parts {
    if sp, ok := p.(part.SimPart); ok {
        if sp.Tick(ctx) {
            anyChanged = true
        }
    }
}

// 5. If any part changed state (e.g. relay contact flipped), re-resolve.
if anyChanged {
    goto step1  // simplified — loop until stable or max iterations
}
```

### NetUnion and StateGraph

```go
// NetUnion is the union-find interface for merging nets.
type NetUnion interface {
    Union(a, b int)
    Find(a int) int
}

// StateGraph accumulates directed state-propagation edges.
type StateGraph struct {
    Edges []StateEdge
}
type StateEdge struct {
    FromNet, ToNet int
    Drive          int  // NetHigh or NetLow
}
```

### Registration

```go
// TypeInfo describes a registered part type.
type TypeInfo struct {
    New    func(id int, pos core.Pt) Part           // create with defaults
    Decode func(data json.RawMessage) (Part, error)  // deserialize
    Tools  []string                                   // placement tool slots
    Icon   func() *ebiten.Image                       // toolbar icon
}

// Registry maps type IDs to their TypeInfo. Populated by init() functions.
var Registry = map[core.PartTypeID]TypeInfo{}

// Register adds a part type. Called from each component's init().
func Register(id core.PartTypeID, info TypeInfo) {
    Registry[id] = info
}
```

Each component's `init()`:

```go
func init() {
    part.Register(TypeID, part.TypeInfo{
        New:    newIndicator,
        Decode: decodeIndicator,
        Tools:  []string{"main"},
        Icon:   toolbarIcon,
    })
}
```

No generics. No StoreCodec. No RegisterStorePartType. Just a map entry.

### Property System

```go
type PropSpec struct {
    Items []PropItem
}

type PropItem struct {
    Label   string
    Kind    int        // PropText, PropInt, PropChoice, PropBool, PropAction
    Value   any
    Choices []string   // for PropChoice
    Min     int        // for PropInt
    Max     int        // for PropInt
}

type PropAction struct {
    Index    int
    NewValue any
}
```

The property panel UI (in render) reads `PropSpec` from the selected part and
renders it. User interactions produce `PropAction` values fed back through
`ApplyProp`.

---

## Component File Layout

Every component under `part/catalog/<name>/` follows the same file structure:

| File | Contents |
| ---- | -------- |
| `part.go` | `TypeID`, `init()`, struct definition, `New()`, `Decode()`, `Clone()`, `MarshalJSON()` |
| `draw.go` | `Draw()`, `Bounds()`, `Anchors()`, `HitTest()`, geometry constants |
| `props.go` | `PropSpec()`, `ApplyProp()` |
| `sim.go` | `Tick()` and any other sim interface methods (omit if part has no sim behavior) |
| `assets.go` | Asset selector — picks the right generated vector asset for current state |
| `*_gen.go` | Generated Ebiten vector commands from SVG sources (do not edit by hand) |

This layout is identical across all "real" components. A contributor opening
any component finds the same files containing the same categories of code.

Each part draws itself in `Draw()` by selecting the appropriate generated
vector asset for its current state and rendering it into the `DrawContext`.
The part owns the decision of what it looks like — the renderer just
iterates parts and calls `Draw`.

### Example: Indicator

```go
// part.go
package indicator

const TypeID core.PartTypeID = "indicator"

type Indicator struct {
    core.BasePart
    // Runtime state (zero during editing, set by sim).
    Lit bool
}

func init() {
    part.Register(TypeID, part.TypeInfo{
        New:    newIndicator,
        Decode: decodeIndicator,
        Tools:  []string{"main"},
        Icon:   toolbarIcon,
    })
}

func newIndicator(id int, pos core.Pt) part.Part {
    return &Indicator{BasePart: core.BasePart{ID: id, TypeID: TypeID, Pos: pos}}
}

func (ind *Indicator) Base() *core.BasePart { return &ind.BasePart }
func (ind *Indicator) Segments() []core.Seg { return nil }

func (ind *Indicator) Clone(newID int, allocPin func() core.PinID) part.Part {
    c := *ind
    c.ID = newID
    c.Lit = false
    // remap pins...
    return &c
}
```

```go
// sim.go — indicator implements SimPart
func (ind *Indicator) Tick(ctx part.SimContext) bool {
    net := ctx.NetByPin(ind.PinA)
    wasLit := ind.Lit
    ind.Lit = ctx.NetState(net) == core.NetHigh
    return ind.Lit != wasLit
}
```

### Wire

Wire is a first-class component. It implements `Part` exactly like any other
component, so select, move, copy, paste, delete, undo/redo all work through
the standard paths with no wire-specific editor code.

```go
package wire

type Wire struct {
    core.BasePart
    Half  core.Pt     // half-extent from center; segment = (Pos-Half, Pos+Half)
    PinA  core.PinID
    PinB  core.PinID
    // Runtime state.
    State int         // NetFloat/Low/High/Short for coloring
}

func (w *Wire) Segments() []core.Seg {
    return []core.Seg{{
        A: core.Pt{X: w.Pos.X - w.Half.X, Y: w.Pos.Y - w.Half.Y},
        B: core.Pt{X: w.Pos.X + w.Half.X, Y: w.Pos.Y + w.Half.Y},
    }}
}
```

Wire implements `SimPart` so it can update its own display color from its
net state:

```go
func (w *Wire) Tick(ctx part.SimContext) bool {
    prev := w.State
    if net := ctx.NetByPin(w.PinA); net >= 0 {
        w.State = ctx.NetState(net)
    }
    return w.State != prev
}
```

Wire has no `SeedNets`, `AddConductive`, or other sim behavior. The sim
iterates all parts, finds that wire satisfies `SimPart`, and calls `Tick`.
The wire reads its net state and updates its own display color. No filtering,
no special cases.

Wire routing (continuous segment placement) is a special editor mode. See
the Editor section.

### Relay

Relay is dynamic — its geometry depends on the current pole count. The struct
has a variable-length `Poles` slice. All geometry methods (`Bounds`,
`Anchors`, `HitTest`, `Draw`) compute their results from the current state.
No special "dynamic geometry" system is needed.

```go
package relay

type Relay struct {
    core.BasePart
    Poles      []Pole
    PickupMs   int
    ReleaseMs  int
    FlightMs   int
    JitterMs   int

    // Runtime state.
    CoilActive          bool
    Contacts            []ContactState
    PendingContacts     []ContactState
    TransitionDueTick   uint64
    TransitionScheduled bool
}

type Pole struct {
    PinCommon core.PinID
    PinNC     core.PinID
    PinNO     core.PinID
}
```

Relay implements `Conductor` and `SimPart`:

```go
func (r *Relay) AddConductive(union part.NetUnion, netByPin func(core.PinID) int) {
    // For each pole, union the common pin with NC or NO based on current
    // contact state. The relay reads its own Contacts slice — the sim
    // never inspects it directly.
}

func (r *Relay) Tick(ctx part.SimContext) bool {
    changed := false

    // 1. Fire any pending transitions that are due at this tick.
    if r.TransitionScheduled && ctx.Tick >= r.TransitionDueTick {
        for i := range r.Contacts {
            r.Contacts[i] = r.PendingContacts[i]
        }
        r.TransitionScheduled = false
        changed = true
    }

    // 2. Sense coil net state.
    coilNet := ctx.NetByPin(r.PinCoilA)
    coilHigh := ctx.NetState(coilNet) == core.NetHigh
    if coilHigh != r.CoilActive {
        r.CoilActive = coilHigh
        // Schedule a future contact state change.
        delay := r.pickupTicks(ctx.TickMicros)
        if !coilHigh {
            delay = r.releaseTicks(ctx.TickMicros)
        }
        if ctx.EnableJitter {
            delay += r.jitterTicks(ctx.TickMicros, ctx.Rand)
        }
        r.TransitionDueTick = ctx.Tick + uint64(delay)
        r.TransitionScheduled = true
        // Compute what the contacts will be after transition.
        r.computePendingContacts()
    }

    return changed
}
```

The relay manages its own transition scheduling internally. The sim just
calls `Tick()` each step. If `Tick` returns true (contacts flipped), the
sim knows it needs to re-resolve nets because `AddConductive` will now
return different unions.

---

## Editor (`editor/`)

The editor owns all schematic editing logic. It reads and mutates
`world.Parts` directly. It does not import `sim`.

### State

Editor state is package-level variables:

```go
package editor

var (
    Selection    []int        // indices into world.Parts
    HoverIndex   int          // part under cursor, -1 if none
    PlaceMode    bool
    PlaceTool    core.PartTypeID
    PlacePreview part.Part    // ghost being placed
    WireMode     bool
    WireDraft    []core.Pt    // in-progress wire path points
    Dragging     bool
    DragStart    core.Pt
    BoxSelecting bool
    BoxRect      core.Rect
    UndoStack    []Snapshot
    RedoStack    []Snapshot
    Clipboard    []part.Part
    LabelEditing bool
    LabelIndex   int          // index of part being label-edited
    LabelBuffer  []rune
)
```

### Input Handling

The app converts the mouse position to world coordinates and calls:

```go
func HandleClick(pt core.Pt, button int)
func HandleRelease(pt core.Pt, button int)
func HandleDrag(pt core.Pt)
func HandleKey(key ebiten.Key)
func HandleScroll(delta float64) // only if app delegates zoom here
```

All points are in world coordinates. The editor never does screen-to-world
conversion itself.

### Placement

```go
func StartPlacement(typeID core.PartTypeID) {
    info := part.Registry[typeID]
    PlacePreview = info.New(world.AllocPartID(), core.Pt{})
    PlaceMode = true
}

func commitPlacement(pos core.Pt) {
    pushUndo()
    PlacePreview.Base().Pos = snapToGrid(pos)
    world.Parts = append(world.Parts, PlacePreview)
    PlacePreview = nil
    PlaceMode = false
}
```

### Wire Routing

Wire insertion uses a dedicated mode. The user presses `W` to enter wire
mode, clicks to lay points, and presses `Escape` or `W` to exit.

```go
func handleWireClick(pt core.Pt) {
    snapped := snapToGridOrPin(pt)
    WireDraft = append(WireDraft, snapped)

    if len(WireDraft) >= 2 {
        from := WireDraft[len(WireDraft)-2]
        to := WireDraft[len(WireDraft)-1]
        pushUndo()
        w := wire.New(world.AllocPartID(), from, to, world.AllocPinID, world.AllocPinID)
        world.Parts = append(world.Parts, w)
    }
    // Stay in wire mode — next click continues from last point.
}
```

Each segment creates a separate Wire part in `world.Parts`. The wire mode
continues from the last placed point so the user can draw a multi-segment path
without re-entering the mode.

### Selection

Clicking a part selects it. Box-drag selects all parts whose `Bounds()`
intersect the box. Since wires implement `Part` with `Bounds()` and
`HitTest()`, they participate identically.

Wire multi-click ladder (editor-specific UX):
- Single click: select one wire
- Double click: select all physically connected wires (shared endpoints)
- Triple click: select all wires on the same derived net

### Move / Rotate / Mirror

```go
func MoveSelected(delta core.Pt) {
    pushUndo()
    for _, idx := range Selection {
        b := world.Parts[idx].Base()
        b.Pos.X += delta.X
        b.Pos.Y += delta.Y
    }
}

func RotateSelected() {
    pushUndo()
    for _, idx := range Selection {
        b := world.Parts[idx].Base()
        b.Rotation = (b.Rotation + 1) % 4
    }
}

func MirrorSelected() {
    pushUndo()
    for _, idx := range Selection {
        b := world.Parts[idx].Base()
        b.Mirror = !b.Mirror
    }
}
```

These work identically for every part type including wires.

### Delete

```go
func DeleteSelected() {
    pushUndo()
    // Remove selected indices from world.Parts (reverse order).
    sort.Sort(sort.Reverse(sort.IntSlice(Selection)))
    for _, idx := range Selection {
        world.Parts = slices.Delete(world.Parts, idx, idx+1)
    }
    Selection = nil
}
```

### Clipboard (Copy / Paste)

```go
func CopySelected() {
    Clipboard = nil
    for _, idx := range Selection {
        Clipboard = append(Clipboard, world.Parts[idx])
    }
}

func Paste(offset core.Pt) {
    pushUndo()
    pinMap := map[core.PinID]core.PinID{}
    for _, orig := range Clipboard {
        cloned := orig.Clone(world.AllocPartID(), world.AllocPinID)
        cloned.Base().Pos.X += offset.X
        cloned.Base().Pos.Y += offset.Y
        world.Parts = append(world.Parts, cloned)
    }
    // Select newly pasted parts.
}
```

Because wires implement `Clone()` like every other part, clipboard works
without wire-specific code.

### Undo / Redo

A snapshot captures the full schematic state:

```go
type Snapshot struct {
    Parts      []byte  // JSON-marshaled world.Parts
    NextPartID int
    NextPinID  core.PinID
}

func pushUndo() {
    snap := captureSnapshot()
    UndoStack = append(UndoStack, snap)
    RedoStack = nil
}

func Undo() {
    if len(UndoStack) == 0 { return }
    RedoStack = append(RedoStack, captureSnapshot())
    restoreSnapshot(UndoStack[len(UndoStack)-1])
    UndoStack = UndoStack[:len(UndoStack)-1]
}
```

Restoring a snapshot deserializes the parts list back into `world.Parts`.
This uses the same `part.Registry[typeID].Decode` path as file loading.

### Labels

Label editing is editor state. The part's `Base().Label` field is the stored
text. The editor handles the text editing lifecycle:

```go
func StartLabelEdit(partIdx int) {
    LabelEditing = true
    LabelIndex = partIdx
    LabelBuffer = []rune(world.Parts[partIdx].Base().Label)
}

func CommitLabelEdit() {
    pushUndo()
    world.Parts[LabelIndex].Base().Label = string(LabelBuffer)
    LabelEditing = false
}
```

### Overlay Drawing

The editor provides an overlay drawing function called by app:

```go
func DrawOverlays(dst *ebiten.Image) {
    // Uses render helpers for drawing.
    for _, idx := range Selection {
        render.DrawSelectionOutline(dst, world.Parts[idx].Bounds())
    }
    if PlacePreview != nil {
        PlacePreview.Draw(part.DrawContext{Dst: dst, ..., Ghost: true})
    }
    if WireMode && len(WireDraft) > 0 {
        render.DrawWireDraft(dst, WireDraft)
    }
    if BoxSelecting {
        render.DrawBoxSelect(dst, BoxRect)
    }
}
```

---

## Flattener (`flatten/`)

The flattener derives nets from the current parts. For flat schematics (no
modules), this is the only job. When modules are added later, the flattener
will also expand module instances.

### Net Derivation

```go
package flatten

func BuildNets() {
    var anchors []core.PinAnchor
    var segs []core.Seg

    for _, p := range world.Parts {
        anchors = append(anchors, p.Anchors()...)
        segs = append(segs, p.Segments()...)
    }

    world.Nets = deriveNets(anchors, segs)
}
```

The `deriveNets` algorithm:

1. Collect all pin anchors and wire segments.
2. Group pins that share the same world position (coincidence).
3. Group pins whose anchor points lie on a wire segment (pin-on-segment).
4. Merge wire segments that share endpoints.
5. Build connected components. Each component is a Net.
6. Assign net IDs and build the `PinID → NetID` mapping.

This algorithm is a pure function of anchors and segments. It lives entirely
in `flatten/` and has no Ebiten dependency.

### Pin-to-Net Mapping

After `BuildNets`, the flattener also populates a lookup that the sim uses:

```go
// PinNet maps each PinID to the NetID it belongs to.
var PinNet map[core.PinID]int
```

This could live in `world` or in `flatten`. Since sim and render both need it,
`world` is the natural home:

```go
world.PinNet = flatten.BuildPinNetMap(world.Nets)
```

### Future: Module Expansion

When modules are added:

1. Walk `world.Parts` looking for module instances.
2. Load each referenced `.cofo` file.
3. Expand module instances into concrete parts with remapped IDs.
4. Add expanded parts to a flat working list.
5. Derive nets from the expanded list.
6. Feed the flat result to the sim.

The key boundary: during editing, the parts list contains module instances as
opaque boxes with exported pins. During simulation, the flattener expands them
into a flat net list plus concrete parts.

---

## Simulator (`sim/`)

The simulator is a clock. It owns tick stepping, net resolution, and the
resolve-tick-resolve loop. It does NOT own part state — each part owns its
own state and updates itself when asked via `Tick()`.

The simulator's job:

1. Resolve net states (who is connected to what, what is driven).
2. Ask each part to update itself given the current net states and tick.
3. If any part changed (e.g. relay flipped contacts), re-resolve and repeat.
4. Advance time to the next interesting tick.

### Lifecycle

```go
package sim

func Start() {
    flatten.BuildNets()
    world.SimTick = 0
    resolveAndTick()  // initial resolve + first Tick pass
}

func Stop() {
    world.Nets = nil
    world.NetStates = nil
    world.PinNet = nil
    world.SimTick = 0
}
```

### Net Resolution

Net resolution builds the electrical state of the circuit without inspecting
any part's internal fields. It only calls the `Conductor`, `NetSeeder`, and
`StateEdger` interface methods, which parts implement by reading their own
state.

```go
func resolveNets() {
    union := newUnionFind(len(world.Nets))
    high := map[int]bool{}
    low := map[int]bool{}

    // Parts declare which of their pins are conductively connected.
    for _, p := range world.Parts {
        if c, ok := p.(part.Conductor); ok {
            c.AddConductive(union, netByPin)
        }
    }

    // Power sources and clocks declare which nets they drive.
    for _, p := range world.Parts {
        if s, ok := p.(part.NetSeeder); ok {
            s.SeedNets(netByPin, high, low)
        }
    }

    // Directed edges (diodes).
    var graph part.StateGraph
    for _, p := range world.Parts {
        if e, ok := p.(part.StateEdger); ok {
            e.AddStateEdges(netByPin, &graph)
        }
    }

    world.NetStates = resolveFromSeeds(union, high, low, &graph)
}
```

### Tick Loop

Fixed-timestep, single-threaded, deterministic.

```go
const TickMicros = 10  // 10 µs per tick

func AdvanceFrame() {
    target := world.SimTick + 1000  // ~10 ms of sim time per frame

    for world.SimTick < target {
        world.SimTick = nextInterestingTick(target)
        resolveAndTick()
    }
}

func resolveAndTick() {
    for iterations := 0; iterations < maxResolveIterations; iterations++ {
        resolveNets()

        ctx := part.SimContext{
            NetByPin:     netByPin,
            NetState:     netStateLookup,
            Tick:         world.SimTick,
            TickMicros:   TickMicros,
            EnableJitter: true,
            Rand:         simRand,
        }

        anyChanged := false
        for _, p := range world.Parts {
            if sp, ok := p.(part.SimPart); ok {
                if sp.Tick(ctx) {
                    anyChanged = true
                }
            }
        }

        if !anyChanged {
            break  // stable — no part changed state
        }
        // A part changed (e.g. relay contact flipped). Loop to re-resolve
        // nets with the new conductive state, then tick again.
    }
}
```

`nextInterestingTick` advances to the earliest tick at which something is
scheduled to happen. Each `SimPart` that has pending timed events (relays)
can be queried for its next due tick. If nothing is pending, we jump
straight to the frame target.

### Relay Timing

Relay timing is the core reason this simulator exists. In v2, the relay
manages all of its own timing internally:

1. `Tick()` checks whether a scheduled transition is due at the current tick.
   If so, the relay flips its own contacts and returns `true`.
2. `Tick()` also senses the coil net state. If it changed, the relay
   schedules a future transition (pickup/release delay + jitter).
3. Because `Tick()` returned `true`, the sim re-resolves nets.
   `AddConductive` now returns different unions reflecting the new contacts.
4. The next `Tick()` pass lets all parts (indicators, wires) see the updated
   net states and update their own display state.

The sim never knows what a "relay" is or what "contacts" are. It just sees
that a `SimPart` changed and a `Conductor` now unions different pins.

### Input Handling (Switch Press)

When the user clicks a switch during simulation:

```go
func HandleClick(pt core.Pt) {
    for i := len(world.Parts) - 1; i >= 0; i-- {
        hit := world.Parts[i].HitTest(pt)
        if !hit.Hit { continue }
        if h, ok := world.Parts[i].(part.InputHandler); ok {
            changed, _ := h.HandleInput(true)
            if changed {
                resolveAndTick()
            }
        }
        return
    }
}
```

---

## Renderer (`render/`)

The renderer owns the frame loop drawing: grid, scene assembly, and chrome.
It builds a `DrawContext` and calls `Part.Draw` on each part. **Parts draw
themselves** — the renderer does not know what a relay or indicator looks
like.

### Visual Asset Pipeline

Source artwork is SVG. At build time (or via a generator script), each SVG
is converted into a generated Go file containing Ebiten vector path commands.
At runtime, each component selects the right pre-generated asset for its
current state and draws it using Ebiten's vector API. No PNG decoding, no
runtime SVG parsing.

```text
icon-sources/<name>/*.svg
    │
    ▼  scripts/regen-component-vectors.sh
part/catalog/<name>/*_gen.go      (generated vector path data)
    │
    ▼  assets.go picks the right asset for current state
Part.Draw(ctx)                    (draws vectors onto ctx.Dst)
```

Each component's `assets.go` is a thin selector:

```go
func (ind *Indicator) asset() VectorAsset {
    if ind.Lit { return indicatorOnAsset }
    return indicatorOffAsset
}
```

The relay composes its visual from layered slices (bottom cap, pole rows,
top cap) because pre-baking every pole-count × contact-state combination
is combinatorially impractical.

### Scene Drawing

```go
func DrawScene(dst *ebiten.Image) {
    drawGrid(dst)

    ctx := part.DrawContext{
        Dst:     dst,
        Cam:     world.Cam,
        Zoom:    world.Zoom,
        ScreenW: world.ScreenW,
        ScreenH: world.ScreenH,
    }

    if world.RunMode {
        ctx.NetState = netStateLookup  // func(PinID) int
    }

    for _, p := range world.Parts {
        p.Draw(ctx)
    }
}
```

During editing, `ctx.NetState` is nil. Parts check this and use default
colors. During simulation, `ctx.NetState` returns the resolved state so
parts can color themselves. For example, a wire's `Draw` looks up the state
of its net and picks the corresponding color. An indicator checks `Lit`
(which it set during its own `Tick`). The renderer never interprets part
state — it just hands each part a canvas and a context.

### Chrome

Toolbar and property panel are screen-space. They are drawn after the
schematic scene and use screen coordinates directly.

```go
func DrawToolbar(dst *ebiten.Image, tools []ToolButton, activeTool int)
func DrawPropPanel(dst *ebiten.Image, spec part.PropSpec)
func DrawStatusBar(dst *ebiten.Image, text string)
```

### Theme

Color constants live in `render/` (or a `render/theme.go` file). Dark/light
scheme toggle is a package-level variable.

```go
var DarkMode = true

func WireColor(state int) color.RGBA { ... }
func GridColor() color.RGBA { ... }
func SelectionColor() color.RGBA { ... }
func GhostTint() color.RGBA { ... }
```

---

## App / Orchestrator (`app/`)

The thin top-level layer that owns Ebiten lifecycle, input polling, mode
dispatch, and file I/O.

### Ebiten Lifecycle

```go
type App struct{}

func (a *App) Update() error {
    world.ScreenW, world.ScreenH = ebiten.WindowSize()

    // 1. Poll input.
    mx, my := ebiten.CursorPosition()

    // 2. Handle chrome in screen space.
    if handleToolbarClick(mx, my) { return nil }

    // 3. Handle pan and zoom.
    handlePanZoom(mx, my)

    // 4. Convert to world coordinates.
    pt := world.ScreenToWorld(mx, my)

    // 5. Dispatch to active mode.
    if world.RunMode {
        sim.HandleClick(pt)
    } else {
        editor.HandleClick(pt, mouseButton)
    }

    // 6. Advance simulation.
    if world.RunMode {
        sim.AdvanceFrame()
    }

    return nil
}

func (a *App) Draw(screen *ebiten.Image) {
    render.DrawScene(screen)

    if world.RunMode {
        // sim has no overlays currently
    } else {
        editor.DrawOverlays(screen)
    }

    render.DrawToolbar(screen, toolbarItems, activeToolIdx)
    if sel := selectedPart(); sel != nil {
        render.DrawPropPanel(screen, sel.PropSpec())
    }
}
```

### Input Pipeline

```text
Raw mouse (screen pixels)
    │
    ├─► Toolbar hit? → handle in screen space, done
    │
    ├─► Pan/zoom? → update world.Cam / world.Zoom, done
    │
    └─► world.ScreenToWorld() → world point
            │
            ├─► RunMode  → sim.HandleClick(pt)
            └─► EditMode → editor.HandleClick(pt)
```

One conversion. Everything below the conversion is world coordinates.

### Mode Switching

```go
func toggleRunMode() {
    if world.RunMode {
        sim.Stop()
        world.RunMode = false
    } else {
        editor.ClearTransient()
        flatten.BuildNets()
        sim.Start()
        world.RunMode = true
    }
}
```

### File I/O

```go
func SaveProject(path string) error {
    file := FileFormat{
        Parts:      world.Parts,
        NextPartID: world.NextPartID,
        NextPinID:  world.NextPinID,
    }
    data, _ := json.MarshalIndent(file, "", "  ")
    return os.WriteFile(path, data, 0644)
}

func LoadProject(path string) error {
    data, _ := os.ReadFile(path)
    var file FileFormat
    json.Unmarshal(data, &file)
    world.Parts = file.Parts
    world.NextPartID = file.NextPartID
    world.NextPinID = file.NextPinID
    editor.Reset()
    return nil
}
```

### Toolbar

Toolbar order and hotkeys are defined in app:

```go
var toolbarItems = []ToolButton{
    {TypeID: "relay",     Hotkey: '1', Label: "Relay"},
    {TypeID: "vcc",       Hotkey: '2', Label: "VCC"},
    {TypeID: "gnd",       Hotkey: '3', Label: "GND"},
    {TypeID: "switch",    Hotkey: '4', Label: "Switch"},
    {TypeID: "indicator", Hotkey: '5', Label: "Indicator"},
    {TypeID: "diode",     Hotkey: '6', Label: "Diode"},
    {TypeID: "rch",       Hotkey: '7', Label: "RCH"},
    {TypeID: "clock",     Hotkey: '8', Label: "Clock"},
}
```

These are string type IDs. No concrete catalog imports.

---

## File Format

Clean break from the current format. JSON, one array of part records:

```json
{
  "nextPartID": 42,
  "nextPinID": 128,
  "parts": [
    {
      "type": "relay",
      "data": {
        "id": 1,
        "pos": [100, 200],
        "rotation": 0,
        "mirror": false,
        "label": "K1",
        "poles": [{ "pinCommon": 10, "pinNC": 11, "pinNO": 12 }],
        "pickupMs": 5,
        "releaseMs": 3
      }
    },
    {
      "type": "wire",
      "data": {
        "id": 2,
        "pos": [140, 200],
        "half": [40, 0],
        "pinA": 20,
        "pinB": 21
      }
    }
  ]
}
```

Each record has a `type` discriminator and a `data` payload. Deserialization
reads the type, looks up `part.Registry[type].Decode`, and passes the data.

Runtime state (e.g. relay `coilActive`, wire `state`) is included if the file
was saved during simulation. On load, any runtime state is present but inert
until the sim is started. If simpler, runtime state can be stripped on save.

---

## Coordinate System

World coordinates are the default everywhere:

- Part positions are in world space.
- Pin anchors are in world space.
- Wire segments are in world space.
- Hit testing receives world-space points.
- Part.Draw receives a DrawContext with camera info for the world-to-screen
  transform.

Screen coordinates are used only for:

- Toolbar hit detection.
- Property panel rendering.
- Status bar rendering.

The single conversion point is `world.ScreenToWorld()`, called once per frame
in the app's input pipeline.

---

## Ebiten Direct Use

Ebiten types (`*ebiten.Image`, `ebiten.Key`, etc.) are used directly. No shim
layer. The ebitshim package is not carried forward into v2.

Testing strategy for Ebiten-linked code:

- Unit tests target pure logic: geometry, net derivation, serialization,
  component registration, property mutations, state transitions.
- Ebiten-linked code (drawing, input polling, frame lifecycle) is tested
  through replay fixtures and manual verification.
- If local `go test` hangs for Ebiten-linked packages, rely on CI and manual
  local test runs outside the sandbox.

Before committing to removing the shim, verify that existing tests pass
locally without it. If a small subset of tests genuinely needs a headless
Ebiten stub, keep that stub minimal and local to the test files.

---

## Touch / Tablet Considerations

The architecture does not require touch support in v1, but avoids decisions
that would block it:

- All schematic interaction is in world space, so touch gestures map naturally
  to the same coordinate system as mouse clicks.
- Toolbar buttons use a `64x64` screen-space pitch, large enough for touch.
- The input pipeline in app can be extended to handle touch events alongside
  mouse events — both convert to world-space points before dispatch.
- Future additions:
  - On-screen tool rails for touch.
  - One-finger tap to select/place.
  - Two-finger pan.
  - Pinch zoom.
  - Larger touch hit areas for pins and small controls.

Desktop keyboard+mouse workflow remains first-class. Touch support is additive.

---

## Testing Strategy

- **Unit tests** for pure logic: `core/` geometry, `flatten/` net derivation,
  component `PropSpec`/`ApplyProp`, serialization round-trips, sim net
  resolution, relay timing.
- **Component registration tests**: verify that each component registers, can
  be created with defaults, serialized, and deserialized.
- **Replay tests**: recorded input sequences replayed against the app, with
  the final schematic state compared to expected output. Covers cross-cutting
  editor workflows.
- **Simulation tests**: start from a known schematic, run N ticks, assert net
  states and part runtime state.

Tests should not depend on Ebiten runtime for pure logic. Only integration
and replay tests need Ebiten.

---

## Non-Goals

These are explicitly out of scope for this rewrite:

- Analog / SPICE-like simulation.
- New component types beyond the current set.
- Module / hierarchical schematic support (architecture allows it, but it is
  not implemented in v2).
- Plugin or dynamic loading of components.
- Multi-document editing.
- Collaborative editing.
- Backend or cloud services.
