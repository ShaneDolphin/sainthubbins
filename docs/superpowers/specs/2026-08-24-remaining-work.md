# Saint Hubbins — Remaining Work (audit, 2026-08-24)

Findings from auditing the codebase against `docs/03-roadmap.md`. Every claim
here was verified by running code, not read off the source.

## Verified state

| Subsystem | State | Evidence |
|---|---|---|
| Pattern engine | Works | 14 packages pass `go test ./... -race` |
| Mini-notation | Works, three deviations | `docs/tutorial/08-limitations.md` |
| Offline WAV render | Works, mono sine only | `internal/audio/webaudio.go` |
| `Clock` (lookahead scheduler) | **Implemented, unused** | `internal/core/scheduler.go:51`; no caller outside its own file |
| `Cyclist` (hap scheduler) | **Implemented, never started** | `internal/core/scheduler.go:130+`; `internal/session/session.go:21` builds it with `OnTrigger = nil` |
| OSC output | **No-op stub** | `internal/osc/osc.go` — `SendSuperDirt` formats a string, discards it, returns nil |
| MIDI output | **Interface + mock only** | `internal/io/midi.go` — `MockMIDI` appends to `[]string` |
| Serial / MQTT / gamepad / motion / csound / sampler / desktop | Stubs, 12–23 LOC each | no external dependencies in any |
| Transpiler | 269 LOC, **not wired** | not called from `cmd/`, `web/` or `core/`; `Evaluate()` gets `nil` everywhere |
| WASM | Builds, **never loaded** | console calls `/api/evaluate` over HTTP instead |
| `go vet ./...` | Clean | — |

## The core finding

The scheduler is finished and the output backends are empty. Three pieces exist
and none are connected:

```
Cyclist (works) --OnTrigger--> nil          <-- the missing link
     |
     +-- Clock (works, lookahead + latency)
     +-- Pattern.QueryArc (works, filters HasOnset already)
```

`Cyclist.Start` already queries the pattern per tick, filters to haps with an
onset, and computes `targetTime`, `duration` and `deadline` for each. It calls
`c.OnTrigger(hap, deadline, duration, cps, targetTime)` — and every construction
site passes `nil`.

Connecting that callback to a real socket is the highest-value work in the
project: it turns a pattern library into a live instrument, and reuses
SuperDirt's samples and stereo engine rather than growing the sine renderer.

## Confirmed engine defects

1. **`SlowCat` drops multi-event arguments.** `internal/core/pattern.go`,
   `SlowCat`: the per-cycle span is mapped with `t.Sub(t.Sam())`, which maps a
   span's *end* to zero — `[1,2)` becomes `0/1 → 0/1`, a zero-width span.
   Whole-cycle arguments survive (they answer a point query); anything with
   events at sub-cycle positions returns nothing. The correct base,
   `cyc.Begin.Sam()`, is already computed nine lines later for the reverse
   shift.
2. **`@` elongation is a no-op.** `internal/mini/mini.go` implements it with
   `WithSteps`, which does not affect timing. `TimeCatWeighted`
   (`internal/core/pattern_weighted.go:9`) is the right primitive.
3. **`!` replicate subdivides in place.** Implemented with `Ply`, so
   `"bd!3 sd"` puts three kicks inside the first half instead of producing four
   equal steps.
4. **`%n` polymeter suffix is ignored.** `internal/mini/mini.go` parses it and
   returns the base pattern unchanged.
5. **`Add` on a wrapped pattern replaces the control bag** with a bare number
   instead of transposing.

## Roadmap hygiene

`docs/03-roadmap.md` phase 5 and `docs/05-execution-checklist.md` both gate on
`saint-hubbins eval 's("bd sd")'` returning haps. No evaluator for that syntax
exists, and the gate passes vacuously today: the command returns one hap whose
value is the literal source text. There is no live status document
(`IMPLEMENTATION_STATUS.md` is archived) and no `CLAUDE.md` recording
conventions — the absence of the "cycle-dependent combinators must call
`SplitQueries`" rule caused a real bug found in review.

## Global constraints

- Go 1.25.0, module `codeberg.org/uzu/saint-hubbins`.
- **No new third-party dependencies.** The project's selling point is a single
  Go binary with no runtime deps; OSC encoding is ~150 lines of pure Go.
- AGPL-3.0-or-later header on every new file, matching existing files.
- Engine packages live under `internal/`; anything importing them must be in
  this module.
- Tests must be hermetic — no hardware, no network peers outside loopback.
- `go test ./... -race -count=1` and `go vet ./...` must stay clean.
