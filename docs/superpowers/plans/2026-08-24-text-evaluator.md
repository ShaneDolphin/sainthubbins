# Text Evaluator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `s("bd sd").fast(2).gain(0.8)` work as text, so the live console and `eval` become usable for composition instead of rhythm inspection only.

**Architecture:** goja is already a dependency and already used by `internal/transpiler`, but nothing binds the pattern API into a VM — `core.Evaluate` is passed a `nil` transpiler everywhere and falls back to mini-notation, then to a literal string. This plan adds `internal/jsapi`: a goja runtime with the pattern constructors as global functions and the transforms as chainable methods on a wrapped `core.Pattern`. Registration is table-driven so a new transform is one line, not a new binding.

**Tech Stack:** Go 1.25, `github.com/dop251/goja` (already in `go.mod`).

**Spec:** `docs/superpowers/specs/2026-08-24-remaining-work.md`

## Global Constraints

- Go 1.25.0, module `codeberg.org/uzu/saint-hubbins`. goja is the only third-party package; add no others.
- AGPL-3.0-or-later header on every new file.
- **Mini-notation must keep working unchanged.** `eval "bd sd"` and every existing template depend on it. The JS path is tried first and falls back to the mini parser, never the reverse.
- A JS error must surface as an error, not as a literal-string hap. That silent fallback is the bug this plan exists to remove.
- `go test ./... -race -count=1` and `go vet ./...` must stay clean.
- goja is not goroutine-safe: one `*goja.Runtime` per evaluation, or guarded by a mutex. The web console handles concurrent requests.

## File Structure

| File | Responsibility |
|---|---|
| `internal/jsapi/runtime.go` (create) | VM construction, global registration, `Evaluate`. |
| `internal/jsapi/pattern.go` (create) | The JS-facing pattern object and its chainable methods. |
| `internal/jsapi/registry.go` (create) | Tables mapping JS names to Go operations. |
| `internal/jsapi/*_test.go` (create) | Behaviour tests per file. |
| `cmd/saint-hubbins/main.go` (modify) | `eval`/`render` try JS first. |
| `web/server.go` (modify) | `/api/evaluate` tries JS first and reports errors. |

---

### Task 1: A VM that returns a pattern

**Files:**
- Create: `internal/jsapi/runtime.go`, `internal/jsapi/pattern.go`, `internal/jsapi/runtime_test.go`

**Interfaces:**
- Produces: `func Evaluate(code string) (core.Pattern, error)`, and unexported `type jsPattern struct { pat core.Pattern }`.

Start with one constructor — `s` — and prove the round trip: JS text in, `core.Pattern` out.

- [ ] **Step 1: Write the failing test**

