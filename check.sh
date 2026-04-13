#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TOOLS_DIR="$ROOT_DIR/.tools"
TOOLS_BIN_DIR="$TOOLS_DIR/bin"
TOOLS_NODE_DIR="$TOOLS_DIR/node"

GOLANGCI_LINT_VERSION="v2.8.0"
MARKDOWNLINT_CLI2_VERSION="0.18.1"

export PATH="$TOOLS_BIN_DIR:$TOOLS_NODE_DIR/node_modules/.bin:$PATH"

usage() {
  cat <<'EOF'
Usage:
  ./check.sh bootstrap   Install repo-local lint/check tools into ./.tools
  ./check.sh go          Run Go lint checks via golangci-lint (includes govet)
  ./check.sh markdown    Run markdownlint-cli2 against Markdown files
  ./check.sh architecture Run LLM-based architecture review
  ./check.sh full        Run all checks

Notes:
  - govet is included through golangci-lint's standard linter set.
  - bootstrap installs tools locally in ./.tools and does not touch global tool state.
  - architecture review requires either COILFORGE_ARCH_LLM_CMD or Copilot CLI with COPILOT_GITHUB_TOKEN.
EOF
}

install_go_tool() {
  local bin_name="$1"
  local package_path="$2"
  local version="$3"

  mkdir -p "$TOOLS_BIN_DIR"

  if [[ -x "$TOOLS_BIN_DIR/$bin_name" ]]; then
    return
  fi

  echo "Installing $bin_name@$version"
  GOBIN="$TOOLS_BIN_DIR" go install "${package_path}@${version}"
}

bootstrap() {
  mkdir -p "$TOOLS_BIN_DIR" "$TOOLS_NODE_DIR"

  install_go_tool "golangci-lint" "github.com/golangci/golangci-lint/v2/cmd/golangci-lint" "$GOLANGCI_LINT_VERSION"

  if [[ ! -x "$TOOLS_NODE_DIR/node_modules/.bin/markdownlint-cli2" ]]; then
    echo "Installing markdownlint-cli2@$MARKDOWNLINT_CLI2_VERSION"
    npm install --no-audit --no-fund --prefix "$TOOLS_NODE_DIR" "markdownlint-cli2@${MARKDOWNLINT_CLI2_VERSION}"
  fi
}

require_tool() {
  local tool_name="$1"
  if ! command -v "$tool_name" >/dev/null 2>&1; then
    echo "Missing required tool: $tool_name" >&2
    echo "Run ./check.sh bootstrap first." >&2
    exit 1
  fi
}

run_go_checks() {
  require_tool "golangci-lint"
  golangci-lint run --config "$ROOT_DIR/.github/golangci.yml" ./...
}

run_markdown_checks() {
  require_tool "markdownlint-cli2"
  markdownlint-cli2 \
    "**/*.md" \
    "#.git" \
    "#.tools" \
    "#node_modules"
}

run_architecture_checks() {
  "$ROOT_DIR/scripts/check-architecture.sh"
}

run_full() {
  run_go_checks
  run_markdown_checks
  if [[ -n "${COILFORGE_ARCH_LLM_CMD:-}" ]]; then
    run_architecture_checks
  fi
}

main() {
  local command="${1:-full}"

  case "$command" in
    bootstrap)
      bootstrap
      ;;
    go)
      run_go_checks
      ;;
    markdown)
      run_markdown_checks
      ;;
    architecture)
      run_architecture_checks
      ;;
    full)
      run_full
      ;;
    -h|--help|help)
      usage
      ;;
    *)
      echo "Unknown command: $command" >&2
      usage >&2
      exit 1
      ;;
  esac
}

main "$@"
