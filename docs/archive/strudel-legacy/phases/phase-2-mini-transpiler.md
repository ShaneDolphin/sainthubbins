# Phase 2 — Mini Notation + Transpiler

> Duration: 2-3 weeks. Depends on Phase 1.

## Objective

Port the **mini-notation parser** and the **JS transpiler** so user strings like `"bd ~ sd ~"` and code like `s("bd sd")` evaluate identically to JS.

## 2.1 Mini Parser (`packages/mini/`)

### Source of truth

- `krill.pegjs` — PEG grammar for mini-notation
- `krill-parser.js` — 2497 LOC **generated** file (do not hand-port)
- `mini.mjs` — `patternifyAST(ast, code, onEnter)` (386 LOC core)

### Go approach

1. Copy `krill.pegjs` to `internal/mini/krill.peg`.
2. Generate parser with `pigeon`:

```bash
go install github.com/mna/pigeon@latest
pigeon -o internal/mini/parser.go internal/mini/krill.peg
# check in both .peg and parser.go
```

3. Implement `mini.go`:

```go
package mini

// AST types mirror krill-parser.js type_ fields
type AST struct { Type string; Source []AST; Options *Options; ... }

func ParseMini(input string) (*AST, error) // wraps generated parser
func PatternifyAST(ast *AST, code string, onEnter func(*AST)) core.Pattern
```

- `PatternifyAST` must replicate `mini.mjs:applyOptions` switch: `stretch` (fast/slow), `replicate` (repeatCycles+fast), `bjorklund` (euclid/euclidRot), `degradeBy`, `tail`, `range`.
- `randOffset = 0.0003` constant preserved.
- `SetStringParser` integration: `core.SetStringParser(func(s string) core.Pattern { ast,_ := ParseMini(s); return PatternifyAST(ast, s, nil) })`.

### Alternatives if pigeon fails

- Use `participle/v2` with custom lexer.
- As last resort, run `krill-parser.js` via `goja` and bridge AST as JSON — acceptable for Phase 2.1, but Phase 2 final must be pure Go.

### Tests

- Port `packages/mini/test/*.mjs` + `bench/mini.bench.mjs` fixtures.
- Golden file: `testdata/mini/*.mini` → expected `[]Hap` JSON (generated from JS).

## 2.2 Transpiler (`packages/transpiler/`)

### JS behavior

`transpiler.mjs` (348 LOC): `transpiler(input, {wrapAsync, addReturn, emitMiniLocations, emitWidgets, blockBased, range})` — parses JS with `acorn`, walks AST with `estree-walker`, applies plugins, emits via `escodegen`.

Plugins:
- `plugin-mini.mjs` — rewrites mini template literals
- `plugin-sample.mjs` — sample string handling
- `plugin-kabelsalat.mjs` — Kabelsalat integration
- `plugin-widgets.mjs` — widget/slider injection

Also `registerLanguage` / `registerTranspilerPlugin` registry.

### Go: Two-phase implementation

**Phase 2a — goja shim (1 week, unblocks tests)**

```go
// internal/transpiler/shim.go //go:build !purego
package transpiler
import "github.com/dop251/goja"
func TranspileJS(input string, opts Options) (Result, error) {
    vm := goja.New()
    // load transpiler.mjs via embed.FS, run with goja
}
```

Lets `tunes.test.mjs` pass immediately.

**Phase 2b — Pure Go transpiler (1-2 weeks)**

```go
// internal/transpiler/transpiler.go
type Options struct { WrapAsync bool; AddReturn bool; EmitMiniLocations bool; EmitWidgets bool; BlockBased bool; Range [2]int }
type Result struct { Output string; MiniLocations []Location; Widgets []Widget }
func Transpile(input string, opts Options) (Result, error)
```

- Parser: `github.com/dop251/goja/parser` or `github.com/robertkrimen/otto/parser` (JS parser in Go).
- Walker: implement `estree`-like walk over `goja/ast` nodes.
- Plugins: Go interfaces `type Plugin interface { Walk(ctx *Context, node ast.Node) }`.
- Emitter: custom JS emitter or `goja`'s.

### Evaluation integration

```go
// internal/core/evaluate.go — now uses Go transpiler
func Evaluate(code string, opts transpiler.Options) (Pattern, error) {
    res, _ := transpiler.Transpile(code, opts)
    return evalWithGoja(res.Output) // or pure Go eval if feasible
}
```

## Acceptance Checklist

- [ ] `go test ./internal/mini -count=1` passes (all ported mini tests)
- [ ] Mini strings produce identical Hap sets to JS (goja comparison) for 30+ cases including `euclid`, `degradeBy`, `tail`, `range`
- [ ] Transpiler shim passes `tunes.test.mjs` fixtures (or pure Go transpiler does)
- [ ] `registerLanguage` / `registerTranspilerPlugin` equivalents exist and are tested
- [ ] `TestTranspileRoundTrip` — JS-in, JS-out, eval, query — matches JS end-to-end

