# Saint Hubbins — conventions

Go live-coding pattern engine (`codeberg.org/uzu/saint-hubbins`). A `Pattern`
is a pure function of a time span: query it once over one cycle, or once over
a thousand, and it must return the same events either way. Most rules below
exist to protect that property, because the renderer only ever issues one
query per render — a bug that only shows up under a wide query reaches
listeners and nothing else.

## Cycle-dependent combinators must call `SplitQueries()`

If a combinator reads the cycle number off the query span — `state.Span.Begin.Sam()`
or `.Floor()` on it — and makes a decision from it (which branch to take, what
transform to apply), it **must** call `.SplitQueries()` on the returned
pattern, or iterate `state.Span.SpanCycles()` itself.

**Why:** a query spanning multiple cycles gives you one `state.Span` for the
whole range. Read the cycle from `Span.Begin` without splitting and you apply
one cycle's decision across the entire span — the combinator is only correct
by accident, when the query happens to be exactly one cycle wide.

**Why this bites in practice, not just in theory:** `internal/audio/webaudio.go:48`
renders with exactly one call, `pattern.QueryArc(0, cycles)`, over the whole
track. `Every` and `LastOf` originally read the cycle from the span start
without splitting. Every unit test passed — the suite queries cycle by cycle
— and every rendered track silently lost or duplicated its variations,
because the renderer's one query spans all of them at once. 1,200 tests, zero
of them wide enough to catch it.

Four non-test sites in `internal/core` do this today, and all four call it correctly:

| Reads the cycle at | Calls `SplitQueries()` at |
|---|---|
| `pattern_misc.go:151` (`LastOf`) | `pattern_misc.go:156` |
| `pattern_random.go:148` (randomness) | `pattern_random.go:158` |
| `pattern_time.go:193` (`Every`) | `pattern_time.go:198` |
| `pattern_vlpf_morph.go:159` (vlpf morph) | `pattern_vlpf_morph.go:174` |

Before adding or changing any combinator that touches the cycle number, run:

```
grep -rn "state\.Span\.Begin\|state\.Span\.End" internal/core/*.go | grep -v _test
```

(`.Sam()` on `Begin` is the usual spelling of "give me the cycle number," but
the search above is deliberately broader in two directions, and don't narrow
it back on the grounds that the extra generality looks unused. First: it
matches `state.Span.Begin` generally, not just `.Sam()` — the rule is "reads
the cycle number off the query span," not "spells it `.Sam()`," and a future
combinator could read `Begin` via `.Floor()` or direct arithmetic without
matching a `.Sam()`-only search. Second: `TimeSpan` has both `Begin` and
`End` (`internal/core/timespan.go:10-13`), and a combinator that reads
`state.Span.End` instead commits the identical bug — so the search covers
both fields. Do not loosen this to a bare `state.Span.` search either: that
matches 16 lines, most of them `SpanCycles()` calls — the compliant,
self-splitting pattern — and a search that flags correct code teaches people
to ignore it.)

Every hit must resolve to a `SplitQueries()` call (or its own `SpanCycles()`
loop). If your new code shows up in that grep without one, that's the bug.

**Test both query shapes.** A test that only queries one cycle at a time
cannot see this class of bug — it's exactly what let `Every`/`LastOf` ship
broken. When you add a cycle-dependent combinator, assert that one N-cycle
query produces the same haps as N separate one-cycle queries stitched
together.

## Only seven controls reach the offline audio

`internal/audio/webaudio.go`'s sine renderer reads exactly: `freq` (:90),
`n` (:106), `note` (:120), `s` (:132), `gain` (:156), `cutoff` (:176),
`lpf` (:189). The vocabulary holds **339** distinct control names —

```
grep -ohE 'createParam\("[^"]+"' internal/core/controls*.go | sort -u | wc -l
```

