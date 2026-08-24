# 04 — Conventions, Testing, CI & Risks

## Go Conventions

### Naming

- Packages: `core`, `mini`, `transpiler`, `audio`, `draw`, `tonal`, `edo`, `xen` — same as JS package names, lowercase, no `strudel` prefix.
- Exported funcs: `CamelCase` matching JS `camelCase` intent (`S` for `s`, `Note` for `note`); add alias `var S = Sound` where JS has aliases.
- Files: `fraction.go`, `pattern_core.go` — `snake_case.go` per Go convention.

### Error Handling

- JS throws; Go returns `error`. Every `Pattern` combinator that can fail (e.g., `Euclid` with 0 steps) returns `(Pattern, error)` — provide `Must*` panic wrapper for tests/examples.
- `Evaluate` returns `(*Pattern, error)` — never panic on user input.

### Logging

- Replace `logger.mjs` / `loggerbridge.mjs` with `log/slog` (Go 1.21+). Levels: Debug for pattern queries, Info for REPL, Warn for audio clipping, Error for transpiler failures.

### Codegen

- `go:generate` directives in `internal/core/controls.go`, `internal/soundfonts/gm.go`.
- Generators under `tools/gen-controls/`, `tools/gen-gm/` — run `make gen`.

### Headers

Every `.go` file starts with AGPL-3.0 header (copy from JS, change year/language):

```go
// Copyright (C) 2026 Strudel contributors — AGPL-3.0-or-later
// Ported from https://codeberg.org/uzu/strudel (JS) to Go.
```

## Testing Strategy

### Unit Tests (`*_test.go`)

- Mirror JS `*.test.mjs` structure. Example: `packages/core/test/pattern.test.mjs` → `internal/core/pattern_test.go`.
- Use `testing` + `github.com/stretchr/testify` (assert/require) — or stdlib only if preferred.
- Fixture files: `testdata/` JSON generated from JS via `node --experimental-vm-modules` script that runs JS and dumps Hap arrays.

### Cross-Validation via goja

For phases 1-2, embed JS fixtures via `goja`:

```go
func TestPatternParity(t *testing.T) {
    vm := goja.New()
    // load js/packages/core/*.mjs via goja
    expected := runJS(t, vm, `pattern("bd sd").query(state(...))`)
    got := core.Cat(core.Pure("bd"), core.Pure("sd")).Query(state)
    assert.Equal(t, expected, toJSON(got))
}
```

Remove `goja` dependency after pure Go tests are green, or keep behind `//go:build with_goja` tag.

### Snapshots

- JS `test/__snapshots__/` → Go `testdata/snapshots/*.json`.
- Use `github.com/sebdah/goldie` or manual `os.ReadFile` + `require.JSONEq`.

### Race & Bench

- `go test -race ./...` in CI.
- `go test -bench=. -benchmem ./internal/core` replaces `vitest bench` (`bench/tunes.bench.mjs`).

## CI Pipeline (replaces `pnpm check`)

```yaml
# .github/workflows/go.yml (or .forgejo/workflows/)
on: [push, pull_request]
jobs:
  ci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: go vet ./...
      - run: golangci-lint run
      - run: go test ./... -race -count=1 -coverprofile=coverage.out
      - run: GOOS=js GOARCH=wasm go build -o web/static/strudel.wasm ./cmd/strudel-wasm
```

### Replacing JS tooling

| JS | Go |
|----|----|
| `eslint` | `golangci-lint` |
| `prettier` | `gofmt` + `goimports` |
| `vitest` | `go test` |
| `jsdoc` | `pkgsite` / `gomarkdoc` |
| `vite` | `go build` + `esbuild` (for static JS) |
| `pnpm i` | `go mod download` |

## Risks & Mitigations

| # | Risk | Impact | Mitigation |
|---|------|--------|------------|
| 1 | Pattern semantics drift (JS dynamic → Go static) | High — silent wrong music | Golden-file cross-validation via goja for every combinator; port `pattern.test.mjs` 1:1. |
| 2 | `fraction.js` exactness vs Go float | High — timing drift | Custom `Fraction` with `int64` + `GCD`; test with JS fixtures; never use `float64` for spans. |
| 3 | Audio fidelity (superdough DSP) | High — audible difference | Offline WAV render + RMS comparison; vendor JS DSP constants. |
| 4 | Mini PEG divergence (`krill.pegjs` → `pigeon`) | Medium — parse errors | Keep `.peg` as truth; round-trip test all `testtunes.mjs` mini strings. |
| 5 | Transpiler complexity (`acorn` + 4 plugins) | Medium — eval breaks | Two-phase: goja shim first, pure Go second; `tunes.test.mjs` as gate. |
| 6 | WASM bundle size | Low — slow load | Build with `-ldflags="-s -w"`; tree-shake via `tinygo` evaluation later if needed. |
| 7 | Tailwind/Astro parity | Low — visual diff | Not a gate; minimal Go templates + reused CSS is sufficient for v0.1. |
| 8 | Tauri → Wails desktop gap | Low — desktop not core | Mark desktop as post-MVP; native CLI + browser REPL are the MVP. |

## Out-of-Scope / Deferred

- `packages/vite-plugin-bundle-audioworklet` — not needed; Go worklets are plain Go files.
- `hs2js` tree-sitter Haskell parser — if Go `tree-sitter` binding is heavy, keep `hs2js` as WASM JS and bridge via `syscall/js`.
- `mondo`/`mondolang` — depends on external `mondolang` package; stub as `Pattern` passthrough until needed.

## Glossary

| Term | Meaning |
|------|---------|
| Hap | Event with whole/part TimeSpans + value bag |
| CPS | Cycles per second (tempo) |
| Mini | Tidal mini-notation (`"bd ~ sd"`) |
| Zyklus | Clock/scheduler loop |
| Superdough | Main DSP engine |
| SoundMap | Registry of sample/synth triggers (`nanostores` → `sync.Map`) |
| WASM bridge | `syscall/js` calls from Go WASM to browser APIs |

