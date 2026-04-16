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
