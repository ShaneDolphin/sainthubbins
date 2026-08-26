```
 ____      _     ___  _   _  _____
/ ___|    / \   |_ _|| \ | ||_   _|
\___ \   / _ \   | | |  \| |  | |
 ___) | / ___ \  | | | |\  |  | |
|____/ /_/   \_\|___||_| \_|  |_|
 _   _  _   _  ____   ____   ___  _   _  ____
| | | || | | || __ ) | __ ) |_ _|| \ | |/ ___|
| |_| || | | ||  _ \ |  _ \  | | |  \| |\___ \
|  _  || |_| || |_) || |_) | | | | |\  | ___) |
|_| |_| \___/ |____/ |____/ |___||_| \_||____/
```

# Saint Hubbins — Go-Native Live Coding Pattern Engine

Saint Hubbins is a Go-native environment for algorithmic music and live coding. Patterns are pure functions of time — `Pattern` values queried over rational time spans — that emit sound, MIDI, and visual events. The system includes a pattern engine, a compact text notation, a control vocabulary, offline audio rendering, and a live console served from a single Go binary.

Patterns are first-class values. You compose them, transform them in time, and layer them. The engine evaluates a pattern for any time window and returns the active events (`haps`) with precise timing, so playback and rendering stay deterministic. These go to eleven.

---

## New here? Start with the tutorial

**[docs/tutorial/](docs/tutorial/README.md)** takes you from a first sound to writing your own tracks, and ships **eight complete example tracks** you can run and edit:

```console
$ go run ./examples/house
house.wav — 160 events over 8 bars, 16.0s, peak 0.81
```

| Style | BPM | Run |
|-------|-----|-----|
| House | 125 | `go run ./examples/house` |
| Chicago House | 120 | `go run ./examples/chicago-house` |
| Techno | 132 | `go run ./examples/techno` |
| Minimal dubstep | 140 | `go run ./examples/minimal-dubstep` |
| Maximal dubstep | 140 | `go run ./examples/maximal-dubstep` |
| Drum and bass | 174 | `go run ./examples/drum-and-bass` |
| Electronica | 110 | `go run ./examples/electronica` |
| Trance | 138 | `go run ./examples/trance` |

Each has a [line-by-line walkthrough](docs/tutorial/templates/README.md) explaining what every line does and what to change to make it your own.

The tutorial also covers [building a track from scratch](docs/tutorial/06-new-song-cli.md), the [mini-notation grammar](docs/tutorial/02-mini-notation.md), and an honest list of [current limitations](docs/tutorial/08-limitations.md).

---

## Features

- **Pattern engine** — exact rational timing (`Fraction`), `TimeSpan`, `Hap`, `State`, and a `Pattern` type with functor / applicative / monadic composition
- **Mini notation** — compact string language for rhythms (`"bd sd ~"`, `"[bd sd]*2"`, `"c3 e3 g3"`, `"bd(3,8)"`)
- **339 controls** — sound selection, pitch, filters, envelopes, spatialization, buses, and synthesis params, all composable via `Set`
- **Transformation core** — time (`Slow`/`Fast`/`Early`/`Late`/`Compress`/`Zoom`), structure (`Stack`/`FastCat`/`SlowCat`/`Arrange`/`Palindrome`/`Jux`), Euclidean (`Euclid`/`Bjorklund`/`Struct`), repetition and slicing (`Ply`/`Chop`/`Striate`/`Segment`), conditional (`When`/`Every`/`Sometimes`/`Degrade`), and alignment variants
- **Live console** — `POST /api/evaluate` and `POST /api/pianoroll` plus a single-page editor, served by `go run ./cmd/saint-hubbins serve`
- **Offline audio** — mono `float32` rendering to WAV with gain, filter, and note-to-frequency mapping
- **Music theory** — scales, chords, voicings, and transposition
- **Visuals** — pianoroll, spiral, and pitch-wheel helpers — Stonehenge edition
- **WASM build target** — the engine compiles under `GOOS=js GOARCH=wasm` to `saint-hubbins.wasm`, kept as a starting point for embedding it in someone else's page. Nothing in this repository loads it — the live console calls `POST /api/evaluate` over HTTP instead — and its `saintHubbins.queryPattern` export is still a stub that echoes its argument and returns an empty `haps` array
- **I/O backends** — real: MIDI file export (`midi`) and real-time OSC to SuperDirt (`play`); still no-op stubs: Serial, MQTT, Gamepad, motion sensing

