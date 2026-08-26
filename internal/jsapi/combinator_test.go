// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

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

// TestCatAlternates above only queries one cycle, which cannot distinguish
// concatenation (bd, then sd, then bd again) from repetition (bd every
// cycle) or from stack (bd and sd both on cycle 0) — all three produce
// exactly one hap in [0,1). Query several cycles and check which sample
// plays on each one.
func TestCatOverMultipleCyclesAlternates(t *testing.T) {
	p, err := Evaluate(`cat(s("bd"), s("sd"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	want := []string{"bd", "sd", "bd", "sd"}
	for i, wantSample := range want {
		haps := p.QueryArc(core.FractionFromInt(int64(i)), core.FractionFromInt(int64(i+1)))
		if len(haps) != 1 {
			t.Fatalf("cycle %d: got %d haps, want 1", i, len(haps))
		}
		m, ok := haps[0].Value.(map[string]any)
		if !ok || m["s"] != wantSample {
			t.Errorf("cycle %d: got %v, want s:%s", i, haps[0].Value, wantSample)
		}
	}
}

// hapKeys is defined in transform_test.go (same package) and compares haps
// by (start, end, value) rather than struct identity.

// TestCatWideQueryMatchesStitchedCycles is the property CLAUDE.md calls out
// for any cycle-dependent combinator: the renderer issues exactly one wide
// QueryArc(0, cycles) call, so a combinator that reads the cycle number
// without splitting is only correct by accident, when a test happens to
// query one cycle at a time. core.Cat is an alias for core.SlowCat, which
// already iterates state.Span.SpanCycles() itself (the self-splitting,
// compliant shape CLAUDE.md describes) — this test is the check that the
// jsapi wiring around it doesn't undo that guarantee.
func TestCatWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`cat(s("bd"), s("sd"), s("hh"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	assertWideQueryMatchesStitchedCycles(t, p, 9)
}

// slowcat is the same combinator cat is an alias for; test it directly
// too since both names are exposed as separate globals.
func TestSlowcatWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`slowcat(s("bd"), s("sd"), s("hh"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	assertWideQueryMatchesStitchedCycles(t, p, 9)
}

func TestFastcatWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`fastcat(s("bd"), s("sd"), s("hh"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	assertWideQueryMatchesStitchedCycles(t, p, 6)
}

func TestSequenceWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`sequence(s("bd"), s("sd"), s("hh"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	assertWideQueryMatchesStitchedCycles(t, p, 6)
}

func TestStackWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`stack(s("bd*2"), s("hh*3"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	assertWideQueryMatchesStitchedCycles(t, p, 5)
}

// assertWideQueryMatchesStitchedCycles queries p once over [0, cycles) and
// compares that to `cycles` separate one-cycle queries stitched together.
// A mismatch means the combinator (or the jsapi wiring around it) is only
// correct when queried one cycle at a time — silent in the cycle-by-cycle
// test suite, audible in the renderer's one wide QueryArc call.
//
// The comparison is order-independent (both key slices are sorted before
// comparing). core.Stack queries each sub-pattern across the *entire*
// requested span before moving to the next one, so a wide query groups
// haps by layer (all of pat[0]'s haps, then all of pat[1]'s) while the
// stitched cycle-by-cycle queries group them by cycle (pat[0] and pat[1]'s
// haps for cycle 0, then both for cycle 1, ...) — same haps, different
// interleaving, and that interleaving isn't a property either caller of
// this helper claims. Sorting doesn't hide a genuine content bug: hapKey
// already encodes (begin, end, value), so a combinator that produced the
// wrong sample at a given time slot would still produce a different sorted
// multiset, not a merely reordered one.
func assertWideQueryMatchesStitchedCycles(t *testing.T, p core.Pattern, cycles int64) {
	t.Helper()
	wide := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles))
	var stitched []core.Hap
	for i := int64(0); i < cycles; i++ {
		stitched = append(stitched, p.QueryArc(core.FractionFromInt(i), core.FractionFromInt(i+1))...)
	}
	wantWide, wantStitched := hapKeys(wide), hapKeys(stitched)
	sortHapKeys(wantWide)
	sortHapKeys(wantStitched)
	if !reflect.DeepEqual(wantWide, wantStitched) {
		t.Fatalf("wide query over %d cycles disagrees with %d stitched one-cycle queries:\nwide:     %v\nstitched: %v",
			cycles, cycles, wantWide, wantStitched)
	}
	if len(wide) == 0 {
		t.Fatalf("wide query produced no haps at all — this test can't distinguish correct splitting from a pattern that produces nothing")
	}
}

func sortHapKeys(keys []hapKey) {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].begin != keys[j].begin {
			return keys[i].begin < keys[j].begin
		}
		if keys[i].end != keys[j].end {
			return keys[i].end < keys[j].end
		}
		return keys[i].value < keys[j].value
	})
}

