// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Regression tests: operators must respect bracket depth.

package mini

import (
	"fmt"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// values returns the distinct hap values of one cycle.
func values(p core.Pattern) []string {
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	out := make([]string, 0, len(haps))
	for _, h := range haps {
		out = append(out, fmt.Sprint(h.Value))
	}
	return out
}

// noGarbage asserts no atom carries a stray bracket or operator character,
// which is the signature of a token being split mid-bracket. Only string atoms
// are inspected: a control bag renders as map[n:1 s:bd] and legitimately
// contains brackets.
func noGarbage(t *testing.T, code string, p core.Pattern) {
	t.Helper()
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	for _, h := range haps {
		v, ok := h.Value.(string)
		if !ok {
			continue
		}
		for _, bad := range []string{"[", "]", "<", ">", "{", "}", "|", ",", "@", "!", "%", ".."} {
			if contains(v, bad) {
				t.Fatalf("%q produced garbage atom %q (contains %q)", code, v, bad)
			}
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// Bug 1: comma-stacking inside square brackets.
func TestMiniBracketCommaStack(t *testing.T) {
	code := "[bd*4, hh*8]"
	p := Mini(code)
	noGarbage(t, code, p)
	// 4 kicks + 8 hats stacked = 12 haps in one cycle
	if got := len(values(p)); got != 12 {
		t.Fatalf("%q expected 12 haps (4 bd + 8 hh), got %d: %v", code, got, values(p))
	}
}

// Bug 2: random choice inside square brackets.
func TestMiniBracketChoice(t *testing.T) {
	code := "[bd|sd|hh]"
	p := Mini(code)
	noGarbage(t, code, p)
	if got := len(values(p)); got != 1 {
		t.Fatalf("%q expected 1 hap, got %d: %v", code, got, values(p))
	}
	v := values(p)[0]
	if v != "bd" && v != "sd" && v != "hh" {
		t.Fatalf("%q chose %q, want one of bd/sd/hh", code, v)
	}
}

// Operators inside brackets must not split the bracket group.
func TestMiniOperatorInsideBrackets(t *testing.T) {
	cases := []struct {
		code string
		want int
	}{
		{"[bd*4 sd]*2", 10},      // (4 bd + 1 sd) * 2
		{"[bd*2 sd]", 3},         // 2 bd + 1 sd
		{"[bd sd]*2", 4},         // already worked — guard against regression
		{"[bd(3,8) hh]", 4},      // euclid inside brackets: 3 bd + 1 hh
		{"[bd*4, hh*8, cp]", 13}, // three-way stack
	}
	for _, tc := range cases {
		p := Mini(tc.code)
		noGarbage(t, tc.code, p)
		if got := len(values(p)); got != tc.want {
			t.Errorf("%q expected %d haps, got %d: %v", tc.code, tc.want, got, values(p))
		}
	}
}

// Nested stacking one level deeper.
func TestMiniNestedStack(t *testing.T) {
	code := "[[bd sd], hh*4]"
	p := Mini(code)
	noGarbage(t, code, p)
	if got := len(values(p)); got != 6 {
		t.Fatalf("%q expected 6 haps (2 + 4), got %d: %v", code, got, values(p))
	}
}

// Top-level comma stacking without brackets.
func TestMiniTopLevelCommaStack(t *testing.T) {
	code := "bd*4, hh*8"
	p := Mini(code)
	noGarbage(t, code, p)
	if got := len(values(p)); got != 12 {
		t.Fatalf("%q expected 12 haps, got %d: %v", code, got, values(p))
	}
}

// A stack must be simultaneous, not sequential: counting haps is not enough,
// so assert both layers start at cycle position 0 and each spans a full cycle.
func TestMiniStackIsSimultaneous(t *testing.T) {
	for _, code := range []string{"[bd, hh]", "bd, hh"} {
		p := Mini(code)
		haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(haps) != 2 {
			t.Fatalf("%q expected 2 haps, got %d", code, len(haps))
		}
		for _, h := range haps {
			if h.Whole == nil {
				t.Fatalf("%q hap %v has no whole", code, h.Value)
			}
			if h.Whole.Begin.Float64() != 0 {
				t.Errorf("%q layer %v starts at %v, want 0 (sequence, not stack)",
					code, h.Value, h.Whole.Begin.Float64())
			}
			if h.Whole.End.Float64() != 1 {
				t.Errorf("%q layer %v ends at %v, want 1 (sequence, not stack)",
					code, h.Value, h.Whole.End.Float64())
			}
		}
	}
}

// Guard the documented grammar so the tutorial cannot drift from the engine.
func TestMiniDocumentedSyntax(t *testing.T) {
	cases := []struct {
		code string
		want int
	}{
		{"bd sd", 2},
		{"bd ~ sd ~", 2},
		{"bd*4", 4},
		{"bd/2", 1},
		{"[bd sd]*2", 4},
		{"bd(3,8)", 3},
		{"<bd sd>", 1},
		{"bd!3 sd", 4},
		{"bd@3 sd", 2},
		{"{bd sd, hh hh hh}", 4},
		{"[bd*4, hh*8]", 12},
		{"c3 e3 g3", 3},
		{"bd:1 sd:2", 2},
		{"{0 .. 3}", 4},
	}
	for _, tc := range cases {
		p := Mini(tc.code)
		noGarbage(t, tc.code, p)
		if got := len(values(p)); got != tc.want {
			t.Errorf("%q expected %d haps, got %d: %v", tc.code, tc.want, got, values(p))
		}
	}
}

// TestMiniSplitStepBaseReattachesRemainderAfterDigits is the white-box half
// of Fix 2: splitWeight/splitReplicate must consume only the leading numeric
// run after "@"/"!" and reattach whatever comes after it (a "?" degrade,
// the other operator) to the base, rather than handing the whole remainder
// to ParseFloat/Atoi and discarding it wholesale when that fails.
func TestMiniSplitStepBaseReattachesRemainderAfterDigits(t *testing.T) {
	cases := []struct {
		tok        string
		wantBase   string
		wantReps   int
		wantWeight float64
	}{
		{"bd@3?", "bd?", 1, 3},
		{"hh*8@2?", "hh*8?", 1, 2},
		{"bd!2?", "bd?", 2, 1},
		{"bd!2@3", "bd", 2, 3}, // deferred case: "!" consumed first, then "@3" strips into weight 3
	}
	for _, tc := range cases {
		base, reps, weight := splitStepBase(tc.tok)
		if base != tc.wantBase || reps != tc.wantReps || weight != tc.wantWeight {
			t.Errorf("splitStepBase(%q) = (%q, %d, %v), want (%q, %d, %v)",
				tc.tok, base, reps, weight, tc.wantBase, tc.wantReps, tc.wantWeight)
		}
	}
}

// TestMiniDegradeSurvivesWeight guards Fix 2's observable behavior for "?"
// after "@": a single cycle can't tell "always degrades" from "never
// degrades" from "degrades half the time" apart, so this asserts over many
// cycles that the hap count actually thins out relative to the
// never-degraded baseline.
//
// The "!" replicate sibling of this case ("bd!2?") is intentionally not
// asserted here — see the doc comment on TestMiniReplicateDegradeMaskedByFastCatCycleBug.
func TestMiniDegradeSurvivesWeight(t *testing.T) {
	const cycles = 32
	count := func(code string) int {
		return len(Mini(code).QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles)))
	}
	base := count("bd@3")
	if base != cycles {
		t.Fatalf("bd@3 over %d cycles: got %d haps, want %d (sanity baseline)", cycles, base, cycles)
	}
	degraded := count("bd@3?")
	if degraded == base {
		t.Errorf("%q over %d cycles: got %d haps, same as the never-degraded baseline %d — \"?\" is being swallowed by the \"@\" weight",
			"bd@3?", cycles, degraded, base)
	}
}

// TestMiniReplicateDegradeMaskedByFastCatCycleBug documents a finding from
// verifying Fix 2, not a Fix 2 bug: splitStepBase is provably correct for
// "bd!2?" (see TestMiniSplitStepBaseReattachesRemainderAfterDigits — it
// resolves to base "bd?", reps 2, exactly as intended), but the end-to-end
// hap count never thins out, because core.FastCat's sub-pattern queries
// always carry local span 0/1->1/1 regardless of which outer cycle is being
// queried (verified directly against core.FastCat + DegradeBy, with no
// mini-notation involved). DegradeBy's per-hap coin flip is keyed on the
// hap's begin time, so every replica inside any multi-step sequence gets the
// exact same, cycle-invariant decision forever — "bd@3?" only degrades
// because a lone step bypasses FastCat/TimeCat entirely and is queried at
// the true absolute time.
//
// This is a pre-existing engine defect, reproducible with no mini-notation
// or Fix 2 changes involved, and out of scope for this fix wave — flagged
// here instead of silently narrowing the guard.
func TestMiniReplicateDegradeMaskedByFastCatCycleBug(t *testing.T) {
	const cycles = 32
	base := len(Mini("bd!2").QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles)))
	degraded := len(Mini("bd!2?").QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles)))
	if degraded != base {
		t.Fatalf(`"bd!2?" over %d cycles: got %d haps, want %d (the never-degraded baseline) — `+
			"if this fails, the FastCat cycle-position defect this test documents has been fixed; "+
			"replace this test with a real degrade assertion like TestMiniDegradeSurvivesWeight's",
			cycles, degraded, base)
	}
}