---

## Requirements

- Go 1.25 or newer
- No Node.js, no external audio server required for build or offline render

---

## Installation

```bash
git clone <this-repo>
cd Go
go mod download

# build the CLI
go build -o saint-hubbins ./cmd/saint-hubbins
# also available as hubbins (symlink)
ln -s saint-hubbins hubbins

# or run directly
go run ./cmd/saint-hubbins query
```

### Makefile helpers

```bash
make test        # go test ./... -race -count=1
make lint        # go vet ./...
make wasm        # GOOS=js GOARCH=wasm build -> web/static/saint-hubbins.wasm
make serve       # go run ./cmd/saint-hubbins serve
make gen         # go generate ./... — no-op; the repo has no go:generate directives
make fmt         # gofmt -w . — see below before you run it
```

`make fmt` reformats the whole tree, and **850 files are not currently
`gofmt`-clean** (20 of them non-test files under `internal/core`, the rest
tests), so one run drops all of them into your diff. Format the files you
actually edited instead.

---

## Quick Start

```bash
# 1. Synthetic query — stack two sound events over 2 cycles
go run ./cmd/saint-hubbins query

# 2. Evaluate a mini-notation pattern
go run ./cmd/saint-hubbins eval "bd sd"

# 3. Start the live console server (http://localhost:8080)
go run ./cmd/saint-hubbins serve
# or on a custom address:
go run ./cmd/saint-hubbins serve :3000

# 4. Render 4 cycles of a pattern to a WAV file
go run ./cmd/saint-hubbins render out.wav "bd sd hh cp"
```

Open `http://localhost:8080` after `serve` — the page contains an editor with **Evaluate** and **Hush** that POST to the server.

---

## CLI Reference

The binary is `saint-hubbins` (`./cmd/saint-hubbins`, alias `hubbins`). Six subcommands:

| Command | Usage | Effect |
|---|---|---|
| `query` | `saint-hubbins query` | Demo: `Stack(s("bd"), s("sd"))` queried over 0..2 cycles, pretty-printed JSON with `whole`, `part`, `value` |
| `eval` | `saint-hubbins eval <code>` | Parses `<code>` as mini or control expression, queries 0..1 cycle, prints JSON |
| `serve` | `saint-hubbins serve [addr]` | Starts the console server. Default `addr` is `:8080`. |
| `render` | `saint-hubbins render <out.wav> <code>` | Renders `<code>` for 4 cycles at 48 kHz and writes a 16-bit mono WAV |
| `play` | `saint-hubbins play <code> [host] [port] [secs]` | Streams `<code>` to SuperDirt over OSC. Defaults: `host` `127.0.0.1`, `port` `57120`, `secs` `8`. **Requires SuperCollider with SuperDirt already running and listening on port 57120** — if you hear nothing, that is almost always why. |
| `midi` | `saint-hubbins midi <out.mid> <code> [cycles]` | Renders `<code>` for `cycles` cycles (default 4) at 480 ticks per quarter note and writes a Standard MIDI File |

`eval` and `render` accept **mini-notation** — the rhythm language documented in
[the tutorial](docs/tutorial/02-mini-notation.md). Function-call syntax such as
`s("bd sd")` or `.fast(2)` is *not* implemented as text: there is no script
evaluator, so unrecognised input comes back as a literal string value. Layering,
controls and transformations are the [Go API](docs/tutorial/03-patterns-in-go.md).

`play "0 1 2 3"` sends four OSC messages naming **samples** `0`, `1`, `2`, `3`
(`s "0"`, `s "1"`, ...), not notes: mini-notation stores a bare numeric token
as a string, so it takes the same path as `bd`, not the numeric `n` path. If
you want a sample index, use `bd:3` (see `bd:1` in the mini-notation table
below); if you want actual note numbers over OSC, build the pattern with the
Go API's `core.N` instead of typing bare numbers into mini-notation.

`midi` has the same bare-numeric trap, and one MIDI-specific trap of its own:

