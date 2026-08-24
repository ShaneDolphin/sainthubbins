# Pattern Engine Gaps Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix the five confirmed engine defects so the mini-notation grammar behaves the way Strudel users expect and the tutorial's limitations chapter shrinks.

**Architecture:** Four fixes live in `internal/mini/mini.go` (weight, replicate, polymeter) and `internal/core` (`SlowCat`, `Add`). Each is independent; do them in any order, but Task 1 first — it is a one-line fix with the largest blast radius. Tasks 2–4 change documented behaviour, so each ends by updating `docs/tutorial/02-mini-notation.md` and `docs/tutorial/08-limitations.md` in the same commit.

**Tech Stack:** Go 1.25 standard library. `TimeCatWeighted` (`internal/core/pattern_weighted.go:9`) and `SplitQueries` (`internal/core/pattern.go:354`) already exist and are the right primitives.

**Spec:** `docs/superpowers/specs/2026-08-24-remaining-work.md`

## Global Constraints

- Go 1.25.0, module `codeberg.org/uzu/saint-hubbins`. No new dependencies.
- AGPL-3.0-or-later header on every new file.
- `go test ./... -race -count=1` and `go vet ./...` must stay clean.
- **Documentation is part of the change.** `docs/tutorial/02-mini-notation.md` documents current behaviour including the deviations, and `docs/tutorial/08-limitations.md` lists them. Fixing behaviour without updating both leaves the tutorial lying.
- `internal/mini/mini_nested_ops_test.go::TestMiniDocumentedSyntax` pins hap counts for the documented grammar. When behaviour changes deliberately, update that table in the same commit — never delete the case.

## File Structure

| File | Responsibility |
|---|---|
| `internal/core/pattern.go` (modify) | `SlowCat` cycle mapping. |
| `internal/core/pattern_cycle_test.go` (create) | Regression tests for `SlowCat` across cycle boundaries. |
| `internal/mini/mini.go` (modify) | Weight (`@`), replicate (`!`), polymeter steps (`%`). |
| `internal/mini/mini_grammar_test.go` (create) | Tests for the three grammar fixes. |
| `internal/core/pattern_arith.go` (create) | Control-bag-aware arithmetic for `Add`. |
| `internal/core/pattern_arith_test.go` (create) | Tests for `Add` on wrapped patterns. |

---

### Task 1: SlowCat must not collapse the cycle span

**Files:**
- Modify: `internal/core/pattern.go` (the `SlowCat` function)
- Test: `internal/core/pattern_cycle_test.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: no signature change. `SlowCat(pats ...Pattern) Pattern` behaves correctly for arguments containing several events.

Root cause, confirmed by running it: `SlowCat` maps each per-cycle span into the argument's local time with

```go
mappedSpan := cyc.WithTime(func(t Fraction) Fraction { return t.Sub(t.Sam()) })
```

`Sam()` is the start of the cycle containing `t`. For the span `[1,2)` the begin maps to `1-1 = 0` but the end maps to `2-2 = 0`, producing the zero-width span `0/1 → 0/1`. A whole-cycle argument still answers that point query, which is why `SlowCat(Silence(), Pure("bd"))` appears to work; anything with events at sub-cycle positions returns nothing. The correct base is the cycle's own start, which the function already computes nine lines further down as `base := cyc.Begin.Sam()`.

- [ ] **Step 1: Write the failing test**

```go
// internal/core/pattern_cycle_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// SlowCat must give each argument a full cycle, not a zero-width span.

package core

import "testing"

// fourOnTheFloor has events at sub-cycle positions, which is what exposes the
// collapsed-span bug — a single whole-cycle event does not.
func fourOnTheFloor() Pattern {
	return FastCat(Pure("bd"), Pure("bd"), Pure("bd"), Pure("bd"))
}

func TestSlowCatKeepsMultiEventArguments(t *testing.T) {
	p := SlowCat(Silence(), fourOnTheFloor())
	for cycle := int64(0); cycle < 6; cycle++ {
		got := len(p.QueryArc(FractionFromInt(cycle), FractionFromInt(cycle+1)))
		want := 0
		if cycle%2 == 1 {
			want = 4
		}
		if got != want {
			t.Errorf("cycle %d: got %d haps, want %d", cycle, got, want)
		}
	}
}

