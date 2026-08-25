# Saint Hubbins — Overview & Scope

## What Saint Hubbins Is

Saint Hubbins is a **Go-native live-coding music environment**. Patterns are pure functions `State -> [Hap]` queried over rational time; output drives offline audio, MIDI, OSC, and the live console.

- **Module**: `codeberg.org/uzu/saint-hubbins`, Go 1.25
- **CLI**: `saint-hubbins` (alias `hubbins`) — `eval`, `serve`, `render`, `query`
- **Live console**: Go HTTP server at `http://localhost:8080` (`web/server.go`, page included as an inline Go template), talking to itself over HTTP — see below
- **License**: AGPL-3.0-or-later

## Goals

> Build a complete Go pattern engine with idiomatic replacements where browser APIs do not exist.

Not a port of any single JS framework — Saint Hubbins reimagines the pattern model in Go with its own vocabulary (mini notation, controls, transforms) and its own console.

## Two Execution Targets

1. **Native Go binary** — CLI `saint-hubbins` evaluates patterns, renders audio offline, and serves the live console over HTTP. The console's own JavaScript calls `POST /api/evaluate` — it does not load WASM.
2. **WASM build** — the pattern engine also compiles to WASM (`GOOS=js GOARCH=wasm`, `make wasm`) as a target for embedding the engine in someone else's page; Go serves the artifact at `/static/saint-hubbins.wasm`. This is a deliberate, separate target — see `docs/tutorial/08-limitations.md`.

## Constraints

- AGPL-3.0 compliance — `ATTRIBUTION.md` + `LICENSE`.
- No behavior drift in core query semantics without a migration note.
- No JS runtime embedded in production.

## See Also

- `ATTRIBUTION.md` — upstream inspiration credit (TidalCycles-inspired prototype / TidalCycles) retained in one place.
- `docs/archive/strudel-legacy/` — historical planning docs from the TidalCycles-inspired prototype, kept for reference only. Not part of the Saint Hubbins surface.
