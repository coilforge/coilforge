# Vector Sources

Place per-part SVG source files here before generating `*_gen.go` stubs.

This directory stores the source SVG artwork and related working files for
CoilForge part and toolbar icons.

Notes:

- these files are source art, not runtime assets
- runtime part graphics are generated into catalog-local Go/vector data
- the current icon set was produced with Claude Opus 4.6 as the artwork tool
- this directory is intentionally separate from `internal/part/catalog/.../assets`
  so the source-art workflow stays distinct from the generated runtime path

## Schematic symbol SVGs (part vector sources)

Part drawings (`*.svg` used by `scripts/regen-part-vectors.sh`, not the `btn_*.svg` toolbar set):

- **Width:** Use one nominal viewBox width for all schematic symbols (currently **512** user units). Horizontal pitch and proportions stay comparable after scaling to each part’s world bounds.
- **Height:** **Arbitrary** per symbol. Pick the smallest viewBox height that comfortably contains the artwork (e.g. a short diode may use `512×256`; a tall vertical symbol may use `512×512` or `512×128` for a thin strip).
- **Coordinates:** `gen_part_vectors.go` keeps the SVG **`viewBox` as-is** — no tight crop or origin shift. A small doodle in one corner stays small relative to the full viewport; **`part.Bounds()`** scales the whole **viewBox** rectangle onto the schematic. Padding in the viewBox is real empty space at runtime.
- **Stroke width:** Use the **same** `stroke-width` (and cap/join style) in user units across these files (e.g. `10` at width 512). Then stroke thickness scales with geometry in a consistent way when the whole symbol is scaled into the part rectangle. Do not mix different stroke widths unless you intend a different visual weight.

### Boxy SVG and similar editors

[Boxy SVG](https://boxy-svg.com) often exports primitives with **`style="fill: rgb(...); stroke: rgb(...); ..."`** instead of plain `fill` / `stroke` attributes, and uses **`<ellipse rx ry>`** for circles. `scripts/gen_part_vectors.go` handles that workflow: it merges **`style`** over XML attributes, converts **`rgb(r,g,b)`** to **`#RRGGBB`** for codegen, and treats **ellipse** like a circle using **`r = max(rx, ry)`**. Ignored metadata such as **`bx:grid`** in `<defs>` is harmless.

For pin names, prefer a stable **`id`** on the pin element. Boxy sometimes only adds **`<title>PIN1</title>`** inside the ellipse; until the generator reads `<title>`, set **Object → ID** in Boxy (or equivalent) so exports include **`id="pin_…"`**.

### Pin grid (512-wide viewBox)

- Use a **512** user-unit wide `viewBox` for normal horizontal symbols.
- Place pin centers on a **64-unit grid** so pin positions line up with a coarse lattice that is still compatible with the editor’s **16** world-unit snap (**64 = 4 × 16**).
- One convenient set of **six** interior column positions (away from the left/right frame edges) is **`x ∈ {64, 128, 192, 256, 320, 384}`**. Pins may still use **`0`**, **`512`**, or any other multiple of **64** when the symbol needs edge-aligned leads.
- **`gen_part_vectors.go` snaps** each **pin marker** (filled **`#FF0000` circle**) center to the **nearest** multiple of **64** on **X** and **Y** in SVG user units (same space as `viewBox`). Slight misses from drawing tools still land on the proper grid intersections.
- **Rounded rectangles:** `<rect rx="…" ry="…"/>` is supported (equal corners; radii clamped like SVG). **Boxy SVG** and similar export this form.
- **Cubic Bézier strokes:** `<path d="… C …"/>` segments are flattened to short polylines ( spline corners from Boxy). **`M` / `L` / `A` / `C` / `Z`** are handled in path data. Relative **`m/l/c`** and quadratic **`Q`** are not implemented—use absolute commands or flatten in the authoring tool if needed.

### Color semantics (authoring → generator)

Exactly **three** semantic colors apply to part SVGs (black / red / blue). Use **`#000000`**, **`#FF0000`**, **`#0000FF`** when picking colors directly; **`rgb(...)`** and **`#RGB`** from tools are normalized to the same values during generation.

| Role | Color | Meaning |
| --- | --- | --- |
| **Part (body)** | `#000000` | Visible schematic ink: lines, polygons, circles that are part of the symbol. |
| **Pin** | `#FF0000` | Pin markers (typically small **circles**). **Required:** stable XML `id` per pin so generated code can name anchors. |
| **Label anchor** | `#0000FF` | **Not** drawn as component ink. Marks where the **reference designator / label** should sit (see below). |

**Everything else is ignored** by the generator: any stroke or fill that is not exactly `#000000`, `#FF0000`, or `#0000FF` is skipped (construction guides, accidental grays, etc.). Use other colors for notes in Inkscape if you like; they will not appear in generated drawing code.

Black remains a good choice for the real part: it is easy to pick in any palette, reads clearly on the default canvas, and does not clash with the semantic red (pins) and blue (label anchor).

### Label anchor (optional)

When a part should place its label at a defined spot, use **`#0000FF`** so it is easy to see while editing and unambiguous for the tool.

- **Preferred:** a small **filled circle** (`fill="#0000FF"`) — the **center** `(cx, cy)` is the label anchor (e.g. center of the text box or baseline origin, depending on what the runtime implements). Radius is arbitrary; only the center matters.
- **Alternative:** a **`#0000FF` `<text>`** element (e.g. a single placeholder character). Use the **center of the text bounding box** as the anchor. A **dot is usually better**: one clear center point without depending on font metrics or alignment.
- **Optional `id`:** `cf_label` on that element if your toolchain needs a stable handle; placement can still be driven primarily by the blue color.

The generator should treat **blue** primitives as **metadata only:** do **not** emit them as schematic ink; emit a label offset from their position after the same normalization as pins and body.

Parts with **no** blue label marker can fall back to a type-defined default (e.g. above/beside `Bounds()` center) in code.

**Status:** `scripts/gen_part_vectors.go` does **not** yet filter by these colors, extract pins from red circles, or extract label anchors from blue primitives. Wiring that in is the next step once this convention is locked.

## Toolbar icon PNG generation

Use the script pipeline to regenerate toolbar PNG runtime assets from `btn_*.svg`:

```bash
./scripts/regen-toolbar-icons.sh
```

Current mapping in the script:

- `icon-sources/btn_relay.svg` -> `internal/part/catalog/relay/toolbar_icon.png` (84x84)
- `icon-sources/btn_vcc.svg` -> `internal/part/catalog/power/toolbar_icon_vcc.png` (84x84)
- `icon-sources/btn_gnd.svg` -> `internal/part/catalog/power/toolbar_icon_gnd.png` (84x84)
- `icon-sources/btn_button.svg` -> `internal/part/catalog/switches/toolbar_icon.png` (84x84)
- `icon-sources/btn_indicator.svg` -> `internal/part/catalog/indicator/toolbar_icon.png` (84x84)
- `icon-sources/btn_diode.svg` -> `internal/part/catalog/diode/toolbar_icon.png` (84x84)
- `icon-sources/btn_rch.svg` -> `internal/part/catalog/rch/toolbar_icon.png` (84x84)
- `icon-sources/btn_clock.svg` -> `internal/part/catalog/clock/toolbar_icon.png` (84x84)
