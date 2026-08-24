# Phase 5 — Frontend, REPL, Draw & Desktop

> Duration: 4-6 weeks. Depends on Phases 2 and 3. This is the user-visible finish line.

## Objective

Replace `website/` (Astro + React), `packages/codemirror/`, `packages/draw/`, `packages/repl/`, and `src-tauri/` with Go-served frontend + WASM pattern engine.

## 5.1 Draw / Visualizers (`packages/draw/` ~700 LOC)

| JS file | Go equivalent | Output |
|---------|---------------|--------|
| `pianoroll.mjs` 318 LOC | `internal/draw/pianoroll.go` | SVG string or JSON for frontend |
| `spiral.mjs` | `internal/draw/spiral.go` | SVG |
| `pitchwheel.mjs` | `internal/draw/pitchwheel.go` | SVG |
| `draw.mjs` / `animate.mjs` / `color.mjs` | `internal/draw/draw.go` | helpers |

- Go draw packages produce **data, not DOM** — return SVG strings or JSON (`PianorollData{Notes: []Note{...}}`) that the Go template or WASM JS renders.
- Test: `TestPianoroll` compares SVG output against JS fixture (snapshot).

## 5.2 REPL (`packages/repl/` + `packages/codemirror/`)

### Editor

- **Do not reimplement CodeMirror in Go.** Keep CodeMirror 6 as static JS under `web/static/js/codemirror.js` (bundled via `esbuild` or vendored).
- CodeMirror calls into WASM for evaluation: `strudelWasm.evaluate(code) -> {pattern, error}` via `wasm_exports.go`.

### WASM Bridge (`cmd/strudel-wasm/`)

```go
//go:build js && wasm
package main
import "syscall/js"
func main() {
    js.Global().Set("strudelGo", map[string]any{
        "evaluate": js.FuncOf(evaluateJS),
        "query":    js.FuncOf(queryJS),
    })
    select{}
}
func evaluateJS(this js.Value, args []js.Value) any { /* code -> Pattern -> Haps JSON */ }
```

- Build: `make wasm` → `web/static/strudel.wasm` + `wasm_exec.js`.
- JS loader: `web/static/js/wasm.js` loads WASM and exposes `window.strudelGo`.

### Go REPL Logic (`internal/repl/`)

```go
package repl
type Session struct { Code string; Pattern core.Pattern; Err error }
func (s *Session) Evaluate(code string) error // transpiles + evaluates
func (s *Session) Query(span core.TimeSpan) []core.Hap
func Prebake(code string) (string, error) // mirrors prebake.mjs
```

## 5.3 Website / Server (`website/` → `web/`)

### Pages to port (from `website/src/pages/` + `src/content/`)

- `/` landing, `/learn/*` docs, `/workshop/*`, `/blog/*`.
- Start with minimal: REPL at `/`, docs list at `/learn`, single doc page.

### Server (`web/server.go`)

```go
package web
func NewServer(addr string, opts Options) *http.Server
// Routes:
//   GET  /              → REPL template
//   GET  /learn/*       → docs
//   POST /api/evaluate  → {code} -> {haps, error} (native fallback if WASM not used)
//   POST /api/render    → {code, cycles} -> WAV bytes
//   GET  /static/*      → embedded FS
```

- Templates: `html/template` or `github.com/a-h/templ` (type-safe). Recommend `templ` for REPL complexity.
- Static: `embed.FS` for `web/static/*`.
- Docs: parse Markdown from `js/docs/` or `js/website/src/content/` at build time; render via `github.com/yuin/goldmark`.

### Styling

- Reuse `website/tailwind.config.cjs` → run `tailwindcss` CLI to build `web/static/css/app.css`. No Go CSS generation.

## 5.4 Umbrella `web` Package (`packages/web/`)

- `internal/web/web.go` re-exports `core`, `mini`, `tonal`, `transpiler`, `webaudio` equivalents — mirrors `packages/web/web.mjs` for `go run ./examples/...` ergonomics.

## 5.5 Examples (`examples/`)

Port each JS example to Go:

| JS example | Go equivalent |
|------------|---------------|
| `minimal-repl` | `examples/minimal-repl/main.go` — `core` + `mini` headless query loop |
| `buildless` | `examples/buildless/` — serve REPL without build step |
| `headless-repl` | `examples/headless/main.go` — evaluate + print Haps |
| `superdough` | `examples/superdough/main.go` — offline render to WAV |
| `tidal-repl` | `examples/tidal/main.go` — hs2js + tidal pattern |

## 5.6 Desktop (`src-tauri/` Rust → Go)

- **Recommended**: `github.com/wailsapp/wails/v2`.
- `wails.json` replaces `tauri.conf.json`; `cmd/desktop/main.go` replaces `src-tauri/src/main.rs`.
- Bridges: MIDI (`midir` → `gomidi`), OSC (`rosc` → `go-osc`), filesystem (`os` + `fs`), clipboard (`golang.design/x/clipboard`).
- Keep feature parity: productName `Strudel`, bundle `com.strudel.dev`, icons from `src-tauri/icons/`.

## 5.7 Reference / Docs Generation (`packages/reference/` + `undocumented.json`)

- `internal/reference/` generates `reference.json` from Go doc comments (via `go doc` or `gomarkdoc`), replacing `jsdoc --template jsdoc-json`.

## Acceptance Checklist

- [ ] `internal/draw` produces SVG/JSON matching JS fixtures; `TestDraw` green
- [ ] WASM builds: `make wasm` succeeds; `web/static/strudel.wasm` + `wasm_exec.js` present
- [ ] `go run ./cmd/strudel serve` boots at `http://localhost:8080`; REPL evaluates `s("bd sd")` and shows pianoroll
- [ ] `POST /api/evaluate` returns Haps JSON matching WASM path
- [ ] `POST /api/render` returns WAV
- [ ] `examples/headless` and `examples/superdough` run without error
- [ ] Desktop `wails build` produces binary (optional gate if Wails chosen)

