# CoilForge Architecture Rules

This is the compact, review-oriented version of [ARCHITECTURE.md](./ARCHITECTURE.md).
It is meant for automated review, including LLM-based checks.

## Scope

Review the implementation against these rules:

1. Keep code concrete and obvious.
   Avoid closure-dispatch registries, generic store bridges, service bundles, and indirection-heavy patterns.

2. Put behavior on concrete types.
   Parts should expose normal methods on structs, not function fields assembled at registration time.

3. Prefer direct struct access over trivial getters/setters.
   Public fields are acceptable when they keep the code simpler.

4. Shared application state belongs in `internal/world`.
   Do not thread camera, parts, net state, or mode flags through long call chains when they are shared app state.

5. Respect the package boundaries.
   - `core` is a leaf and must not import internal packages.
   - `part` imports only `core`.
   - `world` imports only `core` and `part`.
   - `part/catalog/*` imports `core` and `part`, and must not import `world`.
   - `partmanifest` only blank-imports catalog packages for registration (and may hold placement manifest data).
   - `editor`, `sim`, `flatten`, `render` may import `core`, `part`, and `world`.
   - `editor` must not import `sim`.
   - `sim` must not import `editor`.
   - `app` may orchestrate other packages but should not import concrete catalog packages.
   - `cmd/coilforge` should import only `app` and `partmanifest`.

6. Use world coordinates for schematic logic.
   Screen-to-world conversion should happen at the app boundary, not deep in editor, sim, flatten, or parts.

7. Keep editor and simulator independent.
   They may share `world` state, but they should not directly import or call each other.

8. Treat wires as normal parts.
   Selection, move, copy, paste, delete, rotate, mirror, undo/redo, serialization, and drawing should use the same `part.Part` contract as other part types.
   Only wire routing is editor-specific.

9. Parts own their own behavior and runtime state.
   The simulator should call part interfaces such as `Tick`, `SeedNets`, `AddConductive`, and `AddStateEdges` instead of mutating part internals directly.

10. Drawing responsibility belongs to the part.
    The renderer should assemble scene/chrome state and call `Part.Draw`; it should not contain per-part-type drawing logic.

11. Keep the catalog part layout consistent.
    Real part types under `internal/part/catalog/<name>/` should follow the standard file split:
    `part.go`, `draw.go`, `props.go`, `sim.go` when needed, `assets.go`, and generated `*_gen.go`.

12. Commit generated catalog vector output (`vectors_gen.go` / `*_gen.go` from the part-vector generator) in git so the default build does not require running codegen; regenerate and commit that file whenever SVG sources or the generator change.

13. Keep file I/O and top-level orchestration in `app`.
    Package lifecycle, Ebiten lifecycle, input polling, mode switching, and save/load belong there.

14. Favor simple whole-schematic state management.
    Undo/redo and file save/load can operate on full snapshots rather than complex incremental object graphs.

## Review Output Contract

When reviewing, prefer high-signal findings over commentary.

- Report only concrete violations or likely violations.
- For each finding, cite file paths and explain which rule it violates.
- Distinguish:
  - `definite violation`
  - `likely violation`
  - `no issue`
- End with a single verdict:
  - `VERDICT: PASS`
  - `VERDICT: FAIL`
