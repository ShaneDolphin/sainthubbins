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
		for _, bad := range []string{"[", "]", "<", ">", "{", "}", "|", ","} {
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
	}
	for _, tc := range cases {
		p := Mini(tc.code)
		noGarbage(t, tc.code, p)
		if got := len(values(p)); got != tc.want {
			t.Errorf("%q expected %d haps, got %d: %v", tc.code, tc.want, got, values(p))
		}
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