```go
// internal/jsapi/runtime_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func hapValues(t *testing.T, p core.Pattern) []any {
	t.Helper()
	var out []any
	for _, h := range p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)) {
		out = append(out, h.Value)
	}
	return out
}

func TestEvaluateSoundConstructor(t *testing.T) {
	p, err := Evaluate(`s("bd sd")`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	vals := hapValues(t, p)
	if len(vals) != 2 {
		t.Fatalf("got %d haps, want 2", len(vals))
	}
	m, ok := vals[0].(map[string]any)
	if !ok {
		t.Fatalf("value is %T, want a control bag", vals[0])
	}
	if m["s"] != "bd" {
		t.Errorf("s = %v, want bd", m["s"])
	}
}

// A syntax error must be reported, not silently turned into a literal hap.
func TestEvaluateReportsErrors(t *testing.T) {
	if _, err := Evaluate(`s("bd" +`); err == nil {
		t.Fatal("want a syntax error, got nil")
	}
	if _, err := Evaluate(`notAFunction("x")`); err == nil {
		t.Fatal("want a reference error, got nil")
	}
}

func TestEvaluateRejectsNonPatternResult(t *testing.T) {
	if _, err := Evaluate(`42`); err == nil {
		t.Fatal("want an error when the result is not a pattern")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/jsapi/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Write minimal implementation**

```go
// internal/jsapi/pattern.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// The JS-facing pattern object.

package jsapi

import "codeberg.org/uzu/saint-hubbins/internal/core"

// jsPattern wraps a core.Pattern so JS can hold and chain it. Methods are
// attached in registry.go rather than declared here, so adding a transform is
// a table entry rather than a new method.
type jsPattern struct {
	pat core.Pattern
}
```

```go
// internal/jsapi/runtime.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// A goja runtime with the pattern API bound into it.

package jsapi

import (
	"fmt"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// Evaluate runs code in a fresh VM and returns the pattern it produced.
//
// A fresh runtime per call keeps this safe under the web console's concurrent
// requests — goja runtimes are not goroutine-safe — and costs little next to
// pattern querying.
func Evaluate(code string) (core.Pattern, error) {
	mini.RegisterStringParser()
	vm := goja.New()
	if err := register(vm); err != nil {
		return core.Silence(), err
	}
	v, err := vm.RunString(code)
	if err != nil {
		return core.Silence(), fmt.Errorf("jsapi: %w", err)
	}
	return unwrap(vm, v)
}

// unwrap converts a JS result into a Pattern. A bare string is treated as
// mini-notation so `"bd sd"` works, but anything else is an error rather than
// a literal-valued hap.
func unwrap(vm *goja.Runtime, v goja.Value) (core.Pattern, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return core.Silence(), fmt.Errorf("jsapi: expression produced no value")
	}
	if obj, ok := v.Export().(*jsPattern); ok {
		return obj.pat, nil
	}
	if s, ok := v.Export().(string); ok {
		return mini.Mini(s), nil
	}
	return core.Silence(), fmt.Errorf("jsapi: expression produced %T, want a pattern", v.Export())
}
```

```go
// internal/jsapi/registry.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Tables binding JS names to engine operations.

package jsapi

import (
	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// controls are the constructors that turn a value into a control pattern.
// The JS name is the key; the Go constructor is the value.
var controls = map[string]func(any) core.Pattern{
	"s": core.S, "sound": core.S, "note": core.Note, "n": core.N,
	"gain": core.Gain, "cutoff": core.Cutoff, "lpf": core.Lpf,
	"pan": core.Pan, "room": core.Room, "speed": core.Speed,
	"attack": core.Attack, "release": core.Release, "shape": core.Shape,
}

// toPattern coerces a JS argument into a Pattern: a wrapped pattern passes
// through, a string is mini-notation, a number is a constant.
func toPattern(v any) core.Pattern {
	switch x := v.(type) {
	case *jsPattern:
		return x.pat
	case string:
		return mini.Mini(x)
	case float64, int, int64:
		return core.Pure(x)
	}
	return core.Silence()
}

// register installs every global into the VM.
func register(vm *goja.Runtime) error {
	wrap := func(p core.Pattern) goja.Value { return vm.ToValue(newJSPattern(vm, p)) }

	for name, ctor := range controls {
		ctor := ctor
		if err := vm.Set(name, func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return wrap(core.Silence())
			}
			arg := call.Argument(0).Export()
			// A string argument is mini-notation, so s("bd sd") is a sequence.
			if str, ok := arg.(string); ok {
				return wrap(ctor(mini.Mini(str)))
			}
			if jp, ok := arg.(*jsPattern); ok {
				return wrap(ctor(jp.pat))
			}
			return wrap(ctor(arg))
		}); err != nil {
			return err
		}
	}
	return nil
}
```

Now the object with methods. Add to `pattern.go`:

```go
// newJSPattern builds the JS object for a pattern: the wrapper plus every
// chainable method from the tables in registry.go.
func newJSPattern(vm *goja.Runtime, p core.Pattern) *goja.Object {
	jp := &jsPattern{pat: p}
	obj := vm.ToValue(jp).(*goja.Object)
	attachMethods(vm, obj, jp)
	return obj
}
```

and to `registry.go`:

```go
// attachMethods is filled in by Task 2. For now it exists so Task 1 compiles.
func attachMethods(vm *goja.Runtime, obj *goja.Object, jp *jsPattern) {}
```

Add the `goja` import to `pattern.go`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/jsapi/ -v`
Expected: PASS (3 tests). If goja exports the wrapper as something other than `*jsPattern`, inspect with `fmt.Printf("%T", v.Export())` and adjust `unwrap` — goja's Go-struct binding is the fiddly part of this task.

- [ ] **Step 5: Commit**

```bash
git add internal/jsapi/
git commit -m "feat(jsapi): evaluate pattern constructors from JS text"
```

---

### Task 2: Chainable transforms

**Files:**
- Modify: `internal/jsapi/registry.go`
- Test: `internal/jsapi/transform_test.go` (create)

**Interfaces:**
- Produces: `attachMethods` gains real behaviour, driven by three tables: `unaryOps`, `numericOps` and the `controls` map reused as setters.

- [ ] **Step 1: Write the failing test**

