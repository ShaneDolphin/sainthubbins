// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-eighth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_EveryOffWhen2(t *testing.T) {
	p := Pure("a").Every(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("every 2 c0 expected 2 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("off 0.25 expected >=2")
	}
	p3 := Pure("a").When(Pure(false), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("when false expected 1 got %d", len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_FilterHapsValues(t *testing.T) {
	p := FastCat(Pure(1), Pure(2), Pure(3)).FilterHaps(func(h Hap) bool {
		if n, ok := h.Value.(int); ok {
			return n > 1
		}
		return false
	})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("filterHaps >1 expected 2 got %d", len(haps))
	}
	p2 := FastCat(Pure(1), Pure(2)).FilterValues(func(v any) bool {
		if n, ok := v.(int); ok {
			return n == 1
		}
		return false
	})
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("filterValues ==1 expected 1 got %d", len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_SpliceFit(t *testing.T) {
	s := Slice(2, Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("slice empty")
	}
	// Fit-like via Fast: use Fast with Pattern
	fc := FastCat(Pure("a"), Pure("b")).Fast(Pure(2))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("fast 2 on fastcat 2 expected 4 got %d", len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}
