# Toolbar And Placement TODO

## Current Status
- Toolbar chrome is in place and stable.
- Left placement toolbar hides in sim mode; right command strip is always visible.
- Hover + active visuals and icon/label rendering are implemented.
- Toolbar clicks are consumed in screen space and do not leak into world clicks.

## Completed (Archived)
- [x] Two-toolbar screen layout and fixed geometry constants.
- [x] Left/right grouped toolbar data in `internal/app/toolbar.go`.
- [x] Hover + active visual states.
- [x] Toolbar click interception before world-space input handling.
- [x] Receiver naming normalization in catalog part files (`self`).
- [x] Local static checks passing (`./check.sh go`).

## Backburner
- [ ] Decide and lock final spacing against prop panel + status bar.
- [ ] Decide whether placement order should be wire-first.
- [ ] Add pressed/down visual state distinct from hover/active.
- [ ] Decide final button content policy: icon-only vs icon+always-label.
- [ ] Add right-toolbar command icons (if desired/available).
- [ ] Review/trim old generated-vector stubs if the generator pipeline is revived.

## Next Focus: Basic Placement On Schematic

### Scope
- Keep app-level orchestration in `internal/app`.
- Keep world-space edit actions in `internal/editor`.
- Keep chrome hit-testing/drawing in `internal/render`.

### Checklist
- [ ] Left-toolbar click selects placement tool explicitly (not only hotkeys).
- [ ] Selected placement tool is reflected immediately in left-toolbar active state.
- [x] Placement ghost follows pointer continuously in edit mode (grid-snapped preview).
- [ ] Clicking schematic in edit mode places the selected part at grid-snapped world position.
- [ ] Placement click path works after toolbar click consumption (no dead input paths).
- [ ] ESC cancels current placement tool cleanly.
- [ ] Switching to run mode cancels/blocks placement interactions.
- [ ] Right-toolbar command clicks are safely no-op placeholders (no accidental editor/sim side effects yet).

### Verification
- [ ] Manual check: select part from toolbar, place part on canvas, repeat.
- [ ] Manual check: toolbar click does not select/move world parts.
- [ ] Manual check: run mode blocks left-toolbar placement flow.
- [ ] `go test ./...` passes.
- [ ] `./check.sh go` passes.

## Parts / catalog draw (follow-up)

- **Wire:** treat as an almost internal-only special case (routing/editor-centric; not the same vector-baked stem pattern as discrete components).
- **Relay:** special structure handling today; patterns there might theoretically be reused by another authored part, but nothing obvious yet — revisit if a second relay-like type appears.
- [ ] **Extract generic vector draw helpers:** move most repeated `draw.go` boilerplate (layout name suffix, `drawBase` with rotation cleared, bounds/anchors/hit-test/draw wiring using generated `vectors_gen.go` / `pins_gen.go`) into a **shared helper** under `internal/part` or a thin subpackage, so catalog packages mostly supply stem selection + type-specific quirks (relay/wire excluded or handled via small hooks).
