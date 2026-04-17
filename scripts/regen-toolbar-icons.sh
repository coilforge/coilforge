#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ICON_SOURCES_DIR="$ROOT_DIR/icon-sources"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/regen-toolbar-icons.sh

Regenerates runtime toolbar PNG assets from icon-sources/btn_*.svg.
Current mapping:
  btn_relay.svg -> internal/part/catalog/relay/toolbar_icon.png
  btn_vcc.svg -> internal/part/catalog/power/toolbar_icon_vcc.png
  btn_gnd.svg -> internal/part/catalog/power/toolbar_icon_gnd.png
  btn_button.svg -> internal/part/catalog/switches/toolbar_icon.png
  btn_indicator.svg -> internal/part/catalog/indicator/toolbar_icon.png
  btn_diode.svg -> internal/part/catalog/diode/toolbar_icon.png
  btn_rch.svg -> internal/part/catalog/rch/toolbar_icon.png
  btn_clock.svg -> internal/part/catalog/clock/toolbar_icon.png
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
  echo "resvg $src_rel -> $dst_rel (84x84)"
  resvg "$src" "$dst" --width 84 --height 84
}

main() {
  if [[ "${1:-}" =~ ^(-h|--help|help)$ ]]; then
    usage
    return
  fi

  require_cmd resvg

  if [[ ! -d "$ICON_SOURCES_DIR" ]]; then
    echo "Missing icon source directory: icon-sources/" >&2
    exit 1
  fi

  render_toolbar_icon "icon-sources/btn_relay.svg" "internal/part/catalog/relay/toolbar_icon.png"
  render_toolbar_icon "icon-sources/btn_vcc.svg" "internal/part/catalog/power/toolbar_icon_vcc.png"
  render_toolbar_icon "icon-sources/btn_gnd.svg" "internal/part/catalog/power/toolbar_icon_gnd.png"
  render_toolbar_icon "icon-sources/btn_button.svg" "internal/part/catalog/switches/toolbar_icon.png"
  render_toolbar_icon "icon-sources/btn_indicator.svg" "internal/part/catalog/indicator/toolbar_icon.png"
  render_toolbar_icon "icon-sources/btn_diode.svg" "internal/part/catalog/diode/toolbar_icon.png"
  render_toolbar_icon "icon-sources/btn_rch.svg" "internal/part/catalog/rch/toolbar_icon.png"
  render_toolbar_icon "icon-sources/btn_clock.svg" "internal/part/catalog/clock/toolbar_icon.png"
}

main "$@"