```go
// internal/jsapi/transform_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func countHaps(t *testing.T, code string) int {
	t.Helper()
	p, err := Evaluate(code)
	if err != nil {
		t.Fatalf("Evaluate(%q): %v", code, err)
	}
	return len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
}

func TestFastDoublesEvents(t *testing.T) {
	if got := countHaps(t, `s("bd sd").fast(2)`); got != 4 {
		t.Errorf("got %d haps, want 4", got)
	}
}

func TestChainedCalls(t *testing.T) {
	p, err := Evaluate(`s("bd").gain(0.5).cutoff(800)`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	m := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0].Value.(map[string]any)
	if m["s"] != "bd" || m["gain"] != 0.5 || m["cutoff"] != 800.0 {
		t.Errorf("got %v, want s:bd gain:0.5 cutoff:800", m)
	}
}

func TestUnaryTransform(t *testing.T) {
	p, err := Evaluate(`s("bd sd").rev()`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	first := haps[0].Value.(map[string]any)
	if first["s"] != "bd" || haps[0].Part.Begin.Float64() != 0.5 {
		t.Errorf("rev did not move bd to the second half: %v @ %v",
			first, haps[0].Part.Begin.Float64())
	}
}

func TestEuclidFromJS(t *testing.T) {
	if got := countHaps(t, `s("bd").euclid(3, 8)`); got != 3 {
		t.Errorf("got %d haps, want 3", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/jsapi/ -run "Fast|Chained|Unary|Euclid" -v`
Expected: FAIL — `Object has no member 'fast'`

- [ ] **Step 3: Write minimal implementation**

Replace the stub `attachMethods` in `registry.go`:

```go
// unaryOps take no arguments.
var unaryOps = map[string]func(core.Pattern) core.Pattern{
	"rev":        func(p core.Pattern) core.Pattern { return p.Rev() },
	"palindrome": func(p core.Pattern) core.Pattern { return p.Palindrome() },
	"degrade":    func(p core.Pattern) core.Pattern { return p.Degrade() },
	"hush":       func(core.Pattern) core.Pattern { return core.Silence() },
}

// numericOps take a single number.
var numericOps = map[string]func(core.Pattern, float64) core.Pattern{
	"fast":      func(p core.Pattern, v float64) core.Pattern { return p.FastF(core.FractionFromFloat(v)) },
	"slow":      func(p core.Pattern, v float64) core.Pattern { return p.SlowF(core.FractionFromFloat(v)) },
	"ply":       func(p core.Pattern, v float64) core.Pattern { return p.Ply(int(v)) },
	"segment":   func(p core.Pattern, v float64) core.Pattern { return p.Segment(v) },
	"late":      func(p core.Pattern, v float64) core.Pattern { return p.Late(v) },
	"early":     func(p core.Pattern, v float64) core.Pattern { return p.Early(v) },
	"degradeBy": func(p core.Pattern, v float64) core.Pattern { return p.DegradeBy(v) },
	"add":       func(p core.Pattern, v float64) core.Pattern { return p.Add(v) },
}

func attachMethods(vm *goja.Runtime, obj *goja.Object, jp *jsPattern) {
	wrap := func(p core.Pattern) goja.Value { return vm.ToValue(newJSPattern(vm, p)) }

	for name, op := range unaryOps {
		op := op
		_ = obj.Set(name, func(goja.FunctionCall) goja.Value { return wrap(op(jp.pat)) })
	}

	for name, op := range numericOps {
		op := op
		_ = obj.Set(name, func(call goja.FunctionCall) goja.Value {
			return wrap(op(jp.pat, call.Argument(0).ToFloat()))
		})
	}

	// Controls double as setters when called on a pattern: .gain(0.5) merges
	// a gain control into every event.
	for name, ctor := range controls {
		ctor := ctor
		_ = obj.Set(name, func(call goja.FunctionCall) goja.Value {
			arg := call.Argument(0).Export()
			if str, ok := arg.(string); ok {
				return wrap(jp.pat.Set(ctor(mini.Mini(str))))
			}
			return wrap(jp.pat.Set(ctor(arg)))
		})
	}

	// euclid takes two arguments, so it is not in the numeric table.
	_ = obj.Set("euclid", func(call goja.FunctionCall) goja.Value {
		return wrap(jp.pat.Euclid(int(call.Argument(0).ToInteger()), int(call.Argument(1).ToInteger())))
	})

	// every and off take a callback, which must be re-entered as a pattern op.
	_ = obj.Set("every", func(call goja.FunctionCall) goja.Value {
		n := int(call.Argument(0).ToInteger())
		fn, ok := goja.AssertFunction(call.Argument(1))
		if !ok {
			return wrap(jp.pat)
		}
		return wrap(jp.pat.Every(n, func(p core.Pattern) core.Pattern {
			res, err := fn(goja.Undefined(), vm.ToValue(newJSPattern(vm, p)))
			if err != nil {
				return p
			}
			if out, ok := res.Export().(*jsPattern); ok {
				return out.pat
			}
			return p
		}))
	})
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/jsapi/ -v`
Expected: PASS (all tests)