- `midi out.mid "0 1 2 3"` writes a *valid but empty* MIDI file — no note
  events at all, exit code 0. The same mini-notation rule applies: a bare
  numeric token is a string, not a pitch, so every hap is skipped as
  pitchless. Unlike `play`'s `bd:3` sample-index workaround, MIDI export
  needs a resolvable note: reference a bare drum name like `bd` for a
  percussive hit, or use `core.N`/`core.Note` from the Go API for real note
  numbers. The CLI reports the emitted note count (`wrote out.mid (1 cycles,
  0 notes)`) and warns on stderr when it is zero, so this should not be
  silent.
- `midi out.mid "bd:3"` exports **note 3 on channel 0**, not a drum hit on
  channel 9 as `bd:3` in mini-notation table entries suggest. `HapToNote`
  resolves `n` before `s` (`bd:3` sets both `s=bd` and `n=3`), so the numeric
  index wins over the drum name. This is a known limitation, not a bug: the
  offline audio renderer (`internal/audio/webaudio.go`) has the identical
  `n`-before-`s` precedence, so MIDI export intentionally matches what
  `render`/`play` would sound like. If you want a drum note over MIDI,
  reference the drum name without a `:n` sample-index suffix (e.g. `bd`, not
  `bd:3`).

Examples:

```bash
saint-hubbins eval 'bd ~ sd cp'
saint-hubbins eval 'c3 e3 g3'
saint-hubbins eval '[bd*4, hh*8]'
saint-hubbins render /tmp/test.wav 'bd(3,8)'
```

---

## HTTP API

When `serve` is running, the server exposes:

| Method | Path | Body | Returns |
|---|---|---|---|
| `GET` | `/` | — | Console HTML (editor + JS `fetch` to `/api/evaluate`) |
| `GET` | `/health` | — | `ok` |
| `POST` | `/api/evaluate` | `{"code":"bd sd"}` | `{"haps":[{"whole":"0/1 → 1/2","part":"0/1 → 1/2","value":"bd"}, ...]}` |
| `POST` | `/api/pianoroll` | `{"code":"..."}` | `{"haps":[... with time/duration ...]}` queried over 0..2 cycles |
| `GET` | `/static/*` | — | Files under `web/static/`. `saint-hubbins.wasm` and `wasm_exec.js` are build products — run `make wasm` to generate them |

CORS headers (`Access-Control-Allow-Origin: *`) are set for API routes; `OPTIONS` preflight returns 204.

```bash
curl -s -X POST http://localhost:8080/api/evaluate \
  -H 'Content-Type: application/json' \
  -d '{"code":"bd sd"}' | jq .

curl -s -X POST http://localhost:8080/api/pianoroll \
  -H 'Content-Type: application/json' \
  -d '{"code":"bd sd cp"}' | jq .
```

---

## Web Console

- `web/server.go` — `Server` with `Handler()` and `Start()`; the console page is a Go template literal in that file (`consoleTemplate`), not a separate asset
- `web/static/` — `saint-hubbins.wasm` + `wasm_exec.js` produced by `make wasm`
- `cmd/saint-hubbins-wasm` — `//go:build js && wasm` entry exporting `saintHubbins.queryPattern(code)` and `version` on `js.Global()`. The console never loads it; it is here for embedders, and `queryPattern` is a stub (it returns `{code, length, haps: []}` without touching the pattern engine)

Embed the server in another Go program:

```go
import "codeberg.org/uzu/saint-hubbins/web"

srv := web.NewServer(":8080")
log.Fatal(srv.Start())
```

---

## Using the Engine as a Go Library

Core packages live under `internal/` (`core`, `mini`, `transpiler`, `audio`, `draw`, `tonal`, etc.). They are idiomatic Go and have no runtime dependency on a browser.

### Querying a pattern

```go
package main

import (
    "fmt"
    "codeberg.org/uzu/saint-hubbins/internal/core"
    "codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
    mini.RegisterStringParser() // enable string → Pattern via mini

    pat := core.Stack(core.S("bd"), core.S("sd"))
    haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(2))
    fmt.Println(len(haps), "haps")

    pat2 := mini.Mini("bd ~ sd cp")
    haps2 := pat2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
    fmt.Printf("%+v\n", haps2[0].Value)

    fast := pat2.FastF(core.FractionFromInt(2))
    _ = fast

    pat3 := core.S("bd").Set(core.Gain(0.8).Set(core.Pan(0.2)))
    hv := pat3.FirstCycle()[0].Value.(map[string]any)
    fmt.Println(hv["s"], hv["gain"], hv["pan"])
}
```