func TestSlowCatWideQueryMatchesPerCycle(t *testing.T) {
	p := SlowCat(fourOnTheFloor(), Silence(), fourOnTheFloor(), Silence())
	whole := len(p.QueryArc(FractionFromInt(0), FractionFromInt(8)))
	split := 0
	for c := int64(0); c < 8; c++ {
		split += len(p.QueryArc(FractionFromInt(c), FractionFromInt(c+1)))
	}
	if whole != split {
		t.Errorf("one 8-cycle query gave %d haps, eight 1-cycle queries gave %d", whole, split)
	}
	if whole != 8 {
		t.Errorf("got %d haps over 8 cycles, want 8 (four on two of every four bars)", whole)
	}
}

// Events must land at the right offset inside their cycle, not just be present.
func TestSlowCatPreservesEventPositions(t *testing.T) {
	p := SlowCat(Silence(), fourOnTheFloor())
	haps := p.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(haps) != 4 {
		t.Fatalf("got %d haps, want 4", len(haps))
	}
	want := []float64{1.0, 1.25, 1.5, 1.75}
	for i, w := range want {
		if got := haps[i].Part.Begin.Float64(); got != w {
			t.Errorf("hap %d begins at %v, want %v", i, got, w)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/core/ -run SlowCat -v`
Expected: FAIL — `TestSlowCatKeepsMultiEventArguments` reports 0 haps where 4 are wanted.

- [ ] **Step 3: Write minimal implementation**

In `internal/core/pattern.go`, inside `SlowCat`, hoist the base and use it for the mapping. Replace:

```go
			mappedSpan := cyc.WithTime(func(t Fraction) Fraction {
				return t.Sub(t.Sam())
			})
```

with:

```go
			// The whole span is mapped relative to the start of this cycle.
			// Using each instant's own Sam() collapses the span, because the
			// end of cycle n has Sam() == n+1 and maps to zero.
			base := cyc.Begin.Sam()
			mappedSpan := cyc.WithTime(func(t Fraction) Fraction {
				return t.Sub(base)
			})
```

Then delete the now-duplicate `base := cyc.Begin.Sam()` that appears a few lines below, inside the hap-shifting loop — it is the same value and the compiler will reject the redeclaration.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/core/ -run SlowCat -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Check for fallout across the suite**

`SlowCat` underpins `Cat`, `<>` alternation and `Arrange`. Run the whole suite before assuming this is safe.

Run: `go test ./... -race -count=1`
Expected: all packages pass. If a test in `internal/core` now fails, read it: a test that pinned the collapsed-span behaviour was encoding the bug and should be updated, but one that fails for another reason means this fix is incomplete.

- [ ] **Step 6: Update the limitations chapter**

Remove the "`SlowCat` drops patterns that do not fill a cycle" section from `docs/tutorial/08-limitations.md`. In `docs/tutorial/03-patterns-in-go.md`, the `SlowCat` row currently warns it is unreliable and redirects to `Every`/`LastOf`; restore it as a normal arrangement tool while keeping the `Every`/`LastOf` advice as the idiom for occasional events.

- [ ] **Step 7: Commit**

```bash
git add internal/core/pattern.go internal/core/pattern_cycle_test.go docs/tutorial/08-limitations.md docs/tutorial/03-patterns-in-go.md
git commit -m "fix(core): give each SlowCat argument a full cycle span"
```

---

### Task 2: `@` elongation must weight a step

**Files:**
- Modify: `internal/mini/mini.go`
- Test: `internal/mini/mini_grammar_test.go`

**Interfaces:**
- Consumes: `core.TimeCatWeighted(pairs ...any) Pattern` — alternating duration and pattern, e.g. `TimeCatWeighted(3, patA, 1, patB)`.
- Produces: unexported `func splitWeight(tok string) (base string, weight float64)` in `internal/mini`. Returns weight `1` when the token carries no `@`.

`"bd@3 sd"` should give `bd` three quarters of the bar and `sd` the last quarter. Today `@` is parsed and routed to `WithSteps`, which changes step metadata but not timing, so the output is identical to `"bd sd"`. The fix moves weight handling up to `parseSequence`, where the relative durations of sibling steps are known.

- [ ] **Step 1: Write the failing test**

```go
// internal/mini/mini_grammar_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// onsets returns each hap's value and start time for one cycle.
func onsets(p core.Pattern) map[string]float64 {
	out := map[string]float64{}
	for _, h := range p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)) {
		if s, ok := h.Value.(string); ok {
			out[s] = h.Part.Begin.Float64()
		}
	}
	return out
}

func TestMiniWeightElongatesAStep(t *testing.T) {
	got := onsets(Mini("bd@3 sd"))
	if len(got) != 2 {
		t.Fatalf("got %d haps, want 2: %v", len(got), got)
	}
	if got["bd"] != 0 {
		t.Errorf("bd starts at %v, want 0", got["bd"])
	}
	if got["sd"] != 0.75 {
		t.Errorf("sd starts at %v, want 0.75 — bd should hold three quarters", got["sd"])
	}
}

func TestMiniWeightIsNotAlwaysThree(t *testing.T) {
	got := onsets(Mini("bd@1 sd@3"))
	if got["sd"] != 0.25 {
		t.Errorf("sd starts at %v, want 0.25", got["sd"])
	}
}

func TestMiniWithoutWeightIsUnchanged(t *testing.T) {
	got := onsets(Mini("bd sd"))
	if got["bd"] != 0 || got["sd"] != 0.5 {
		t.Errorf("unweighted sequence changed: %v", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/mini/ -run MiniWeight -v`
Expected: FAIL — `sd starts at 0.5, want 0.75`

- [ ] **Step 3: Write minimal implementation**

Add to `internal/mini/mini.go`:

```go
// splitWeight separates a trailing @n weight from a token. "bd@3" yields
// ("bd", 3). A token with no weight yields a weight of 1, so callers can treat
// every step uniformly.
func splitWeight(tok string) (string, float64) {
	i := indexAtDepth0(tok, "@")
	if i <= 0 {
		return tok, 1
	}
	w, err := strconv.ParseFloat(strings.TrimSpace(tok[i+1:]), 64)
	if err != nil || w <= 0 {
		return tok[:i], 1
	}
	return tok[:i], w
}
```

Then, in `parseSequence`, replace the final fall-through that builds the sequence:

```go
	pats := make([]core.Pattern, 0, len(tokens))
	for _, tok := range tokens {
		if tok == ".." || tok == "|" || tok == "," {
			continue
		}
		pats = append(pats, parseToken(tok))
	}
	return core.FastCat(pats...)
```

with a weight-aware version:

```go
	// Collect steps with their weights. When every weight is 1 this is exactly
	// FastCat; otherwise the steps divide the cycle in proportion.
	var (
		weighted []any
		weightedAny bool
	)
	for _, tok := range tokens {
		if tok == ".." || tok == "|" || tok == "," {
			continue
		}
		base, w := splitWeight(tok)
		if w != 1 {
			weightedAny = true
		}
		weighted = append(weighted, w, parseToken(base))
	}
	if len(weighted) == 0 {
		return core.Silence()
	}
	if !weightedAny {
		pats := make([]core.Pattern, 0, len(weighted)/2)
		for i := 1; i < len(weighted); i += 2 {
			pats = append(pats, weighted[i].(core.Pattern))
		}
		if len(pats) == 1 {
			return pats[0]
		}
		return core.FastCat(pats...)
	}
	return core.TimeCatWeighted(weighted...)
```

Finally, remove `@` from `parseToken`'s suffix handling so it is not applied twice: in the `containsAtDepth0(tok, "@_!?")` block, drop `@` from the `containsAtDepth0(tok, "@_")` branch, leaving `_` handled as before.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/mini/ -run MiniWeight -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Update the pinned grammar table and docs**

`TestMiniDocumentedSyntax` has the case `{"bd@3 sd", 2}` — the count is still 2, so it should still pass. Run the whole mini package to confirm.

Run: `go test ./internal/mini/ -v`
Expected: PASS

Then in `docs/tutorial/02-mini-notation.md`, move `@` out of the "Two differences from Strudel" section into the main grammar with a worked example, and update the reference table row. Remove the `@` row from the table in `docs/tutorial/08-limitations.md`. In `README.md`, restore the `@n` row to "Elongate / weight".

- [ ] **Step 6: Commit**

```bash
git add internal/mini/mini.go internal/mini/mini_grammar_test.go docs/tutorial/02-mini-notation.md docs/tutorial/08-limitations.md README.md
git commit -m "fix(mini): make @ weight a step's duration"
```

---

### Task 3: `!` replicate must add steps

**Files:**
- Modify: `internal/mini/mini.go`
- Test: `internal/mini/mini_grammar_test.go`

**Interfaces:**
- Consumes: `splitWeight` from Task 2 (the same `parseSequence` loop is edited).
- Produces: unexported `func splitReplicate(tok string) (base string, count int)`. Returns count `1` when the token carries no `!`.

`"bd!3 sd"` should be four equal steps — `bd bd bd sd` — with onsets at 0, 1/4, 1/2 and 3/4. Today it is implemented with `Ply`, which subdivides the token's own slot, giving three kicks inside the first half at 0, 1/6 and 1/3.

- [ ] **Step 1: Write the failing test**

```go
// append to internal/mini/mini_grammar_test.go

func TestMiniReplicateAddsEqualSteps(t *testing.T) {
	p := Mini("bd!3 sd")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("got %d haps, want 4", len(haps))
	}
	want := []float64{0, 0.25, 0.5, 0.75}
	for i, w := range want {
		if got := haps[i].Part.Begin.Float64(); got != w {
			t.Errorf("hap %d begins at %v, want %v", i, got, w)
		}
	}
	if v, _ := haps[3].Value.(string); v != "sd" {
		t.Errorf("last hap is %v, want sd", haps[3].Value)
	}
}

func TestMiniReplicateMatchesWritingItOut(t *testing.T) {
	a := Mini("bd!3 sd").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	b := Mini("bd bd bd sd").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(a) != len(b) {
		t.Fatalf("!3 gave %d haps, writing it out gave %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Part.Begin.Float64() != b[i].Part.Begin.Float64() {
			t.Errorf("hap %d: !3 at %v, written out at %v",
				i, a[i].Part.Begin.Float64(), b[i].Part.Begin.Float64())
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/mini/ -run MiniReplicate -v`
Expected: FAIL — hap 1 begins at 0.1666… rather than 0.25

- [ ] **Step 3: Write minimal implementation**

Add to `internal/mini/mini.go`:

```go
// splitReplicate separates a trailing !n from a token. "bd!3" yields ("bd", 3);
// a bare "bd!" yields ("bd", 2). A token with no ! yields a count of 1.
func splitReplicate(tok string) (string, int) {
	i := indexAtDepth0(tok, "!")
	if i <= 0 {
		return tok, 1
	}
	rest := strings.TrimSpace(tok[i+1:])
	if rest == "" {
		return tok[:i], 2
	}
	n, err := strconv.Atoi(rest)
	if err != nil || n < 1 {
		return tok[:i], 1
	}
	return tok[:i], n
}
```

In the `parseSequence` loop written in Task 2, expand replicated tokens before applying weight:

```go
	for _, tok := range tokens {
		if tok == ".." || tok == "|" || tok == "," {
			continue
		}
		base, reps := splitReplicate(tok)
		base, w := splitWeight(base)
		if w != 1 {
			weightedAny = true
		}
		// A replicated step becomes that many sibling steps, so "bd!3 sd" is
		// four equal quarters rather than three kicks crammed into one slot.
		for i := 0; i < reps; i++ {
			weighted = append(weighted, w, parseToken(base))
		}
	}
```

Then remove the `!` branch from `parseToken` so replication is not applied twice: delete the `if containsAtDepth0(tok, "!")` block that calls `pat.Ply(reps)`, and drop `!` from the `containsAtDepth0(tok, "@_!?")` guard.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/mini/ -run MiniReplicate -v`
Expected: PASS (2 tests)

- [ ] **Step 5: Update the pinned grammar table**

`TestMiniDocumentedSyntax` pins `{"bd!3 sd", 4}` — still 4, so it should pass. Confirm the whole package.

Run: `go test ./internal/mini/ -v`
Expected: PASS

Then update `docs/tutorial/02-mini-notation.md`: remove `!` from "Two differences from Strudel" (which now has only `%` left, so retitle it), and show `"bd!3 sd"` as four equal steps in the grammar section. Remove the `!` row from `docs/tutorial/08-limitations.md`. In `README.md`, restore the `!` row to "Replicate".

- [ ] **Step 6: Commit**

```bash
git add internal/mini/mini.go internal/mini/mini_grammar_test.go docs/tutorial/02-mini-notation.md docs/tutorial/08-limitations.md README.md
git commit -m "fix(mini): make ! add sibling steps instead of subdividing"
```

---

### Task 4: `%n` polymeter steps per cycle

**Files:**
- Modify: `internal/mini/mini.go`
- Test: `internal/mini/mini_grammar_test.go`

**Interfaces:**
- Consumes: `core.SlowCat`, `core.Stack`, `Pattern.FastF`, `core.NewFraction`.
- Produces: no new exported names; the `{` branch of `parseToken` gains real `%n` handling.

**Read this before starting — it changes existing behaviour.** A polymeter plays every layer at the same number of steps per cycle, with each layer cycling through its own elements. `{bd sd, hh hh hh}%4` gives four steps per cycle in both layers. Without an explicit `%`, the step count comes from the first layer.

Today the `{` branch stacks the layers as plain sequences, so `{bd sd, hh hh hh}` yields 5 haps (2 + 3). Correct polymeter with 2 steps per cycle yields 4 (2 + 2), and the second layer's `hh` selection advances across cycles. `TestMiniDocumentedSyntax` currently pins `{"{bd sd, hh hh hh}", 5}` and must be updated to 4 as part of this task — that is a deliberate correction, not a regression.

- [ ] **Step 1: Write the failing test**

```go
// append to internal/mini/mini_grammar_test.go

func TestMiniPolymeterUsesFirstLayerStepCount(t *testing.T) {
	// Two steps per cycle, taken from the first layer's length.
	p := Mini("{bd sd, hh hh hh}")
	got := len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if got != 4 {
		t.Fatalf("got %d haps, want 4 (2 steps per cycle in each of 2 layers)", got)
	}
}

func TestMiniPolymeterExplicitStepCount(t *testing.T) {
	p := Mini("{bd sd, hh hh hh}%4")
	got := len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if got != 8 {
		t.Errorf("got %d haps, want 8 (4 steps per cycle in each of 2 layers)", got)
	}
}

// The point of a polymeter is that a layer whose length does not divide the
// step count lands on different elements from one cycle to the next.
func TestMiniPolymeterLayersDriftAcrossCycles(t *testing.T) {
	p := Mini("{bd sd, a b c}")
	seen := map[string]bool{}
	for c := int64(0); c < 3; c++ {
		for _, h := range p.QueryArc(core.FractionFromInt(c), core.FractionFromInt(c+1)) {
			if s, ok := h.Value.(string); ok {
				seen[s] = true
			}
		}
	}
	for _, want := range []string{"a", "b", "c"} {
		if !seen[want] {
			t.Errorf("the three-element layer never played %q across three cycles: %v", want, seen)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/mini/ -run MiniPolymeter -v`
Expected: FAIL — the first test reports 5 haps, the second 5.

- [ ] **Step 3: Write minimal implementation**

In `internal/mini/mini.go`, replace the body of the `{` branch in `parseToken` with:

```go
	if len(tok) >= 2 && tok[0] == '{' {
		closeIdx := strings.LastIndex(tok, "}")
		if closeIdx > 0 {
			inner := tok[1:closeIdx]
			suffix := strings.TrimSpace(tok[closeIdx+1:])

			// Each comma-separated layer is a list of steps.
			var layers [][]core.Pattern
			for _, part := range splitAtDepth0(inner, ',') {
				toks := splitMiniTokens(strings.TrimSpace(part))
				if len(toks) == 0 {
					continue
				}
				steps := make([]core.Pattern, len(toks))
				for i, s := range toks {
					steps[i] = parseToken(s)
				}
				layers = append(layers, steps)
			}
			if len(layers) == 0 {
				return core.Silence()
			}

			// Steps per cycle: %n if given, otherwise the first layer's length.
			// This is what makes it a polymeter — every layer runs at the same
			// rate, and a layer whose length does not divide that rate lands on
			// different elements each cycle.
			stepsPerCycle := len(layers[0])
			if strings.HasPrefix(suffix, "%") {
				if v, err := strconv.Atoi(strings.TrimSpace(suffix[1:])); err == nil && v > 0 {
					stepsPerCycle = v
				}
				suffix = ""
			}

			pats := make([]core.Pattern, 0, len(layers))
			for _, steps := range layers {
				// SlowCat gives one element per cycle; speeding it up by the
				// step count gives that many elements per cycle, wrapping
				// around the layer's own length.
				pats = append(pats, core.SlowCat(steps...).
					FastF(core.NewFraction(int64(stepsPerCycle), 1)))
			}
			base := pats[0]
			if len(pats) > 1 {
				base = core.Stack(pats...)
			}

			// A leftover *n or /n still applies on top.
			if strings.HasPrefix(suffix, "*") {
				if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
					return base.FastF(core.FractionFromFloat(v))
				}
			} else if strings.HasPrefix(suffix, "/") {
				if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
					return base.SlowF(core.FractionFromFloat(v))
				}
			}
			return base
		}
	}
```

This depends on Task 1: `SlowCat` must give each element a full cycle, or every layer will be silent.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/mini/ -run MiniPolymeter -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Update the pinned grammar table**

In `internal/mini/mini_nested_ops_test.go`, change the `TestMiniDocumentedSyntax` case `{"{bd sd, hh hh hh}", 5}` to `{"{bd sd, hh hh hh}", 4}`. Also check `internal/mini/mini_test.go` for `TestMiniPolymeterCurly` and `TestMiniPolymeterSimple` and update their expectations to true polymeter counts.

Run: `go test ./... -race -count=1`
Expected: all packages pass

- [ ] **Step 6: Update the docs**

In `docs/tutorial/02-mini-notation.md`, rewrite the polymeter section to explain steps-per-cycle and drift, document `%n`, and correct the example's hap count. Remove the `%n` row from the differences table in `docs/tutorial/08-limitations.md` — that section should now be empty, so delete the heading too.

- [ ] **Step 7: Commit**

```bash
git add internal/mini/ docs/tutorial/02-mini-notation.md docs/tutorial/08-limitations.md
git commit -m "fix(mini): implement polymeter steps-per-cycle including %n"
```

---

### Task 5: `Add` must transpose a wrapped pattern

**Files:**
- Create: `internal/core/pattern_arith.go`
- Modify: `internal/core/pattern_composers.go` (or wherever `Add` is defined — find with `grep -rn "func (p Pattern) Add(" internal/core/`)
- Test: `internal/core/pattern_arith_test.go`

**Interfaces:**
- Consumes: `toFloat`, `Reify`, `Fmap`, `AppBoth` — all already in `internal/core`.
- Produces: unexported `func addValues(a, b any) any`. Given two control bags, or a bag and a number, it adds into the bag's numeric field rather than flattening it.

`core.Note(mini.Mini("0 4 7")).Add(12)` currently produces the bare number `12` for every event, because `Add` calls `toFloat` on the whole value and a `map[string]any` has no float form. It should produce `map[note:12]`, `map[note:16]`, `map[note:19]`.

- [ ] **Step 1: Write the failing test**

```go
// internal/core/pattern_arith_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package core

import "testing"

func bagOf(t *testing.T, v any) map[string]any {
	t.Helper()
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("value is %T (%v), want a control bag", v, v)
	}
	return m
}

func TestAddTransposesANotePattern(t *testing.T) {
	p := Note(FastCat(Pure(0), Pure(4), Pure(7))).Add(12)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("got %d haps, want 3", len(haps))
	}
	want := []float64{12, 16, 19}
	for i, w := range want {
		m := bagOf(t, haps[i].Value)
		if got := toFloat(m["note"]); got != w {
			t.Errorf("hap %d note = %v, want %v", i, got, w)
		}
	}
}

func TestAddKeepsOtherControls(t *testing.T) {
	p := Note(Pure(60)).Set(Gain(0.5)).Add(12)
	m := bagOf(t, p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value)
	if toFloat(m["note"]) != 72 {
		t.Errorf("note = %v, want 72", m["note"])
	}
	if toFloat(m["gain"]) != 0.5 {
		t.Errorf("gain = %v, want 0.5 — Add must not disturb other controls", m["gain"])
	}
}

func TestAddOnBareNumbersStillWorks(t *testing.T) {
	haps := Pure(60).Add(12).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if got := toFloat(haps[0].Value); got != 72 {
		t.Errorf("got %v, want 72", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/core/ -run TestAdd -v`
Expected: FAIL — `value is float64 (12), want a control bag`

- [ ] **Step 3: Write minimal implementation**

Create `internal/core/pattern_arith.go`:

```go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Arithmetic that understands control bags.

package core

// numericControls are the fields arithmetic targets, in priority order. A bag
// carries several values and only the pitch-like one should move when a
// pattern is transposed.
var numericControls = []string{"note", "n", "freq"}

// addValues adds b into a. Bare numbers add directly. When either side is a
// control bag, the addition lands on the bag's numeric field and every other
// control is preserved — transposing a note must not discard its gain.
func addValues(a, b any) any {
	am, aIsBag := a.(map[string]any)
	bm, bIsBag := b.(map[string]any)

	switch {
	case !aIsBag && !bIsBag:
		return toFloat(a) + toFloat(b)

	case aIsBag && !bIsBag:
		out := cloneBag(am)
		key := primaryNumeric(am)
		out[key] = toFloat(am[key]) + toFloat(b)
		return out

	case !aIsBag && bIsBag:
		out := cloneBag(bm)
		key := primaryNumeric(bm)
		out[key] = toFloat(a) + toFloat(bm[key])
		return out

	default:
		out := cloneBag(am)
		for k, v := range bm {
			if existing, ok := out[k]; ok {
				out[k] = toFloat(existing) + toFloat(v)
			} else {
				out[k] = v
			}
		}
		return out
	}
}

// primaryNumeric picks the field arithmetic should target, defaulting to note
// so that adding to a bag without one still produces something musical.
func primaryNumeric(m map[string]any) string {
	for _, k := range numericControls {
		if _, ok := m[k]; ok {
			return k
		}
	}
	return "note"
}

func cloneBag(m map[string]any) map[string]any {
	out := make(map[string]any, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}
```

Then change `Add` to use it. Replace its body with:

```go
func (p Pattern) Add(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return addValues(a, b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/core/ -run TestAdd -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Check for fallout**

`Add` is widely aliased. Run the full suite and read any failure carefully — a test asserting the old flattening behaviour was encoding the bug.

Run: `go test ./... -race -count=1`
Expected: all packages pass

- [ ] **Step 6: Update the docs**

In `docs/tutorial/05-transformations.md`, the "Pitch" section tells readers to add before wrapping because `Add` on a wrapped pattern does not work. Rewrite it to show `core.Note(mini.Mini("c3 e3 g3")).Add(12)` working directly. Remove the "`Add` does not transpose a wrapped pattern" section from `docs/tutorial/08-limitations.md`. Update the `TestDocumentedIdioms/transpose_before_wrapping` subtest in `examples/examples_test.go` to assert the new behaviour, keeping a case for bare numbers.

- [ ] **Step 7: Commit**

```bash
git add internal/core/ docs/tutorial/05-transformations.md docs/tutorial/08-limitations.md examples/examples_test.go
git commit -m "fix(core): make Add transpose control bags instead of flattening them"
```

---

## Verification after all tasks

```bash
go vet ./...
go test ./... -race -count=1
for d in house chicago-house techno minimal-dubstep maximal-dubstep drum-and-bass electronica trance mytrack; do
  go run ./examples/$d
done
```

The templates use `@` nowhere and `!` nowhere, so their event counts should be unchanged. If any moves, find out why before updating the documented numbers in `docs/tutorial/templates/`. `minimal-dubstep` and `maximal-dubstep` both use `LastOf`, which Task 1 touches indirectly through `SlowCat`; confirm their fills still fire on the fourth bar.

After Tasks 1–4, `docs/tutorial/08-limitations.md` should have lost the `SlowCat` section and the whole "Mini-notation differences from Strudel" table. Re-read that chapter end to end and delete any now-stale framing.
