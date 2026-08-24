# Saint Hubbins — Overview & Scope

## What Saint Hubbins Is

Saint Hubbins is a **Go-native live-coding music environment**. Patterns are pure functions `State -> [Hap]` queried over rational time; output drives offline audio, MIDI, OSC, and the live console.

- **Module**: `codeberg.org/uzu/saint-hubbins`, Go 1.25
- **CLI**: `saint-hubbins` (alias `hubbins`) — `eval`, `serve`, `render`, `query`
- **Live console**: Go HTTP server at `http://localhost:8080` (`web/server.go`, `web/templates/console.html`, `web/static/saint-hubbins.wasm`)
- **License**: AGPL-3.0-or-later

## Goals

> Build a complete Go pattern engine with idiomatic replacements where browser APIs do not exist.

Not a port of any single JS framework — Saint Hubbins reimagines the pattern model in Go with its own vocabulary (mini notation, controls, transforms) and its own console.

## Two Execution Targets

1. **Native Go binary** — CLI `saint-hubbins` evaluates patterns, renders audio offline, and serves the live console over HTTP.
2. **WASM + Go server** — Pattern engine compiled to WASM (`GOOS=js GOARCH=wasm`) for the browser console. Go serves `saint-hubbins.wasm`.

## Constraints

- AGPL-3.0 compliance — `ATTRIBUTION.md` + `LICENSE`.
- No behavior drift in core query semantics without a migration note.
- No JS runtime embedded in production.

## See Also

- `ATTRIBUTION.md` — upstream inspiration credit (TidalCycles-inspired prototype / TidalCycles) retained in one place.
- `docs/archive/strudel-legacy/` — historical planning docs from the TidalCycles-inspired prototype, kept for reference only. Not part of the Saint Hubbins surface.
