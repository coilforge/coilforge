#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_FONT="$(mktemp -t coilforge-inter-XXXXXX.ttf)"
cleanup() {
  rm -f "$TMP_FONT"
}
trap cleanup EXIT

INTER_URL="https://github.com/google/fonts/raw/main/ofl/inter/Inter%5Bopsz,wght%5D.ttf"
echo "Downloading Inter variable font..."
curl -L "$INTER_URL" -o "$TMP_FONT"

go run "$ROOT_DIR/scripts/gen_ui_font_atlas.go" \
  --font "$TMP_FONT" \
  --out-png "$ROOT_DIR/internal/render/fonts/ui_label_atlas.png" \
  --out-json "$ROOT_DIR/internal/render/fonts/ui_label_atlas.json" \
  --size 10 \
  --dpi 72 \
  --tex-w 512 \
  --padding 1
