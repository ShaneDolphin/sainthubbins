# Roadmap — Saint Hubbins

This file is kept honest by running its gates, not by asserting them. Every
gate below was executed against this branch before being written down; see
`.superpowers/sdd/2026-08-24-repo-hygiene/task-1-report.md` for the exact
commands and their output.

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

Historical roadmap that referenced Strudel JS is archived at
`docs/archive/strudel-legacy/03-roadmap.md`.

## Remaining work

Two tracks are genuinely open. (Three tracks named in earlier drafts of the
hygiene plan — real-time OSC output, the pattern-engine gaps, and MIDI file
export — have already shipped and are folded into Phase 4/5 above; do not
re-open them.)

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
building the target, stop advertising it in the console.** The console
footer and package comment in `web/server.go` no longer claim a WASM
bridge; `docs/tutorial/08-limitations.md` states plainly that the console
uses HTTP by choice. Wiring WASM into the console (loading `wasm_exec.js`
and calling `saintHubbins.queryPattern` from the browser) was rejected: it
is real front-end work outside this plan and would create a second
evaluation path that must stay in sync with the Go one. Deleting the
target was also rejected: it costs nothing to keep and preserves embedding
the engine in a page as a future option. Revisit only if the console grows
enough that HTTP round-trip latency becomes a real problem — that is the
condition that would change this answer.

### Repo hygiene

This file has just been rewritten to remove the vacuous gate described
above. `docs/05-execution-checklist.md`, `scripts/check.sh`, and `CLAUDE.md`
are still outstanding — a later task in the same plan.

- **Gate (not yet passing):** `./scripts/check.sh` exits 0.
- **Plan:** `docs/superpowers/plans/2026-08-24-repo-hygiene.md`
