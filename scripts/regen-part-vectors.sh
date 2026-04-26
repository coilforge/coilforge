#!/usr/bin/env bash
# Regenerates internal/part/catalog/<pkg>/vectors_gen.go via gen_part_vectors.go
# (full SVG paint: #hex, #RRGGBBAA, rgb/rgba, transparent, currentColor, SVG color names).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

shopt -s nullglob

CATALOG_DIR="$ROOT_DIR/internal/part/catalog"
CATALOG_PKGS=()

for pkg_dir in "$CATALOG_DIR"/*; do
  [[ -d "$pkg_dir" ]] || continue
  pkg="$(basename "$pkg_dir")"
  assets_dir="$pkg_dir/assets"
  [[ -d "$assets_dir" ]] || continue
  svgs=("$assets_dir"/*.svg)
  (( ${#svgs[@]} > 0 )) || continue
  CATALOG_PKGS+=("$pkg")
done

IFS=$'\n' CATALOG_PKGS=($(printf '%s\n' "${CATALOG_PKGS[@]}" | sort))
unset IFS

for pkg in "${CATALOG_PKGS[@]}"; do
  echo "# gen vectors: $pkg"
  go run -tags genpartvectors ./scripts/gen_part_vectors.go "$pkg"
done
