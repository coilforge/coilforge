#!/usr/bin/env bash
# Regenerates internal/part/catalog/<pkg>/vectors_gen.go via gen_part_vectors.go
# (full SVG paint: #hex, #RRGGBBAA, rgb/rgba, transparent, currentColor, SVG color names).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

CATALOG_PKGS=(indicator gnd vcc clock relay)

for pkg in "${CATALOG_PKGS[@]}"; do
  echo "# gen vectors: $pkg"
  go run -tags genpartvectors ./scripts/gen_part_vectors.go "$pkg"
done
