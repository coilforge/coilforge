# CoilForge

[![CI](https://github.com/coilforge/coilforge/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/coilforge/coilforge/actions/workflows/ci.yml)
[![golangci-lint](https://img.shields.io/badge/linted_with-golangci--lint-00ADD8?logo=golangci-lint&logoColor=white)](https://golangci-lint.run/)
![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/coilforge/coilforge?utm_source=oss&utm_medium=github&utm_campaign=coilforge%2Fcoilforge&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)
[![Go 1.24](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)

CoilForge is an Ebiten-based schematic editor and simulator for relay logic style circuits. The current app supports part placement, wire drawing, selection and transforms, save/load to JSON, and edit/run mode switching.

## Current part catalog

- Relay
- VCC
- GND
- Switch
- Indicator
- Diode
- RCH
- Clock
- Wire

## Run locally

Requirements:

- Go 1.24.x or newer compatible with `go.mod`
- Node.js 22 for repo-local Markdown lint tooling
- Linux only: X11/OpenGL development packages matching the CI workflow when building/running in a Linux environment

Start the desktop app:

```bash
go run ./cmd/coilforge
```

## Checks and lint

The repository uses repo-local tooling through [`check.sh`](/Users/mats/Projects/CoilForge/check.sh).

Bootstrap local tools:

```bash
./check.sh bootstrap
```

Run the same core checks used in CI:

```bash
./check.sh full
```

Available check commands:

- `./check.sh go` runs `golangci-lint` with [`.github/golangci.yml`](/Users/mats/Projects/CoilForge/.github/golangci.yml)
- `./check.sh markdown` runs `markdownlint-cli2` against Markdown files
- `./check.sh architecture` runs the optional architecture review script
- `./check.sh full` runs Go lint and Markdown lint, and includes architecture review when `COILFORGE_ARCH_LLM_CMD` is configured

## Generated vector code

Schematic SVGs live under `internal/part/catalog/<pkg>/assets/`. The generator (`scripts/gen_part_vectors.go`, invoked via `./scripts/regen-part-vectors.sh` or `go run -tags genpartvectors ./scripts/gen_part_vectors.go <pkg>`) writes:

- **`vectors_gen.go`** — registered draw funcs, pin layouts from red marker positions, hit bounds.
- **`pins_gen.go`** (when circles have `id="..."`) — `…PinIDs` struct fields named after those ids (exported identifiers, e.g. `TerminalA`), JSON tags, `…PinMarkerMap(*T)` for `Anchors`, and `assignNew…Pins` for clone pin allocation — embed the struct in your part type and use `self.TerminalA`-style fields in sim code.

**Committed generated files:** these `*_gen.go` outputs are checked into git so a plain `go build` / clone does not require running codegen, PRs can show the full emitted diff, and git history reflects the exact registrations shipped. After editing SVGs or the generator, regenerate and commit the updated generated files with those changes.

## CI notes

GitHub Actions in [`.github/workflows/ci.yml`](/Users/mats/Projects/CoilForge/.github/workflows/ci.yml) currently:

- runs on pushes to `main` and `master`, plus pull requests and manual dispatch
- uses Go `1.26.2` in CI
- uses Node.js `22`
- installs Linux native dependencies required by Ebiten/GLFW
- bootstraps local tools with `./check.sh bootstrap`
- runs `./check.sh full`

## Controls

- `1` to `8`: start placing relay, VCC, GND, switch, indicator, diode, RCH, and clock parts
- `W`: toggle wire mode
- `R`: rotate selection
- `M`: mirror selection
- `Delete` or `Backspace`: delete selection
- `Z` / `Y`: undo / redo
- `C` / `V`: copy / paste selected parts
- `L`: edit the selected part label
- `Enter`: commit label edits
- `F5`: toggle run mode
- `F6`: save to `coilforge.json`
- `F7`: load from `coilforge.json`
- `Escape`: clear the current transient tool state

## Architecture

Architecture notes live in [`Docs/ARCHITECTURE.md`](Docs/ARCHITECTURE.md), with compact contributor rules in [`Docs/ARCHITECTURE_RULES.md`](Docs/ARCHITECTURE_RULES.md).
