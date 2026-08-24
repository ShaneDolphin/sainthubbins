// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Eleventh batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SuperimposeLayer(t *testing.T) {
	p := Pure("a").Superimpose(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("superimpose expected >=2 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := Pure("a").Layer(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) }, func(p Pattern) Pattern { return p.Rev() })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("layer empty")
	}
}

func TestMJS_DiscreteOnsets(t *testing.T) {
	p := Pure("a").DiscreteOnly()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("discreteOnly empty")
	}
	p2 := Pure("a").OnsetsOnly()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("onsetsOnly empty")
	}
	p3 := FastCat(Pure("a"), Pure("b")).FilterHaps(func(h Hap) bool { return h.Value == "a" })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("filterHaps empty")
	}
}

func TestMJS_RollWithChildren(t *testing.T) {
	p := Pure("a").Chunk(3, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunk empty")
	}
	p2 := Pure("a").Chunk(2, func(p Pattern) Pattern { return p.Rev() })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunk rev empty")
	}
}
