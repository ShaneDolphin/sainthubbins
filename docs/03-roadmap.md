# Roadmap — Saint Hubbins

This file is kept honest by running its gates, not by asserting them. Every
gate below was executed against this branch before being written down, and
`./scripts/check.sh` re-runs the automatable ones on demand — run it rather
than taking this page's word for anything.

## Phases (shipped)

| Phase | Focus | Key gate | Status |
|---|---|---|---|
| 0 | Foundations (`Fraction`, `TimeSpan`, `Hap`, `State`) | `go test ./internal/core/... -run Fraction\|TimeSpan` green | Done |
| 1 | Core engine (`Pattern`, controls, scheduler) | `go vet ./...` clean; `go test ./internal/core/...` green | Done |
| 2 | Mini-notation + transpiler | `mini.Mini("bd sd")` parses; `go test ./internal/mini/... ./internal/transpiler/...` green | Done |
| 3 | Offline audio (mono sine-oscillator WAV render) | `saint-hubbins render out.wav "bd sd"` writes a WAV `file(1)` identifies as RIFF/WAVE PCM | Done |
| 4 | I/O — real-time OSC and MIDI file export shipped; Serial/MQTT/Gamepad/motion remain no-op stubs | `saint-hubbins play` round-trips a real OSC bundle over UDP (`TestRunPlaySendsToSuperDirt` in `cmd/saint-hubbins`); `saint-hubbins midi out.mid "c3 e3 g3"` writes a file `file(1)` identifies as Standard MIDI data; `go test ./internal/osc/... ./internal/io/...` green | Done |
| 5 | Live console + WASM build | `saint-hubbins serve` answers `POST /api/evaluate` for mini-notation; `make wasm` builds `web/static/saint-hubbins.wasm` | Done, with a caveat (below) |

**Phase 5 caveat, verified against a running server:** the console evaluates
*mini-notation only*. `POST /api/evaluate {"code":"bd sd"}` returns two haps
(`bd`, `sd`). `POST /api/evaluate {"code":"s(\"bd sd\")"}` returns **one**
hap whose value is the literal string `s("bd sd")` — there is no text
evaluator for function-call syntax yet (tracked below). The WASM binary
builds cleanly but nothing loads it in the browser; the console talks to the
Go server over HTTP instead. Both are documented in
`docs/tutorial/08-limitations.md`.

**Phase 2 note:** the five mini-notation/pattern-engine gaps that used to be
tracked separately (`SlowCat`'s collapsed cycle span, `@` elongation, `!`
replicate, `%n` polymeter, and `Add` on control bags) are fixed. The
tutorial's former "Mini-notation differences from Strudel" table is gone —
`grep -rn "Mini-notation differences" docs/tutorial/` returns nothing.

**Repo hygiene, shipped:** this file, `docs/05-execution-checklist.md`,
`scripts/check.sh` and `CLAUDE.md` were rewritten together to remove the
vacuous gate described under "Text evaluator" below and to give the repo one
runnable entry point instead of a checklist a human ticks. `./scripts/check.sh`
exits 0 with all nine gates passing, and each gate asserts on output content
rather than exit status — see the script's comments for the specific bug each
one was built to catch. Plan:
`docs/superpowers/plans/2026-08-24-repo-hygiene.md`.

Historical roadmap that referenced Strudel JS is archived at
`docs/archive/strudel-legacy/03-roadmap.md`.

## Remaining work

One track is genuinely open; the WASM question below is settled and recorded
here so it is not asked again. (Three tracks named in earlier drafts of the
hygiene plan — real-time OSC output, the pattern-engine gaps, and MIDI file
export — have already shipped and are folded into Phase 4/5 above; repo
hygiene has shipped too, see the note above it. Do not re-open them.)

### Text evaluator

`saint-hubbins eval` and the web console's `/api/evaluate` parse bare
mini-notation (`"bd sd"` → 2 haps) but do not implement JS-like function-call
syntax (`s("bd sd")`, `.fast(2)`, `.gain(0.5)`) — that input comes back as a
single hap whose value is the literal source string. This is the exact gap
the old roadmap's Phase 5 gate obscured by testing `s("bd sd")` against an
evaluator that doesn't exist.

- **Gate (not yet passing):** `saint-hubbins eval 's("bd sd")'` returns two
  haps carrying `s` controls (`map[s:bd]`, `map[s:sd]`), not one hap whose
  value is the literal string `s("bd sd")`.
- **Plan:** `docs/superpowers/plans/2026-08-24-text-evaluator.md`

### WASM build target — decided, do not re-litigate

`make wasm` builds `web/static/saint-hubbins.wasm`, but the live console
never loaded it — it calls `/api/evaluate` over HTTP. Decision: **keep
building the target, stop advertising it as a bridge.** The console footer
and package comment in `web/server.go`, the `README.md` feature list and
project layout, and the package comment in `cmd/saint-hubbins-wasm` no longer
claim one; `docs/tutorial/08-limitations.md` states plainly that the console
uses HTTP by choice. Note that `saintHubbins.queryPattern` is a stub even for
an embedder — it echoes the code it is handed and returns an empty `haps`
array, without calling the pattern engine. Wiring WASM into the console
(loading `wasm_exec.js` and implementing `saintHubbins.queryPattern` for
real) was rejected: it is front-end work outside this plan and would create a
second evaluation path that must stay in sync with the Go one. Deleting the
target was also rejected: it costs nothing to keep and preserves embedding
the engine in a page as a future option. Revisit only if the console grows
enough that HTTP round-trip latency becomes a real problem — that is the
condition that would change this answer.