- [ ] **Step 5: Commit**

```bash
git add internal/jsapi/
git commit -m "feat(jsapi): add chainable transforms and control setters"
```

---

### Task 3: Top-level combinators

**Files:**
- Modify: `internal/jsapi/registry.go`
- Test: `internal/jsapi/combinator_test.go` (create)

**Interfaces:**
- Produces: globals `stack`, `cat`, `slowcat`, `fastcat`, `sequence`, `silence`, `mini`.

- [ ] **Step 1: Write the failing test**

```go
// internal/jsapi/combinator_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import "testing"

func TestStackLayers(t *testing.T) {
	if got := countHaps(t, `stack(s("bd*4"), s("hh*8"))`); got != 12 {
		t.Errorf("got %d haps, want 12", got)
	}
}

func TestCatAlternates(t *testing.T) {
	if got := countHaps(t, `cat(s("bd"), s("sd"))`); got != 1 {
		t.Errorf("got %d haps in one cycle, want 1", got)
	}
}

func TestSilenceIsEmpty(t *testing.T) {
	if got := countHaps(t, `silence()`); got != 0 {
		t.Errorf("got %d haps, want 0", got)
	}
}

func TestMiniHelper(t *testing.T) {
	if got := countHaps(t, `mini("bd sd hh")`); got != 3 {
		t.Errorf("got %d haps, want 3", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/jsapi/ -run "Stack|Cat|Silence|MiniHelper" -v`
Expected: FAIL — `ReferenceError: stack is not defined`

- [ ] **Step 3: Write minimal implementation**

Append inside `register`, before the closing `return nil`:

```go
	// Variadic combinators.
	variadic := map[string]func(...core.Pattern) core.Pattern{
		"stack":    core.Stack,
		"cat":      core.Cat,
		"slowcat":  core.SlowCat,
		"fastcat":  core.FastCat,
		"sequence": core.Sequence,
	}
	for name, fn := range variadic {
		fn := fn
		if err := vm.Set(name, func(call goja.FunctionCall) goja.Value {
			pats := make([]core.Pattern, 0, len(call.Arguments))
			for _, a := range call.Arguments {
				pats = append(pats, toPattern(a.Export()))
			}
			return wrap(fn(pats...))
		}); err != nil {
			return err
		}
	}

	if err := vm.Set("silence", func(goja.FunctionCall) goja.Value {
		return wrap(core.Silence())
	}); err != nil {
		return err
	}
	// mini() is an explicit escape hatch to the rhythm language.
	if err := vm.Set("mini", func(call goja.FunctionCall) goja.Value {
		return wrap(mini.Mini(call.Argument(0).String()))
	}); err != nil {
		return err
	}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/jsapi/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/jsapi/
git commit -m "feat(jsapi): add stack, cat and the other top-level combinators"
```

---

### Task 4: Wire it into the CLI and the console

**Files:**
- Modify: `cmd/saint-hubbins/main.go`, `web/server.go`
- Test: `web/server_jsapi_test.go` (create)

**Interfaces:**
- Consumes: `jsapi.Evaluate`.
- Produces: `func evaluateCode(code string) (core.Pattern, error)` in both call sites — JS first, mini-notation fallback, error surfaced.

The fallback order matters. Mini-notation like `bd sd` is not valid JS and must still work, so: try JS; if it fails **and** the mini parser yields haps, use mini; otherwise report the JS error.

- [ ] **Step 1: Write the failing test**

```go
// web/server_jsapi_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postEvaluate(t *testing.T, code string) map[string]any {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	NewServer("").Handler().ServeHTTP(rec, req)

	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("bad JSON: %v (%s)", err, rec.Body.String())
	}
	return out
}

func TestConsoleEvaluatesJS(t *testing.T) {
	out := postEvaluate(t, `s("bd sd")`)
	haps, _ := out["haps"].([]any)
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2 — the console should evaluate JS now", len(haps))
	}
	first := haps[0].(map[string]any)
	val, ok := first["value"].(map[string]any)
	if !ok || val["s"] != "bd" {
		t.Errorf("value = %v, want a control bag carrying s:bd", first["value"])
	}
}

func TestConsoleStillEvaluatesMiniNotation(t *testing.T) {
	out := postEvaluate(t, "bd*4")
	haps, _ := out["haps"].([]any)
	if len(haps) != 4 {
		t.Errorf("got %d haps, want 4 — mini-notation must keep working", len(haps))
	}
}

func TestConsoleReportsRealErrors(t *testing.T) {
	out := postEvaluate(t, `s("bd").nope()`)
	if out["error"] == nil || out["error"] == "" {
		t.Errorf("a bad method should be reported, got %v", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./web/ -run Console -v`
