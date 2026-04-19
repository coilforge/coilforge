#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/regen-toolbar-icons.sh

Regenerates embedded toolbar PNGs from the indicator catalog assets/btn_*.svg (168×168).
Outputs go to the same assets/ folder (matches //go:embed paths in assets.go).
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

render_toolbar_icon() {
  local src_rel="$1"
  local dst_rel="$2"
  local src="$ROOT_DIR/$src_rel"
  local dst="$ROOT_DIR/$dst_rel"

  if [[ ! -f "$src" ]]; then
    echo "Missing source SVG: $src_rel" >&2
    exit 1
  fi

  mkdir -p "$(dirname "$dst")"
  echo "resvg $src_rel -> $dst_rel (168x168)"
  resvg "$src" "$dst" --width 168 --height 168
}

main() {
  if [[ "${1:-}" =~ ^(-h|--help|help)$ ]]; then
    usage
    return
  fi

  require_cmd resvg

  render_toolbar_icon "internal/part/catalog/indicator/assets/btn_indicator.svg" "internal/part/catalog/indicator/assets/toolbar_icon.png"
}

main "$@"
