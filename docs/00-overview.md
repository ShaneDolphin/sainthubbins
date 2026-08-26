# Saint Hubbins — Overview & Scope

## What Saint Hubbins Is

Saint Hubbins is a **Go-native live-coding music environment**. Patterns are pure functions `State -> [Hap]` queried over rational time; output drives offline audio, MIDI, OSC, and the live console.

- **Module**: `codeberg.org/uzu/saint-hubbins`, Go 1.25
- **CLI**: `saint-hubbins` (alias `hubbins`) — `eval`, `serve`, `render`, `play`, `query`, `midi` (full table in [README](../README.md#cli-reference))
- **Live console**: Go HTTP server at `http://localhost:8080` (`web/server.go`, page included as an inline Go template), talking to itself over HTTP — see below
- **Text input**: every command that takes `<code>` resolves it through `internal/jsapi.EvaluateCode` — JS pattern code first, mini-notation as the fallback, an error when it is neither
- **License**: AGPL-3.0-or-later

## Goals

> Build a complete Go pattern engine with idiomatic replacements where browser APIs do not exist.

Not a port of any single JS framework — Saint Hubbins reimagines the pattern model in Go with its own vocabulary (mini notation, controls, transforms) and its own console.

## Two Execution Targets

1. **Native Go binary** — CLI `saint-hubbins` evaluates patterns, renders audio offline, exports Standard MIDI Files, streams OSC to SuperDirt, and serves the live console over HTTP. The console's own JavaScript calls `POST /api/evaluate` — it does not load WASM.
2. **WASM build** — the pattern engine also compiles to WASM (`GOOS=js GOARCH=wasm`, `make wasm`) as a target for embedding the engine in someone else's page; Go serves the artifact at `/static/saint-hubbins.wasm`. This is a deliberate, separate target — see `docs/tutorial/08-limitations.md`.

## Constraints

- AGPL-3.0 compliance — `ATTRIBUTION.md` + `LICENSE`.
- No behavior drift in core query semantics without a migration note.
- **One** embedded JS runtime, and only as a front end for user input.
  `github.com/dop251/goja` is the module's single direct dependency; it backs
  `internal/jsapi`, which parses what a user types and hands back a
  `core.Pattern`. No engine internals run in JS, and nothing in the query path
  depends on it. This constraint used to read "no JS runtime embedded in
  production", which was true before the text evaluator shipped; the point it
  was protecting — the pattern engine stays pure Go — still holds.

## See Also

- `ATTRIBUTION.md` — upstream inspiration credit (TidalCycles-inspired prototype / TidalCycles) retained in one place.
- `docs/archive/strudel-legacy/` — historical planning docs from the TidalCycles-inspired prototype, kept for reference only. Not part of the Saint Hubbins surface.
