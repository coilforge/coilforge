#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

ATLAS_PNG="$ROOT_DIR/internal/render/fonts/ui_label_atlas.png"
ATLAS_JSON="$ROOT_DIR/internal/render/fonts/ui_label_atlas.json"

if [[ ! -f "$ATLAS_PNG" || ! -f "$ATLAS_JSON" ]]; then
  echo "UI label atlas missing. Regenerating..."
  "$ROOT_DIR/scripts/regen-ui-font-atlas.sh"
fi

exec go run "$ROOT_DIR/cmd/coilforge" "$@"
