# 02 — Target Architecture (Go)

## Module

```
module codeberg.org/uzu/strudel-go   // or github.com/<you>/strudel-go
go 1.23
```

## Directory Layout

```
Go/
  go.mod
  go.sum
  cmd/
    strudel/          # native CLI: `strudel eval`, `strudel serve`, `strudel render`
      main.go
    strudel-wasm/     # WASM entry: exports Pattern API to JS
      main.go         # //go:build js && wasm
      exports.go
  internal/
    core/             # ← @strudel/core
      fraction.go / fraction_test.go
      timespan.go
      hap.go
      state.go
      pattern.go      # 4k LOC equivalent — split into pattern_core.go, pattern_time.go, pattern_structure.go, pattern_random.go
      controls.go     # generated from controls.mjs
      controls_gen.go # go:generate
      signal.go
      cyclist.go
      zyklus.go
      scheduler.go    # clockworker + schedulerState
      evaluate.go     # evalScope / evaluate
      util.go
      euclid.go
      pick.go
    mini/             # ← @strudel/mini
      parser.go       # pigeon-generated from krill.pegjs
      mini.go         # patternifyAST
      krill.peg
    transpiler/       # ← @strudel/transpiler
      transpiler.go   # Go parser (or goja shim phase 1)
      plugins.go      # plugin-mini, plugin-sample, etc.
    audio/            # ← superdough + supradough + webaudio + dough
      context.go      # AudioContext interface
      context_native.go  // oto/malgo impl  //go:build !js && !wasm
      context_wasm.go    // js bridge      //go:build js && wasm
      superdough.go
      sampler.go
      synth.go
      effects.go      # filter, reverb, delay, distortion, vowel
      wavetable.go
      worklets.go
      nodepool.go
      webaudio.go     # webaudioOutput, renderPatternAudio
    soundfonts/       # ← @strudel/soundfonts
      loader.go
      gm.go           # generated
      list.go
    tonal/            # ← @strudel/tonal
      tonal.go
      ireal.go
      voicings.go
    edo/              # ← @strudel/edo
    xen/              # ← @strudel/xen
    draw/             # ← @strudel/draw
      pianoroll.go
      spiral.go
      draw.go
    io/               # ← midi, osc, serial, mqtt, gamepad, motion, hydra, csound
      midi/
      osc/
      serial/
      mqtt/
      gamepad/
      hydra/
    repl/             # ← @strudel/repl
      repl.go
      prebake.go
    web/              # ← @strudel/web umbrella
      web.go
    reference/        # ← @strudel/reference
  web/                # ← website/ replacement
    templates/        # Go html/template or templ
      layouts/
      pages/
      repl/           # REPL HTML
    static/
      js/
      css/
      strudel.wasm    # built from cmd/strudel-wasm
      wasm_exec.js
    server.go         # http.ServeMux, handlers, static, API
    handlers_repl.go
  samples/            # ← samples/
  examples/           # ← examples/ ported to Go examples
  tools/
    gen-controls/     # generates controls_gen.go from JS controls.mjs
    gen-gm/           # generates soundfonts/gm.go
```

## Key Design Decisions

### 1. Pattern Engine: Generics vs `any`

JS Pattern is `Pattern<T>` in spirit but untyped. Go choice:

- `type Pattern struct { query func(State) []Hap; steps *Fraction }`
- `type Hap struct { Whole *TimeSpan; Part TimeSpan; Value any; Context map[string]any }`

Use `any` for `Value` to match JS bag semantics (`{s: "bd", gain: 0.5}`). Provide typed helpers `ValueAsMap()`, `ValueAsFloat()`. Generics (`Pattern[V]`) are optional for callers but not required for correctness; start with `any` to avoid premature abstraction.

`withValue(func(any) any)` maps to `func MapValue(fn func(any) any) Pattern`.

### 2. Fraction: Exact Rational

JS `fraction.js` is exact. Go: do **not** use `float64`. Options:
- Wrap `math/big.Rat` (exact, but heavy)
- Custom `struct { Num, Den int64 }` with normalization (faster, matches JS int range)
- Recommendation: custom `Fraction` with `int64` + `big.Int` fallback for overflow, plus `MarshalJSON`.

Implement `Add/Sub/Mul/Div/Mod/Cmp/LCM/GCD/ToFloat`.

### 3. Audio Abstraction

```go
type AudioContext interface {
  CurrentTime() float64
  SampleRate() int
  CreateGain() GainNode
  CreateBufferSource() BufferSource
  // ...
  Destination() AudioNode
}
```

- Native: `oto` (https://github.com/ebitengine/oto) or `malgo` (miniaudio) + manual DSP.
- WASM: `syscall/js` bridges to `window.AudioContext`.

This keeps `superdough` logic testable without a real audio device.

### 4. Transpiler: Two-Phase

- **Phase 1 (bootstrap)**: Use `goja` (pure Go JS runtime) to run the existing `transpiler.mjs` inside Go tests. Lets core tests pass while Go transpiler is built.
- **Phase 2**: Port transpiler to Go using `github.com/robertkrimen/otto` or hand-rolled parser with `acorn` equivalent (`github.com/dop251/goja/parser`). Final target: `pigeon` or `participle/v2` for mini-like template handling.

### 5. Mini Parser

`krill.pegjs` → Go: use `github.com/mna/pigeon` (PEG). Check in both `.peg` and generated `.go`. Keep JS `krill.pegjs` as source of truth until Go PEG is proven.

### 6. Frontend Strategy

- **Do not port React/Astro to Go component-for-component**. Replace with:
  - `net/http` + `html/template` (or `templ`) for docs/marketing pages.
  - REPL editor: **keep CodeMirror in JS** (it's the best editor). Load it as static JS; have it call into WASM (`strudel.wasm`) for pattern evaluation. Alternative: `monaco` via JS interop — same pattern.
  - Styling: reuse Tailwind CSS (static build), not Go-generated CSS.

### 7. Desktop (Tauri → Wails/Fyne)

- JS uses Tauri (Rust). Go equivalent: `wails` (Go + WebView) or `fyne` (pure Go). Recommend `wails` v2 to reuse existing `web/` frontend with minimal change. Bridges: MIDI (`gitlab.com/gomidi/midi`), OSC (`github.com/hypebeast/go-osc`), filesystem, clipboard.

### 8. Concurrency

- JS `zyklus` uses `setTimeout`/`AudioWorklet`. Go uses `time.Ticker` + goroutine scheduler.
- `Cyclist` → `type Scheduler struct { cps float64; mu sync.Mutex; tick *time.Ticker }`
- Pattern queries are pure and goroutine-safe; scheduler fans out Haps to outputs via channels.

## Build & CI

- `go.mod` replaces `pnpm-workspace.yaml` + `lerna.json`.
- `Makefile` targets: `make test`, `make lint`, `make wasm`, `make serve`, `make gen` (codegen).
- WASM: `GOOS=js GOARCH=wasm go build -o web/static/strudel.wasm ./cmd/strudel-wasm && cp $(go env GOROOT)/misc/wasm/wasm_exec.js web/static/`.
- Lint: `golangci-lint run` (replaces `eslint`).
- Format: `gofmt` + `goimports` (replaces `prettier`).
- Docs: `pkgsite` or `gomarkdoc` (replaces `jsdoc`).