// A group must be closed by its own bracket. Accepting any closer turns a typo
// into a plausible but different pattern instead of an obvious literal.
func TestMiniMismatchedBracketsAreNotGroups(t *testing.T) {
	for _, code := range []string{"<bd]", "[bd)", "{bd]", "[bd>"} {
		haps := Mini(code).QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(haps) != 1 {
			t.Fatalf("%q: got %d haps, want 1 literal", code, len(haps))
		}
		got, ok := haps[0].Value.(string)
		if !ok || got != code {
			t.Errorf("%q parsed as %v; a mismatched pair should stay a literal atom", code, haps[0].Value)
		}
	}
	// The matched forms must still parse as groups.
	for _, code := range []string{"<bd sd>", "[bd sd]"} {
		if n := len(values(Mini(code))); n == 0 {
			t.Errorf("%q should still parse as a group, got %d haps", code, n)
		}
	}
}

// Depth-0 commas inside <> stack alternations, the same way they stack
// sequences inside []. Previously the comma rode along on a token and produced
// an empty-valued hap.
func TestMiniAngleBracketCommaStacks(t *testing.T) {
	code := "<bd sd, hh cp>"
	p := Mini(code)
	noGarbage(t, code, p)

	// Cycle 0 takes the first item of each layer, cycle 1 the second.
	for cycle, want := range map[int64][]string{0: {"bd", "hh"}, 1: {"sd", "cp"}} {
		haps := p.QueryArc(core.FractionFromInt(cycle), core.FractionFromInt(cycle+1))
		if len(haps) != 2 {
			t.Fatalf("%q cycle %d: got %d haps, want 2", code, cycle, len(haps))
		}
		seen := map[string]bool{}
		for _, h := range haps {
			s, ok := h.Value.(string)
			if !ok || s == "" {
				t.Fatalf("%q cycle %d: empty or non-string value %v", code, cycle, h.Value)
			}
			seen[s] = true
		}
		for _, w := range want {
			if !seen[w] {
				t.Errorf("%q cycle %d: missing %q, got %v", code, cycle, w, seen)
			}
		}
	}

	// A choice operator alongside the comma must not leave a stray empty hap.
	for _, h := range Mini("<bd|sd, hh|cp>").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)) {
		if s, ok := h.Value.(string); ok && s == "" {
			t.Errorf(`"<bd|sd, hh|cp>" produced an empty-valued hap`)
		}
	}
}
