// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Regression tests: cycle-dependent combinators must decide per cycle, not once
// per query span. The offline renderer queries many cycles in a single call, so
// a combinator that reads the cycle from the span start silently applies its
// transformation to the whole render.

package core

import (
	"sort"
	"testing"
)

// spanVsSplit returns the hap count for one query covering n cycles, and the
// sum of n single-cycle queries. They must agree.
func spanVsSplit(p Pattern, n int64) (whole int, split int) {
	whole = len(p.QueryArc(FractionFromInt(0), FractionFromInt(n)))
	for c := int64(0); c < n; c++ {
		split += len(p.QueryArc(FractionFromInt(c), FractionFromInt(c+1)))
	}
	return whole, split
}

func fourBeats() Pattern {
	return FastCat(Pure("bd"), Pure("bd"), Pure("bd"), Pure("bd"))
}

func TestEveryIsPerCycleUnderWideQuery(t *testing.T) {
	silence := func(Pattern) Pattern { return Silence() }

	// Every(2, silence) must silence alternate cycles, not the whole query.
	p := fourBeats().Every(2, silence)
	whole, split := spanVsSplit(p, 8)
	if whole != split {
		t.Errorf("Every(2, silence): one 8-cycle query gave %d haps, eight 1-cycle queries gave %d", whole, split)
	}
	if whole == 0 {
		t.Errorf("Every(2, silence) silenced every cycle under a wide query")
	}
	if want := 16; whole != want {
		t.Errorf("Every(2, silence): got %d haps over 8 cycles, want %d", whole, want)
	}
}

func TestEveryFiresOnlyOnMatchingCycles(t *testing.T) {
	// Starting from silence, the function must apply on cycles 0 and 4 only.
	p := Silence().Every(4, func(Pattern) Pattern { return fourBeats() })
	for c := int64(0); c < 8; c++ {
		got := len(p.QueryArc(FractionFromInt(c), FractionFromInt(c+1)))
		want := 0
		if c%4 == 0 {
			want = 4
		}
		if got != want {
			t.Errorf("cycle %d: got %d haps, want %d", c, got, want)
		}
	}
	if whole, split := spanVsSplit(p, 8); whole != split {
		t.Errorf("wide query gave %d haps, per-cycle sum %d", whole, split)
	}
}

func TestLastOfFiresOnFinalCycleOfEachGroup(t *testing.T) {
	// LastOf(4, ...) is how a fill lands at the end of a four-bar phrase.
	p := Silence().LastOf(4, func(Pattern) Pattern { return fourBeats() })
	for c := int64(0); c < 8; c++ {
		got := len(p.QueryArc(FractionFromInt(c), FractionFromInt(c+1)))
		want := 0
		if c%4 == 3 {
			want = 4
		}
		if got != want {
			t.Errorf("cycle %d: got %d haps, want %d", c, got, want)
		}
	}
	if whole, split := spanVsSplit(p, 8); whole != split {
		t.Errorf("wide query gave %d haps, per-cycle sum %d", whole, split)
	}
}

// Transformations that preserve hap count still have to land on the right
// cycles, which a count-only check cannot detect.
func TestEveryAppliesToCorrectCycleOrdering(t *testing.T) {
	base := FastCat(Pure("a"), Pure("b"))
	p := base.Every(2, func(q Pattern) Pattern { return q.Rev() })

	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps) != 4 {
		t.Fatalf("want 4 haps over 2 cycles, got %d", len(haps))
	}
	// Haps come back in the pattern's own order, so sort by onset before
	// comparing: it is the times that Rev changes, not the traversal order.
	sort.Slice(haps, func(i, j int) bool {
		return haps[i].Part.Begin.Float64() < haps[j].Part.Begin.Float64()
	})
	// Cycle 0 is reversed (b then a); cycle 1 is untouched (a then b).
	want := []string{"b", "a", "a", "b"}
	for i, w := range want {
		got, ok := haps[i].Value.(string)
		if !ok || got != w {
			t.Errorf("position %d (onset %v): got %v, want %q",
				i, haps[i].Part.Begin.String(), haps[i].Value, w)
		}
	}
}
