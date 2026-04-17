# Two Toolbars TODO

**Progress snapshot:** Milestone 1 layout and drawing are largely in place (`internal/render/chrome.go`, `internal/app/app.go`, `internal/app/toolbar.go`, `internal/partmanifest`). Two vertical strips (left placement, right commands); buttons show debug **hit** + **icon slot** outlines; renderer draws a scaled toolbar bitmap when `part.Registry[type].Icon()` returns non-nil (most catalogs still return nil). **Active** placement tool uses a thicker/highlight hit outline from `activeToolIndex`. Milestone 2 (input/hover/press) and milestone 3 (icons + labels for all buttons) remain.

## Goal
- Put two toolbars on screen while keeping package boundaries from `AGENTS.md` and `Docs/ARCHITECTURE.md`.
- Left toolbar contains the wire button plus part placement buttons.
- Right toolbar contains app/control actions such as run/stop, step, load, save, home, settings, and similar commands.
- Hide the left toolbar during simulation.
- Focus this slice on UI presence and visual behavior only, not on real command functionality.

## Scope Rules
- Keep chrome handling in screen space.
- Keep orchestration and input dispatch in `internal/app`.
- Keep drawing in `internal/render`.
- Do not move toolbar behavior into `internal/editor`.
- Do not add unnecessary global state unless layout or command state truly needs it.
- Do not wire toolbar buttons to real editor, sim, save/load, or settings functionality in this slice.

## Files To Touch
- [x] `internal/app/app.go`
- [x] `internal/app/toolbar.go`
- [x] `internal/render/chrome.go`
- [x] Optional: command/helper files in `internal/app` if the right toolbar needs explicit action wrappers. *(placeholders in `rightToolbarButtons()`.)*
- [ ] Optional later: part catalog icon hooks when milestone 3 begins. *(hooks exist; most `toolbarIcon()` still return nil.)*

## Milestone 1: Draw Toolbars And Plain Buttons
- [x] Define a fixed screen-space layout for two vertical toolbars.
- [x] Reserve the left side for placement tools.
- [x] Reserve the right side for app/control commands.
- Define exact toolbar geometry:
- [x] Toolbar outer width and height rules. *(strip width, full-height panel minus margins; see `chrome.go` constants.)*
- [x] Button width and height. *(`toolbarButtonHitPx` square hit target.)*
- [x] Button-to-button spacing. *(`toolbarButtonGapPx`.)*
- [x] Toolbar padding and screen-edge margins. *(`chromeEdgeMargin`, `toolbarPanelInnerPadPx`.)*
- [x] Hide/show rules for edit mode vs sim mode. *(Left: edit only. Right: always.)*
- [x] Define exact hitbox rules for each button in screen space. *(48×48 px hit rect; same as debug stroke.)*
- [x] Decide whether the hitbox matches the visible button rect exactly or includes a small tolerance margin. *(Exact match to drawn hit outline for now.)*
- [ ] Decide exact spacing, button size, and margins so toolbars do not collide with the property panel or status bar. *(Not explicitly laid out against prop panel yet.)*
- [x] Build grouped toolbar data in `internal/app/toolbar.go`.
- [ ] Left group: wire first, then part-type buttons. *(`partmanifest` order is relay…clock then wire last; reorder manifest if wire-first is required.)*
- [x] Right group: run/stop, step, load, save, home, settings, and similar commands as visual placeholder buttons only. *(Currently three placeholders: Run / Save / Load — expand to match goal.)*
- [x] Draw both toolbars in `internal/render/chrome.go` using plain text buttons. *(Superseded by debug hit + icon-slot outlines and optional icons; **text labels on toolbar buttons not drawn yet.**)*
- [x] Show the left toolbar only in edit mode.
- [x] Keep the right toolbar visible in both edit mode and sim mode unless a specific command needs to be hidden.

### Checkpoint 1
- [x] Both toolbars render in stable positions.
- [ ] Buttons are plain rectangles with text labels only. *(Debug outlines + optional icon in slot; labels deferred.)*
- [x] Toolbar and button dimensions are explicitly defined in code, not implied by ad hoc drawing math.
- [x] Button hitboxes are explicitly defined in code, even before real functionality is attached.
- [x] Left toolbar disappears during sim.
- [x] Rendering the chrome does not affect scene drawing underneath.

## Milestone 2: Interactive States
- [ ] Add screen-space hit testing for both toolbars in `internal/app/app.go`.
- [ ] Handle toolbar clicks before `world.ScreenToWorld()`.
- [ ] Prevent toolbar clicks from falling through into editor or sim world interactions.
- [ ] Add hover visualization for buttons.
- [ ] Add pressed visualization for buttons.
- [ ] Add active visualization using toolbar-owned visual state only, not real tool selection or real command state.
- [ ] Add disabled visualization for buttons that are intentionally present but not implemented.
- [ ] Keep text labels visible on all buttons.
- Define what visual behavior a click produces for this milestone:
- [ ] Hover on pointer enter.
- [ ] Pressed while mouse is down.
- [ ] Optional mock active state for presentation/testing only.
- [ ] Disabled state for placeholder commands.
- [ ] Decide whether clicks should do nothing or toggle mock visual state for demonstration purposes.
- [ ] Keep all toolbar interaction self-contained so no real editor or sim behavior changes.

### Checkpoint 2
- [ ] Hover, pressed, active, and disabled states are visually distinct.
- [ ] Hover and press behavior follows the defined hitboxes.
- [ ] Clicking buttons does not trigger real tool selection or real app commands.
- [ ] Disabled buttons remain non-functional.
- [ ] Toolbar clicks never leak into schematic editing or sim clicks.

## Milestone 3: Add Icons
- [ ] Audit existing `Icon` hooks and current `toolbarIcon()` implementations.
- [ ] Decide whether icons can be sourced immediately from existing generated assets or need new generated artwork.
- [ ] Add icons to left-toolbar tool buttons. *(Renderer scales `Icon()` into the slot when non-nil; most `toolbarIcon()` still return nil.)*
- [ ] Add icons to right-toolbar command buttons where assets exist.
- [ ] Keep text labels alongside icons unless the UI proves clear without them.
- [ ] Preserve hover, pressed, active, and disabled visuals after icons are added.

### Checkpoint 3
- [ ] Icons render crisply and consistently.
- [ ] Buttons remain readable with labels.
- [ ] Active and disabled states still read clearly with icons present.

## Final Verification
- [ ] Edit mode: both toolbars behave correctly.
- [x] Sim mode: left toolbar is hidden.
- [x] Right toolbar still behaves correctly in sim mode. *(Strip + placeholders render.)*
- [ ] Toolbar hitboxes consume clicks without changing real editor or sim behavior.
- [x] Build and lint checks pass for touched files.
- [x] Keep the result explicit and simple, with no unnecessary abstraction.
