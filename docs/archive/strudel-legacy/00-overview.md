# 00 — Overview, Scope & Constraints

## What Strudel Is

Strudel (https://strudel.cc) is a **live-coding music environment — a JavaScript port of TidalCycles**. Users write pattern code in a browser REPL; patterns are pure functions `State -> [Hap]` queried over time; output drives Web Audio / MIDI / OSC / SuperDirt.

- **Monorepo**: `pnpm` + `lerna`, 34 packages under `packages/*`, plus `website/` (Astro + React + Tailwind), `src-tauri/` (Rust desktop), `examples/`, `test/`.
- **Size**: ~46k LOC JS/MJS + ~750 files; core pattern engine ~11k LOC.
- **License**: AGPL-3.0-or-later. Go port must remain AGPL-3.0.

## Goal of the Go Port

> **Factor the *entirety* of the JS application into idiomatic Go**, preserving user-visible behavior, with Go-idiomatic replacements where browser APIs do not exist.

Not a partial port. Every package listed in `01-deconstruction.md` must have a Go equivalent or an explicit, justified exclusion.

## Two Execution Targets (both required)

1. **Native Go binary** — CLI `strudel` that can evaluate patterns, render audio offline, and serve the REPL over HTTP. Uses Go audio backends (no browser).
2. **WASM + Go server** — Pattern engine compiled to WASM (`GOOS=js GOARCH=wasm`) so the existing browser REPL can call Go logic. Go HTTP server serves the frontend (`web/`). This is the path that achieves visual/audible parity in the browser.

Why both: Strudel's scheduling core (`clock`, `zyklus`, `schedulerState`) is timing-sensitive and benefits from Go goroutines natively, but the established UX is browser-based. Shipping only one target would leave half the product unreachable.

## Hard Constraints

- **AGPL-3.0 compliance** — keep license headers; do not relicense.
- **No behavior drift** — pattern query semantics (haskell-style `fmap`/`applicative`/`bind`, `TimeSpan`, `Fraction`, `Hap` whole/part) must be bit-identical to JS for the test suite (`packages/core/test/*.mjs`, `test/tunes.test.mjs`, `test/examples.test.mjs`).
- **Audio fidelity** — `superdough` DSP must produce equivalent output (within float tolerance) or the port is not credible.
- **No JS runtime embedded in Go production** — the point is Go; do not shell out to Node for core logic. Using `otto`/`goja` as a stopgap for the transpiler is acceptable only if flagged and replaced.

## Non-Goals (explicitly out of scope)

- Reimplementing the entire npm ecosystem (CodeMirror, React) in Go — replace with Go-native equivalents (see architecture doc).
- Pixel-perfect clone of the Astro website — functional parity is required, not CSS parity.

## Success Criteria (gates)

Each phase has an exit gate in `03-roadmap.md`. Global gates:

1. `go test ./...` passes for ported core (mirrors `vitest` suite).
2. `go vet` + `golangci-lint` clean.
3. WASM build succeeds: `GOOS=js GOARCH=wasm go build -o web/static/strudel.wasm ./cmd/strudel-wasm`.
4. `go run ./cmd/strudel serve` boots REPL at `http://localhost:8080` with pattern evaluation + audio (native or via WASM bridge).
5. Example tunes from `test/testtunes.mjs` render offline to WAV without error.

## Risks (see also `04-conventions.md`)

| Risk | Mitigation |
|------|------------|
| JS dynamic dispatch → Go static types is lossy | Generics + `any` + codegen for `controls.mjs` (295 exports) |
| Web Audio API has no Go equivalent | Abstract `AudioContext` interface; native impl via `oto`/`malgo`, WASM impl delegates to browser |
| Floating rational timing (`fraction.js`) | Port `Fraction` exactly; use `big.Rat` or custom struct, not `float64` |
| Mini PEG parser (`krill-parser.js` 2497 LOC generated) | Regenerate from `krill.pegjs` via `pigeon` or hand-port; keep PEG source as truth |

