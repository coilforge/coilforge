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
