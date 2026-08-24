# 03 — Phased Roadmap & Gates

> Execute in order. Each phase lists inputs, outputs, and a **gate** that must pass before proceeding. Phases 0-2 are strictly sequential; 3-4 can overlap after 2; 5 requires 0-3.

## Timeline (conservative, 1 engineer)

| Phase | Name | Duration | Depends |
|-------|------|----------|---------|
| 0 | Foundations (Fraction/TimeSpan/Hap/State + repo) | 1-2 wks | — |
| 1 | Core Engine (Pattern + Controls + Signals + Scheduler) | 3-4 wks | 0 |
| 2 | Mini + Transpiler | 2-3 wks | 1 |
| 3 | Audio (superdough/webaudio/soundfonts) | 4-6 wks | 1 |
| 4 | I/O (MIDI/OSC/Serial/MQTT/Gamepad/Hydra/CSound) | 2-3 wks | 1 |
| 5 | Frontend + REPL + Desktop | 4-6 wks | 2, 3 |
| **Total** | | **~16-24 wks** | |

## Phase 0 — Foundations

- **Input**: `js/packages/core/fraction.mjs`, `timespan.mjs`, `hap.mjs`, `state.mjs`, `util.mjs` (subset), `fraction.js` dep.
- **Output**: `internal/core/fraction.go`, `timespan.go`, `hap.go`, `state.go`, `util.go` + tests mirroring `packages/core/test/*`.
- **Gate**: `go test ./internal/core -run TestFraction` etc. All rational arithmetic and span ops match JS fixtures.

## Phase 1 — Core Engine

- **Input**: `pattern.mjs` (4191 LOC), `controls.mjs` (3319 LOC, 295 params), `signal.mjs`, `cyclist.mjs`, `zyklus.mjs`, `evaluate.mjs`, `schedulerState.mjs`, `euclid.mjs`, `pick.mjs`.
- **Output**: `internal/core/pattern*.go`, `controls.go` (generated), `signal.go`, `scheduler.go`, `evaluate.go`.
- **Gate**: Ported `packages/core/test/pattern.test.mjs` (1394 LOC) to Go; `go test ./internal/core` passes. Pattern queries for `stack`/`cat`/`fast`/`slow`/`euclid`/`rand` produce identical Hap sets to JS (snapshot comparison via `goja` bridge in tests).

## Phase 2 — Mini + Transpiler

- **Input**: `packages/mini/krill.pegjs` + `mini.mjs`, `packages/transpiler/transpiler.mjs` + 4 plugins.
- **Output**: `internal/mini/parser.go` (pigeon), `mini.go`, `internal/transpiler/transpiler.go`.
- **Phase 1 interim**: use `goja` to run JS transpiler inside Go for test parity.
- **Gate**: `test/tunes.test.mjs` and `packages/mini/test/*.mjs` ported; mini strings like `"bd ~ sd ~"` parse to identical Patterns in Go vs JS. Transpiled JS snippets evaluate equivalently.

## Phase 3 — Audio

- **Input**: `packages/superdough/*`, `packages/supradough/*`, `packages/webaudio/*`, `packages/soundfonts/*`, `packages/dough/*`.
- **Output**: `internal/audio/*`, `internal/soundfonts/*`, `internal/audio/webaudio.go` (OfflineAudioContext render).
- **Gate**: `renderPatternAudio` equivalent in Go renders 4-cycle test pattern to WAV; byte-compare within float epsilon against JS offline render. `go run ./cmd/strudel render --cycles 2 "s(\"bd sd\")"` produces valid WAV.

## Phase 4 — I/O

- **Input**: `packages/midi/*`, `osc/*`, `serial/*`, `mqtt/*`, `gamepad/*`, `motion/*`, `hydra/*`, `csound/*`, `desktopbridge/*`.
- **Output**: `internal/io/{midi,osc,serial,mqtt,gamepad,hydra}/` each with `Interface` + `native`/`wasm` builds + `mock` for tests.
- **Gate**: Each I/O package has a mock that records Haps; integration test shows `note("c3 e3").midi()` emits correct MIDI bytes.

## Phase 5 — Frontend + REPL + Desktop

- **Input**: `packages/codemirror/*`, `draw/*`, `repl/*`, `website/*`, `src-tauri/*`, `examples/*`.
- **Output**: `web/server.go` + `web/templates/*` + `web/static/*` (CodeMirror JS + strudel.wasm), `internal/draw/*` (SVG/JSON), `internal/repl/*`, `cmd/strudel serve`.
- **Gate**: `go run ./cmd/strudel serve` → `http://localhost:8080` shows REPL, evaluates `sound("bd sd")` audibly (WASM audio path), displays pianoroll. `examples/headless-repl` equivalent works as `go run ./examples/headless/main.go`.

## Cross-Cutting Gates (apply to every phase)

- [ ] `go vet ./...` clean
- [ ] `golangci-lint run` clean (or documented suppressions)
- [ ] `go test ./... -race` clean
- [ ] No `// TODO(port)` left in shipped code without tracking issue
- [ ] AGPL-3.0 header on every new `.go` file

## Suggested Milestones (for project tracking)

```
M0.1  go.mod + Fraction + TimeSpan + Hap + State merged
M1.1  Pattern query core merged (fmap/bind/cat/stack)
M1.2  Controls (295 params) generated + tested
M1.3  Scheduler (zyklus/cyclist) + evaluate merged
M2.1  Mini PEG parser merged + mini.test.go green
M2.2  Transpiler (goja shim) merged; tunes.test.go green
M3.1  AudioContext interface + superdough sampler/synth merged
M3.2  Offline WAV render green
M4.1  MIDI + OSC merged
M5.1  web/server.go + WASM build merged
M5.2  REPL end-to-end demo + release tag v0.1.0-go
```

## De-Risking Order

If time is constrained, prioritize **0 → 1 → 2 → 3 (offline render only) → 5 (minimal REPL without I/O)**. Phase 4 I/O can ship after 5 without blocking the core demo.

