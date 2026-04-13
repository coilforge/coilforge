#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARTIFACT_DIR="$ROOT_DIR/.artifacts"
PROMPT_FILE="$ARTIFACT_DIR/architecture-review.prompt.md"
REPORT_FILE="$ARTIFACT_DIR/architecture-review.report.md"
RULES_FILE="$ROOT_DIR/Docs/ARCHITECTURE_RULES.md"
ARCH_FILE="$ROOT_DIR/Docs/ARCHITECTURE.md"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/check-architecture.sh

Required environment:
  One of:
  - COILFORGE_ARCH_LLM_CMD
      Shell command that reads a review prompt from stdin and writes a report to stdout.
  - COPILOT_GITHUB_TOKEN with the `copilot` CLI installed
      Uses GitHub Copilot CLI in programmatic mode.

Optional environment:
  COILFORGE_ARCH_DIFF_BASE
    Git revision to diff against. Defaults to the merge base with origin/main when available,
    otherwise HEAD~1, otherwise HEAD.

Notes:
  - The report must contain one of:
      VERDICT: PASS
      VERDICT: FAIL
EOF
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

resolve_base_ref() {
  if [[ -n "${COILFORGE_ARCH_DIFF_BASE:-}" ]]; then
    printf '%s\n' "$COILFORGE_ARCH_DIFF_BASE"
    return
  fi

  if git rev-parse --verify origin/main >/dev/null 2>&1; then
    git merge-base HEAD origin/main
    return
  fi

  if git rev-parse --verify HEAD~1 >/dev/null 2>&1; then
    printf 'HEAD~1\n'
    return
  fi

  printf 'HEAD\n'
}

collect_changed_files() {
  local base_ref="$1"
  git diff --name-only "$base_ref"...HEAD -- '*.go' '*.md' '.github/workflows/*.yml' '.github/workflows/*.yaml' 'check.sh' 'scripts/*.sh'
}

write_prompt() {
  local base_ref="$1"
  local changed_files="$2"

  mkdir -p "$ARTIFACT_DIR"

  {
    echo "# Task"
    echo
    echo "Review this repository for violations of the compact architecture rules."
    echo "Report only architectural issues or likely issues."
    echo
    echo "# Verdict Format"
    echo
    echo "End the report with exactly one of:"
    echo "- VERDICT: PASS"
    echo "- VERDICT: FAIL"
    echo
    echo "# Compact Rules"
    echo
    cat "$RULES_FILE"
    echo
    echo "# Full Architecture Reference"
    echo
    cat "$ARCH_FILE"
    echo
    echo "# Changed Files Compared To $base_ref"
    echo
    if [[ -n "$changed_files" ]]; then
      printf '%s\n' "$changed_files"
    else
      echo "(no changed files detected; review the current repository state)"
    fi
    echo
    echo "# Important File Excerpts"
    echo
  } > "$PROMPT_FILE"

  if [[ -n "$changed_files" ]]; then
    while IFS= read -r file; do
      [[ -z "$file" ]] && continue
      [[ ! -f "$ROOT_DIR/$file" ]] && continue

      {
        echo
        echo "## File: $file"
        echo
        echo '```'
        sed -n '1,240p' "$ROOT_DIR/$file"
        echo '```'
      } >> "$PROMPT_FILE"
    done <<< "$changed_files"
  else
    while IFS= read -r file; do
      {
        echo
        echo "## File: $file"
        echo
        echo '```'
        sed -n '1,200p' "$ROOT_DIR/$file"
        echo '```'
      } >> "$PROMPT_FILE"
    done < <(find "$ROOT_DIR/internal" -name '*.go' | sort)
  fi
}

run_review() {
  if [[ -n "${COILFORGE_ARCH_LLM_CMD:-}" ]]; then
    /bin/sh -lc "$COILFORGE_ARCH_LLM_CMD" < "$PROMPT_FILE" > "$REPORT_FILE"
  elif command -v copilot >/dev/null 2>&1 && [[ -n "${COPILOT_GITHUB_TOKEN:-}" ]]; then
    local prompt
    prompt="$(cat "$PROMPT_FILE")"
    copilot -p "$prompt" --no-ask-user > "$REPORT_FILE"
  else
    echo "No architecture review backend configured." >&2
    echo "Set COILFORGE_ARCH_LLM_CMD or install Copilot CLI and set COPILOT_GITHUB_TOKEN." >&2
    exit 1
  fi

  cat "$REPORT_FILE"

  if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
    {
      echo "## Architecture Review"
      echo
      cat "$REPORT_FILE"
    } >> "$GITHUB_STEP_SUMMARY"
  fi

  if grep -q 'VERDICT: FAIL' "$REPORT_FILE"; then
    exit 1
  fi

  if ! grep -q 'VERDICT: PASS' "$REPORT_FILE"; then
    echo "Architecture review report did not contain a valid verdict." >&2
    exit 1
  fi
}

main() {
  if [[ "${1:-}" =~ ^(-h|--help|help)$ ]]; then
    usage
    return
  fi

  require_cmd git
  require_cmd sed

  local base_ref
  local changed_files

  base_ref="$(resolve_base_ref)"
  changed_files="$(collect_changed_files "$base_ref" || true)"

  write_prompt "$base_ref" "$changed_files"
  run_review
}

main "$@"