Expected: FAIL — the value is the literal string `s("bd sd")`

- [ ] **Step 3: Write minimal implementation**

Add to `web/server.go`:

```go
// evaluateCode resolves user input to a pattern. JS is tried first; input that
// is not valid JS but parses as mini-notation falls back to the rhythm
// language. Anything else returns the JS error rather than a literal hap.
func evaluateCode(code string) (core.Pattern, error) {
	mini.RegisterStringParser()
	if pat, err := jsapi.Evaluate(code); err == nil {
		return pat, nil
	} else {
		if m := mini.Mini(code); m.Query != nil {
			if len(m.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) > 0 {
				return m, nil
			}
		}
		return core.Silence(), err
	}
}
```

Replace the evaluation block in both `handleEvaluate` and `handlePianoroll` with a call to it, and set `resp.Error = err.Error()` when it fails. Import `codeberg.org/uzu/saint-hubbins/internal/jsapi`.

Make the same change in `cmd/saint-hubbins/main.go` for `demoEval` and `demoRender`, printing the error to stderr and exiting non-zero rather than rendering silence.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./web/ ./cmd/... -v`
Expected: PASS

Note `web/server_test.go` and `web/server_extra3_test.go` post `s("bd sd")` and assert only that the endpoint does not crash. They should still pass, and now exercise the real path.

- [ ] **Step 5: Update the console UI**

`web/server.go`'s console template says mini-notation only and lists rhythm examples. Update the hint to show both: `s("bd*4").gain(0.8)` alongside `bd*4`. Change the default textarea to `stack(s("bd*4"), s("hh*8").gain(0.4))`.

- [ ] **Step 6: Full verification and commit**

```bash
go vet ./... && go test ./... -race -count=1
go run ./cmd/saint-hubbins eval 's("bd sd")'
```

Expected: two haps carrying `s` controls — the roadmap gate, met for real.

```bash
git add web/ cmd/ internal/jsapi/
git commit -m "feat: evaluate JS pattern code in the CLI and console"
```

---

### Task 5: Retire the documented limitation

**Files:**
- Modify: `docs/tutorial/02-mini-notation.md`, `docs/tutorial/03-patterns-in-go.md`, `docs/tutorial/07-new-song-web.md`, `docs/tutorial/08-limitations.md`, `README.md`, `docs/05-execution-checklist.md`

Several chapters state plainly that `s("bd sd")` does not work and that anything beyond rhythm is Go-only. That was true and is now the single largest doc change in the project.

- [ ] **Step 1: Find every claim**

```bash
grep -rn 's("bd sd")\|not implemented as text\|mini-notation only\|no script evaluator' docs/ README.md
```

- [ ] **Step 2: Rewrite each**

- `02-mini-notation.md`: the closing "What mini-notation cannot do" section should now say the console accepts both, and point at the JS API for layering and controls.
- `03-patterns-in-go.md`: "Why the split exists" needs rewriting — Go is still how songs are written to files, but text is no longer rhythm-only.
- `07-new-song-web.md`: the console is no longer a rhythm sketchpad; rewrite the session walkthrough around JS, keeping the mini-notation examples.
- `08-limitations.md`: delete "The console evaluates mini-notation only".
- `README.md`: update the CLI Reference note added earlier.
- `docs/05-execution-checklist.md`: the `s("bd sd")` gate can now be a real assertion — see the repo-hygiene plan's `scripts/check.sh` and strengthen that check to assert an `s` control rather than any hap.

- [ ] **Step 3: Verify and commit**

```bash
./scripts/check.sh   # if the repo-hygiene plan has landed
go test ./... -race -count=1
git add docs/ README.md
git commit -m "docs: the console now evaluates JS, not just mini-notation"
```

---

## Scope note

This is the largest of the five plans and the least load-bearing: everything else works without it, and the templates never use JS. If effort is limited, do the real-time OSC plan and the engine-gaps plan first — they change what the software can do, where this changes how comfortable it is to drive.