// TestVariadicNoArgumentsIsSilence documents a deliberate choice: calling a
// combinator with zero pattern arguments is not a caller error. Each
// core.Stack/Cat/SlowCat/FastCat/Sequence already special-cases zero
// patterns and returns Silence() — a coherent "combine nothing" identity —
// so the jsapi binding passes an empty argument list straight through
// rather than rejecting it.
func TestVariadicNoArgumentsIsSilence(t *testing.T) {
	for _, name := range []string{"stack", "cat", "slowcat", "fastcat", "sequence"} {
		t.Run(name, func(t *testing.T) {
			if got := countHaps(t, name+"()"); got != 0 {
				t.Errorf("%s(): got %d haps, want 0 (silence)", name, got)
			}
		})
	}
}

// TestVariadicAcceptsNumbersAndStrings decides the two non-pattern
// argument shapes that are coherent: a bare string is mini-notation
// (matching every other place in this API that accepts a Pattern argument
// — s(), the control setters), and a bare number becomes a constant
// Pattern via core.Pure, matching toPattern's number branch. Mixing a real
// pattern with a bare number in the same call must also work.
func TestVariadicAcceptsNumbersAndStrings(t *testing.T) {
	if got := countHaps(t, `stack(42)`); got != 1 {
		t.Errorf("stack(42): got %d haps, want 1", got)
	}
	if got := countHaps(t, `stack("bd sd")`); got != 2 {
		t.Errorf(`stack("bd sd"): got %d haps, want 2`, got)
	}
	if got := countHaps(t, `stack(s("bd"), 42)`); got != 2 {
		t.Errorf(`stack(s("bd"), 42): got %d haps, want 2`, got)
	}
}

// TestVariadicRejectsInvalidArguments is the other half of that same
// decision: null, undefined, a plain object literal and a function are not
// coherent pattern arguments. Silently treating any of them as an empty
// layer (the old toPattern behavior, before this task) is exactly the
// silent-no-op shape this whole plan exists to remove — a stray undefined
// variable or object literal in stack(...) must surface as a JS error, not
// vanish.
func TestVariadicRejectsInvalidArguments(t *testing.T) {
	cases := map[string]string{
		"stack with null":             `stack(s("bd"), null)`,
		"stack with undefined":        `stack(s("bd"), undefined)`,
		"stack with a plain object":   `stack(s("bd"), {})`,
		"stack with a function":       `stack(s("bd"), function(){})`,
		"cat with null":               `cat(null, s("bd"))`,
		"cat with undefined":          `cat(s("bd"), undefined)`,
		"cat with a plain object":     `cat({}, s("bd"))`,
		"cat with a function":         `cat(s("bd"), function(){})`,
		"slowcat with null":           `slowcat(null, s("bd"))`,
		"fastcat with a plain object": `fastcat(s("bd"), {})`,
		"sequence with undefined":     `sequence(s("bd"), undefined)`,
		// A bare boolean is also not one of the three accepted shapes
		// (pattern, string, number); accepting it via ToFloat-style
		// coercion the way requireFiniteNumber does for chained numeric
		// methods would be a different, deliberate design choice this
		// task doesn't need to make, so it is rejected like any other
		// unrecognized type.
		"stack with a boolean": `stack(s("bd"), true)`,
	}
	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Evaluate(code); err == nil {
				t.Errorf("Evaluate(%q): want an error, got nil", code)
			}
		})
	}
}

// TestVariadicDeeplyNested confirms toPattern's *jsPattern branch unwraps
// an arbitrarily nested combinator result the same as any other wrapped
// pattern — stack(stack(stack(p))) must behave exactly like p, not lose
// events or panic walking the wrapper chain.
func TestVariadicDeeplyNested(t *testing.T) {
	if got := countHaps(t, `stack(stack(stack(s("bd sd"))))`); got != 2 {
		t.Errorf("got %d haps, want 2", got)
	}
}

// TestVariadicLargeArgumentList confirms a combinator with a large number
// of pattern arguments neither panics nor silently drops any of them.
func TestVariadicLargeArgumentList(t *testing.T) {
	const n = 500
	args := make([]string, n)
	for i := range args {
		args[i] = `s("bd")`
	}
	code := fmt.Sprintf("stack(%s)", strings.Join(args, ", "))
	if got := countHaps(t, code); got != n {
		t.Errorf("stack of %d patterns: got %d haps, want %d", n, got, n)
	}
}

// TestMiniRejectsNonStringArguments: mini() is an explicit escape hatch to
// mini-notation text, not a general-purpose stringify-then-parse helper.
// call.Argument(0).String() would happily coerce a missing argument to the
// literal string "undefined" and a wrapped pattern to "[object Object]",
// parse that coerced text as mini-notation, and hand back a pattern that
// silently plays a sample named after the coercion rather than erroring —
// exactly the silent-wrong-result shape this plan exists to remove.
func TestMiniRejectsNonStringArguments(t *testing.T) {
	cases := map[string]string{
		"no argument":        `mini()`,
		"a number":           `mini(42)`,
		"a pattern argument": `mini(s("bd"))`,
		"null":               `mini(null)`,
		"undefined":          `mini(undefined)`,
	}
	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Evaluate(code); err == nil {
				t.Errorf("Evaluate(%q): want an error, got nil", code)
			}
		})
	}
}
