# Two Toolbars TODO

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
- `internal/app/app.go`
- `internal/app/toolbar.go`
- `internal/render/chrome.go`
- Optional: command/helper files in `internal/app` if the right toolbar needs explicit action wrappers.
- Optional later: part catalog icon hooks when milestone 3 begins.

## Milestone 1: Draw Toolbars And Plain Buttons
- [ ] Define a fixed screen-space layout for two vertical toolbars.
- [ ] Reserve the left side for placement tools.
- [ ] Reserve the right side for app/control commands.
- Define exact toolbar geometry:
- [ ] Toolbar outer width and height rules.
- [ ] Button width and height.
- [ ] Button-to-button spacing.
- [ ] Toolbar padding and screen-edge margins.
- [ ] Hide/show rules for edit mode vs sim mode.
- [ ] Define exact hitbox rules for each button in screen space.
- [ ] Decide whether the hitbox matches the visible button rect exactly or includes a small tolerance margin.
- [ ] Decide exact spacing, button size, and margins so toolbars do not collide with the property panel or status bar.
- [ ] Build grouped toolbar data in `internal/app/toolbar.go`.
- [ ] Left group: wire first, then part-type buttons.
- [ ] Right group: run/stop, step, load, save, home, settings, and similar commands as visual placeholder buttons only.
- [ ] Draw both toolbars in `internal/render/chrome.go` using plain text buttons.
- [ ] Show the left toolbar only in edit mode.
- [ ] Keep the right toolbar visible in both edit mode and sim mode unless a specific command needs to be hidden.

### Checkpoint 1
- [ ] Both toolbars render in stable positions.
- [ ] Buttons are plain rectangles with text labels only.
- [ ] Toolbar and button dimensions are explicitly defined in code, not implied by ad hoc drawing math.
- [ ] Button hitboxes are explicitly defined in code, even before real functionality is attached.
- [ ] Left toolbar disappears during sim.
- [ ] Rendering the chrome does not affect scene drawing underneath.

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
- [ ] Add icons to left-toolbar tool buttons.
- [ ] Add icons to right-toolbar command buttons where assets exist.
- [ ] Keep text labels alongside icons unless the UI proves clear without them.
- [ ] Preserve hover, pressed, active, and disabled visuals after icons are added.

### Checkpoint 3
- [ ] Icons render crisply and consistently.
- [ ] Buttons remain readable with labels.
- [ ] Active and disabled states still read clearly with icons present.

## Final Verification
- [ ] Edit mode: both toolbars behave correctly.
- [ ] Sim mode: left toolbar is hidden.
- [ ] Right toolbar still behaves correctly in sim mode.
- [ ] Toolbar hitboxes consume clicks without changing real editor or sim behavior.
- [ ] Build and lint checks pass for touched files.
- [ ] Keep the result explicit and simple, with no unnecessary abstraction.
