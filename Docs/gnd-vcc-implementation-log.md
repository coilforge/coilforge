# GND / VCC catalog — implementation log (for review)

## Source approach

- Duplicated `internal/part/catalog/indicator` to `gnd` and `vcc` and edited in place (not a clean-slate rewrite), following the same file layout: `part.go`, `draw.go`, `sim.go`, `props.go`, `assets.go`, generated `vectors_gen.go` + `pins_gen.go`.
- Placeholder **schematic** SVGs: `gnd/assets/gnd.svg`, `vcc/assets/vcc.svg` (one red pin marker `TerminalA` for vector regen). **Final art** is expected to be supplied by the author; keep the `TerminalA` id (or update pin ids and re-run `scripts/regen-part-vectors.sh`).
- **Toolbar** bitmap: still the **copied indicator** `assets/toolbar_icon.png` in each package—replace with GND/VCC-specific icons when ready.

## Code changes (by file / area)

| Area | Change |
|------|--------|
| `internal/part/part.go` | `NetSeeder` now takes `union NetUnion` as the first argument so seeders can use `Find(netID)` and key `high` / `low` by **union root** (required once wires merge nets). |
| `internal/sim/sim.go` | `SeedNets(union, netByPin, high, low)` — pass the same `*unionFind` as `part.NetUnion`. |
| `internal/part/catalog/gnd` | `Gnd` + `GndPinIDs`, `NetSeeder` only: `low[root]=true` for `TerminalA`. No `SimPart` (no `Tick`). `draw` uses vector stem `gnd-0`…`gnd-3` (4 rotation slots, no `off`/`on` split). |
| `internal/part/catalog/vcc` | `Vcc` + `VccPinIDs`, `NetSeeder` only: `high[root]=true` for `TerminalA`. |
| `internal/partmanifest/all.go` | Blank imports for `gnd` and `vcc`; placement `2` / `3` after indicator. |
| `scripts/regen-part-vectors.sh` | Regen list includes `gnd` and `vcc`. |
| `Docs/ARCHITECTURE.md` | `NetSeeder` signature and example loops updated. |

## Regeneration

```bash
./scripts/regen-part-vectors.sh
# or:
go run -tags genpartvectors ./scripts/gen_part_vectors.go gnd
go run -tags genpartvectors ./scripts/gen_part_vectors.go vcc
```

## Issues / follow-ups

1. **Conflict**: If one merged net is seeded both **high** and **low**, `resolveFromSeeds` yields **short** (`NetShort`) — intentional for conflicting rails on the same net.
2. **Placeholder SVGs**: Geometry is minimal; pins and bounds will change when final SVGs land—**re-run vector gen** after replacing assets.
3. **`union.Find(-1)`**: Seeder skips unconnected pins (`netByPin` &lt; 0).
4. **Hotkeys**: Manifest uses `2` and `3`; app already maps keys `2`/`3` to placement tools.
5. **Wiring**: This tree has no `wire` part in the catalog yet, so you still cannot build a full GND–VCC–indicator path in the editor until wires (or another conductor) exist. GND/VCC + `NetSeeder` are in place for when that lands.
