// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 38th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SelectByFraction(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	q := p.Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("sometimes empty")
	}
	// Often = SometimesBy 0.75, Rarely = SometimesBy 0.25 via SometimesBy
	r := p.SometimesBy(0.75, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps2 := r.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("often (0.75) empty")
	}
	s := p.SometimesBy(0.25, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps3 := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) == 0 {
		t.Fatalf("rarely (0.25) empty")
	}
}

func TestMJS_PatStack(t *testing.T) {
	a := Pure("a")
	b := Pure("b")
	c := Pure("c")
	stacked := Stack(a, b, c)
	haps := stacked.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("stack 3 expected 3 got %d", len(haps))
	}
	cat := Cat(a, b, c)
	haps2 := cat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Cat = SlowCat: only first pat per cycle, so 1 hap over 1 cycle
	if len(haps2) != 1 {
		t.Fatalf("cat 3 over 1 cycle expected 1 (slowcat) got %d", len(haps2))
	}
	// FastCat gives 3
	fast := FastCat(a, b, c)
	haps3 := fast.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) != 3 {
		t.Fatalf("fastcat 3 expected 3 got %d", len(haps3))
	}
	// Cat over 3 cycles should give 3 (one per cycle)
	haps4 := cat.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps4) != 3 {
		t.Fatalf("cat 3 over 3 cycles expected 3 got %d", len(haps4))
	}
}

func TestMJS_PatternDecay(t *testing.T) {
	// Decay is createParam("decay") control, not Pattern method — use Pure("bd").Set(Decay(0.1)) or Stack
	p := Pure("bd").Set(Decay(0.1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("decay control via Set empty")
	}
	plain := Pure("bd").QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != len(plain) {
		t.Fatalf("decay len vs plain %d vs %d", len(haps), len(plain))
	}
	// Decay value bag merged?
	v, ok := haps[0].Value.(map[string]any)
	if !ok {
		t.Fatalf("decay value not map %T", haps[0].Value)
	}
	if _, exists := v["decay"]; !exists {
		t.Fatalf("decay key missing in %v", v)
	}
}