— and everything outside the seven above is carried through as data and
silently ignored by this renderer. Don't "fix" a control that appears to do
nothing to the WAV; check this list first.
`docs/tutorial/08-limitations.md` states the same seven for users;
`README.md` and `docs/tutorial/04-controls.md` quote the same 339, so adding
controls moves all three numbers. Two counts you will meet and should not
use: 477 exported `core.X` constructors (`Lpf`/`Cutoff` and `Sound`/`S` are
separate identifiers for one control), and the 295 at
`internal/core/controls.go:3` and `:79`, which is the *upstream JS* param
count and is stale.

There is no single priority order across the seven: `freq` (:90) only
short-circuits the three *pitch* fields — `n` (:106), `note` (:120), and `s`
(:132) are each guarded by `&& freq == 220.0`, so they're skipped once `freq`
has set a real value. `gain` (:156) and `cutoff`/`lpf` (:176/:189) carry no
such guard and are always read regardless of `freq`.

Two traps in that arrangement, both confirmed by rendering rather than by
reading:

- **`Freq(220)` is a no-op.** 220.0 is the renderer's default *and* the value
  the guard tests for, so a `freq` of exactly 220 leaves `n`/`note`/`s` in
  charge. `core.Note("c3").Set(core.Freq(220))` renders at c3 (523 zero
  crossings per cycle); with `Freq(1000)` it renders at 1 kHz (3999).
- **`lpf` is not a control name.** `core.Lpf` is an alias for `core.Cutoff`
  and sets `cutoff`. No *control* emits a bare `lpf` field, though
  `Pattern.Vlpf` (`pattern_vlpf_morph.go:29,33`) does — it always sets
  `cutoff` alongside, so the branch at :189 still never fires from engine
  code — only for a hand-built value map. It is seven field *lookups*, six
  of which a control can actually produce.

## Module layout

Engine packages live under `internal/`; Go only lets code inside this module
import them, so a song or tool that needs `internal/core` etc. must live in
this repository — you cannot `go get` this engine into a separate project.
Example/teaching programs go in `examples/`.

## No new third-party dependencies without discussion

`go.mod` has exactly one direct dependency (`github.com/dop251/goja`). A
single-binary build with one dependency is the project's selling point —
that's why the OSC and MIDI encoders here are hand-written instead of
imported. Adding a new direct dependency is a deliberate trade against that,
not a routine choice.

## License header

Every source file starts with:

```go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
```

Match it exactly (including the em dash) on new files. Generated files are
the exception — they carry their generator's banner instead, and you must not
add a header by hand to a file marked `DO NOT EDIT`.

## Generated files, and keeping them with their sources

Two checked-in files are generated, and neither is regenerated by any
automation. **There are no `go:generate` directives anywhere in this repo**,
so `make gen` (`go generate ./...`) is a no-op that exits 0 having done
nothing. Do not read a green `make gen` as "the generated files are current."

| Generated file | Source | Regenerate with |
|---|---|---|
| `internal/mini/parser.go` (`// Code generated by pigeon; DO NOT EDIT.`) | `internal/mini/parser.peg` | `pigeon -o internal/mini/parser.go internal/mini/parser.peg` |
| `internal/core/controls_gen.go` (`// Code generated by tools/gen-controls — DO NOT EDIT.`) | `js/packages/core/controls.mjs` — **not in this repo** | `go run ./tools/gen-controls` (a no-op today; see below) |

`pigeon` is not a module dependency; it is a standalone binary installed with
`go install github.com/mna/pigeon@latest`. Note the grammar: `parser.peg` is
the Go grammar and the real source. The sibling `internal/mini/krill.peg` is
the original pegjs grammar, kept for reference only — pigeon rejects it
(`Identifier reserved word`, JS actions), so it generates nothing.

`controls_gen.go` **cannot currently be regenerated.** `tools/gen-controls`
reads `js/packages/core/controls.mjs`; there is no `js/` directory in this
repository, so the tool prints "controls.mjs not found — skipping generation
(checked-in file is authoritative)" and exits 0 without writing a byte. The
file's `DO NOT EDIT` banner is therefore misleading: it is hand-maintained in
practice. If you must add a control, edit it directly and say so in the
commit message.

