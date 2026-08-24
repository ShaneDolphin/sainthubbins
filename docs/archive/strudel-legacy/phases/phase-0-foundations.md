# Phase 0 — Foundations

> Duration: 1-2 weeks. No dependencies. This phase unblocks everything.

## Objective

Scaffold the Go module and port the **exact rational timing primitives** that all of Strudel depends on.

## 0.1 Repo Scaffold

```bash
cd Go
go mod init codeberg.org/uzu/strudel-go   # or github.com/<you>/strudel-go
cat > go.mod <<'GOMOD'
module codeberg.org/uzu/strudel-go
go 1.23
require (
    github.com/mna/pigeon v1.2.0
    github.com/dop251/goja v0.0.0-20240927123456-abcdef
)
GOMOD
mkdir -p cmd/strudel cmd/strudel-wasm internal/core internal/mini internal/transpiler internal/audio web/static web/templates
cat > Makefile <<'MF'
.PHONY: test lint wasm gen serve
test:  ; go test ./... -race -count=1
lint:  ; go vet ./... && golangci-lint run
wasm:  ; GOOS=js GOARCH=wasm go build -o web/static/strudel.wasm ./cmd/strudel-wasm && cp $$(go env GOROOT)/misc/wasm/wasm_exec.js web/static/
gen:   ; go generate ./...
serve: ; go run ./cmd/strudel serve
MF
cat > .golangci.yml <<'YML'
linters: { enable: [govet, errcheck, staticcheck, ineffassign, unused] }
YML
```

Add AGPL-3.0 `LICENSE` (copy from `js/LICENSE`) and `README.md`.

## 0.2 Fraction (`js/packages/core/fraction.mjs` 147 LOC + `fraction.js` dep)

**JS behavior to preserve**: exact rational with `n/d`, `add/sub/mul/div`, `mod`, `gcd/lcm`, `equals`, `toString`, `valueOf` (float), `parseNumeral`.

**Go file**: `internal/core/fraction.go`

```go
// Fraction — exact rational. n/d always normalized (d>0, gcd=1).
type Fraction struct { Num, Den int64 }
func NewFraction(n, d int64) Fraction
func FractionFromInt(n int64) Fraction
func FractionFromFloat(f float64) Fraction // via continued fraction or big.Rat
func (f Fraction) Add(g Fraction) Fraction
func (f Fraction) Sub(g Fraction) Fraction
func (f Fraction) Mul(g Fraction) Fraction
func (f Fraction) Div(g Fraction) Fraction
func (f Fraction) Mod(g Fraction) Fraction // JS _mod semantics (positive)
func (f Fraction) Cmp(g Fraction) int
func (f Fraction) Equals(g Fraction) bool
func (f Fraction) Float64() float64
func (f Fraction) String() string // "3/4" or "2" if den=1
// helpers: GCD, LCM, ParseNumeral(string) (Fraction, error)
```

- Normalize on construction: divide by `gcd(|n|, d)`, ensure `d>0`.
- Overflow: if `Num`/`Den` exceed `int64`, fallback to `*big.Int` internally or return error — document choice. Simplest: use `int64` and panic on overflow in tests; add `big.Rat` alternative if needed later.
- Tests: `fraction_test.go` — port arithmetic cases from `packages/core/test/util.test.mjs` + exhaustive `fraction.js` edge cases (negative, zero, large).

## 0.3 TimeSpan (`timespan.mjs` 117 LOC)

```go
type TimeSpan struct { Begin, End Fraction }
func NewTimeSpan(b, e Fraction) TimeSpan
func (s TimeSpan) Duration() Fraction
func (s TimeSpan) Equals(other TimeSpan) bool
func (s TimeSpan) Intersection(other TimeSpan) *TimeSpan // nil if disjoint
func (s TimeSpan) SpanCycles() []TimeSpan // splits at cycle boundaries
func (s TimeSpan) WithTime(fn func(Fraction) Fraction) TimeSpan
func (s TimeSpan) Midpoint() Fraction
```

- `Begin <= End` invariant; `Duration = End - Begin`.
- Test against JS `TimeSpan` fixtures.

## 0.4 Hap (`hap.mjs` 178 LOC)

```go
type Hap struct {
    Whole *TimeSpan         // nil for continuous signals
    Part  TimeSpan
    Value any               // usually map[string]any
    Context map[string]any  // source locations
    Stateful bool
}
func NewHap(whole *TimeSpan, part TimeSpan, value any) Hap
func (h Hap) Duration() Fraction
func (h Hap) EndClipped() Fraction
func (h Hap) WithValue(fn func(any) any) Hap
func (h Hap) WithValueMap(m map[string]any) Hap
func (h Hap) HasOnset() bool // Whole != nil && Whole.Begin == Part.Begin
func (h Hap) String() string
```

- Preserve `value.duration` / `value.clip` semantics: if `Value` is `map` with `duration`/`clip`, `Duration()` respects them.
- Immutability: all `With*` return copies.

## 0.5 State (`state.mjs` 28 LOC)

```go
type State struct { Span TimeSpan; Controls map[string]any }
func NewState(span TimeSpan, controls map[string]any) State
func (s State) SetSpan(span TimeSpan) State
func (s State) WithSpan(fn func(TimeSpan) TimeSpan) State
func (s State) SetControls(controls map[string]any) State
```

## 0.6 Util subset (`util.mjs` 508 LOC)

Port only what 0.2-0.5 need: `GCD`, `LCM`, `Mod`, `ID`, `Curry` (as func helpers), `ParseNumeral`, `StringifyValues`. Full `util.mjs` is completed in Phase 1.

## Acceptance Checklist

- [ ] `go test ./internal/core -run TestFraction|TestTimeSpan|TestHap|TestState` passes
- [ ] `fraction_test.go` covers negatives, zero, 0/1, large numerator/den, float round-trip
- [ ] `TimeSpan` intersection/spanCycles match JS fixtures (add `testdata/timespan.json` generated from JS)
- [ ] No `float64` used for timing comparisons — only `Fraction.Cmp`
- [ ] `LICENSE` present, `go vet` clean

