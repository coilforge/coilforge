# CoilForge Agent Guide v3

**Default contract.** Enforce strictly.  
Precedence: ARCHITECTURE_RULES.md > this guide > ARCHITECTURE.md.  
Unclear/conflicting/underspecified → ask immediately. Never guess.

## Core Directives
- Code must be concrete, obvious, readable on first pass.
- Prefer direct struct fields + normal methods.
- Reject: closure registries, generic bridges, service bundles, assembled-at-registration objects, heavy indirection.
- Simpler explicit code > reusable abstraction.

## Package Boundaries (Enforced)
- core: leaf. No internal imports.
- part: imports only core.
- world: imports only core + part.
- part/catalog/*: imports only core + part. Never world.
- partmanifest: blank-imports catalog packages for registration; holds placement order/hotkeys manifest.
- editor/sim/flatten/render: may import core/part/world.
- editor: no sim import.
- sim: no editor import.
- app: orchestrates. No concrete catalog imports.
- cmd/coilforge: imports only app + partmanifest.

## Shared State
- All broad state in internal/world (parts, camera, nets, mode, screen size…).
- Use package-level access. Avoid threading through long chains.

## Coordinates
- Schematic logic: world coordinates only.
- Screen→world conversion: app boundary only.
- No screen-space logic in editor/sim/flatten/render/parts.

## Parts System
- Wires = first-class part.Part.
- All ops (select/move/copy/paste/delete/rotate/mirror/undo/redo/serialize/draw) use uniform part.Part contract.
- Parts own behavior + runtime state.
- Simulator: calls interfaces only (Tick/SeedNets/AddConductive/AddStateEdges…). Never mutates internals.
- Drawing: part responsibility. Renderer supplies context → calls Part.Draw. No per-part-type drawing logic in renderer.

**Catalog layout per part type**  
part.go / draw.go / props.go / sim.go (optional) / assets.go / *_gen.go (generated, never hand-edit).

## Subsystem Rules
- Editor + simulator: independent. Share world state only. No mutual imports/calls.
- Mode switching + orchestration: app only.
- State: full-schematic snapshots preferred for undo/redo/save/load.
- I/O: app only.

## Rendering & Assets
- Parts self-draw via Part.Draw.
- Assets: pre-generated Ebiten vectors from SVGs. Thin selectors in assets.go. No runtime SVG.

## Testing Priority
- Unit tests first: geometry, net derivation, serialization, registration, props, simulation, relay timing.
- Ebiten: replay tests + manual.

## Change & Review Rules
Implement/review using:
- Simpler code, explicit methods, package clarity, world-state access, part-owned logic.
- Reject: boundary violations, screen logic in schematic layers, renderer per-part-type logic, simulator part mutation, wire-specific general cases.

**Review output format (when requested)**  
- Cite file + violated rule.
- Classify: definite violation / likely violation / no issue.
- End with exactly one:  
  VERDICT: PASS  
  or  
  VERDICT: FAIL