`web/static/saint-hubbins.wasm` and `wasm_exec.js` are build products of
`make wasm`, not source, and they are **git-ignored** (`.gitignore:11-12`).
Run `make wasm` before `serve`, or `/static/saint-hubbins.wasm` 404s.

**A generated file and its source must move in the same commit.** This
branch's own `parser.peg` bug was exactly this drift: a rebrand `sed` scoped
to `*.go` rewrote the import path inside the generated `parser.go` and never
touched the `.peg` it came from, so the two disagreed and only the artefact
was right. Fixing the `.peg` then left the artefact three bytes of source
offset behind, which took a second commit to clear — the drift outlived both
halves of the fix that caused it. Regenerating is cheap; verify with a diff
against a temp file before you commit:

```
pigeon -o /tmp/parser_check.go internal/mini/parser.peg
diff /tmp/parser_check.go internal/mini/parser.go
```

## Renaming? Search all text, not just `*.go`

The corollary of the above. When you rename a symbol, an import path or the
product name, `grep -rn` the whole tree — not `--include=*.go`, not `git
grep` restricted by extension. This repo's Go is partly generated from a
`.peg` grammar, the console's HTML lives inside a Go string literal in
`web/server.go`, and build commands are spelled out in the `Makefile` and in
`README.md`. A rename scoped to one extension will look complete and be
wrong, which is precisely how the `parser.peg` import path survived the
rebrand. `./scripts/check.sh`'s rebrand gate catches the specific case of a
leftover Strudel dependency or import path, but it only knows about those
three markers — it is not a general rename check.

## `gofmt`: format only what you touch

`gofmt -l .` reports **850** unclean files — 20 non-test files under
`internal/core` and ~830 tests. Do not run a blanket `gofmt -w` across the tree — it mixes unrelated cosmetic
churn into your diff and makes real changes harder to review. Format the
files you actually edit.

**`make fmt` is `gofmt -w .`.** It is right there in the `Makefile` next to
`make test` and `make lint`, and running it once rewrites all 850 into
your diff. The target stays — someone may deliberately want it — but don't
reach for it as the formatting step of ordinary work. `gofmt -w` the specific
files you changed instead.

This file is authoritative on conventions. `docs/04-conventions.md` covers
naming only and defers here; if the two ever disagree, this one is right and
the other is stale.

## `examples/` is teaching material

These programs are read as documentation, not just run. A comment that
misdescribes what its code does is a defect, the same as a bug — verify a
comment by actually running the pattern it describes, don't take it on faith
from a previous version of the file.

## Documented numbers are assertions, not decoration

`examples/examples_test.go` builds and runs all nine tutorial templates and
asserts each produces non-silent, non-clipping audio of the expected length.
`docs/tutorial/templates/*.md` additionally quotes each template's exact
event count and peak level in prose. Changing engine behavior that affects
any template means re-running it and updating the documented numbers — a
stale number in the docs is as wrong as a stale number in a test.

## Two lessons about tests that look fine

1. **A test can name a property it cannot defend.** Several tests here
   asserted an event's *presence* (`bytes.Contains`) where the actual claim
   was about *order* — and would have passed against the exact bug they were
   meant to catch. When you write a test for a property, break the property
   on purpose and confirm the test actually fails. An assertion nobody has
   watched fail is just slower than no assertion.
2. **Verify with the case that could fail, not the case that confirms.** A
   note-parser comparison passed across eight inputs because every one
   happened to be octave-qualified — bare note names were silently dropped
   and nothing noticed. A sort-stability fixture never reproduced because it
   used two elements, well under Go's insertion-sort threshold of twelve.
   Pick inputs that could expose the bug, not inputs that happen to agree
   with the code.

## Before claiming work is complete

Run `./scripts/check.sh` — nine gates: vet, the race-enabled test suite, the
wasm build, content-asserting checks on eval/render/midi/play, all nine
templates, and the rebrand check. A passing exit code alone is not the bar;
the script inspects actual output for each of those, not just exit status,
and the two gates that delegate to `go test -run` grep for the named test's
own `--- PASS:` line — `go test -run` exits 0 with "no tests to run" when the
pattern matches nothing.