### Time model

```go
a := core.FractionFromInt(1).Div(core.FractionFromInt(2)) // 1/2
b := core.NewFraction(3, 4)                                 // 3/4
span := core.NewTimeSpan(core.FractionFromInt(0), core.FractionFromInt(1))
haps := pat.Query(core.NewState(span))
```

### Mini Notation

`mini.Mini(string) Pattern` parses the compact language. The full, verified grammar is in [the tutorial](docs/tutorial/02-mini-notation.md).

| Syntax | Meaning | Example |
|---|---|---|
| `bd sd cp` | Sequence in one cycle | `mini.Mini("bd sd")` |
| `~` | Rest / silence | `mini.Mini("bd ~ sd")` |
| `*n` / `/n` | Speed up / slow down token | `bd*2`, `bd/2` |
| `(p,s)` / `(p,s,r)` | Euclidean Bjorklund | `bd(3,8)`, `bd(3,8,2)` |
| `@n` | Elongate / weight | `bd@2` |
| `!` / `!n` | Replicate | `bd!2` |
| `?` / `?n` | Degrade chance | `bd?0.5` |
| `[a b]` | Subsequence (group) | `[bd sd]*2` |
| `<a b>` | Alternate each cycle | `<bd sd>` |
| `{a b, c d e}` | Polymeter (stack folded to `Stack` when steps unavailable) | `{a b, c d e}` |
| `a \| b \| c` | Choose one per cycle | `a | b` |
| `0 .. 4` | Range 0..4 inclusive | `0 .. 4` |
| `bd:1` | Sample index via control `n` | `bd:3` → `s=bd n=3` |
| `a,b` | Stack layers | `a,b` or `a , b` |

---

## Audio

```go
import "codeberg.org/uzu/saint-hubbins/internal/audio"

samples, err := audio.RenderPatternAudio(pat, 4, 48000)
err = audio.WriteWAV("out.wav", samples, 48000)
```

---

## Project Layout

```
Go/
  cmd/saint-hubbins/        # native CLI + console server entry (alias hubbins)
  cmd/saint-hubbins-wasm/   # //go:build js && wasm — embedding target, not loaded by the console
  web/
    server.go         # Server.Handler(), /api/evaluate, /api/pianoroll, inline console template
    static/           # saint-hubbins.wasm + wasm_exec.js (generated)
  internal/
    core/             # Fraction, TimeSpan, Hap, State, Pattern, controls, signals, scheduler
    mini/             # mini.Mini, parser (pigeon)
    transpiler/       # Transpile, EvaluateJS
    audio/            # OfflineRenderer, RenderPatternAudio, WriteWAV
    draw/             # Pianoroll, Spiral, pitch wheel, animation
    tonal/            # Scale / Chord / Voicing / Transpose
    session/          # live session (evaluation + scheduler)
  tools/gen-controls/ # generated controls_gen.go; cannot re-run (needs a js/ tree not shipped here)
  go.mod              # module codeberg.org/uzu/saint-hubbins, go 1.25, github.com/dop251/goja
  Makefile
  LICENSE             # AGPL-3.0-or-later
  ATTRIBUTION.md      # upstream credit
```

---

## Development

One command runs every automated gate — vet, the race-enabled suite, the WASM
build, `eval`/`render`/`midi`/`play`, all nine tutorial templates, and a
rebrand check. Each gate asserts on the *output* it gets, not just on a zero
exit status:

```bash
./scripts/check.sh
```

The individual steps, if you want one on its own:

```bash
go test ./... -race -count=1
go vet ./...
go test -tags goja ./...
GOOS=js GOARCH=wasm go build -o web/static/saint-hubbins.wasm ./cmd/saint-hubbins-wasm
go run ./tools/gen-controls   # no-op today — it needs a js/ tree this repo does not ship
```

---

## License

AGPL-3.0-or-later. See `LICENSE`.

Copyright (C) 2026 Saint Hubbins contributors.
