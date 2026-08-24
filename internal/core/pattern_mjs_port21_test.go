// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-first batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SqueezeJoinInnerJoin(t *testing.T) {
	p := Pure(Pure("a")).SqueezeJoin()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("squeezeJoin empty")
	}
	p2 := Pure(Pure("a")).InnerJoin()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("innerJoin empty")
	}
	p3 := Pure(Pure("a")).Fmap(func(v any) any { return v }).SqueezeJoin()
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("squeeze via fmap empty")
	}
}

func TestMJS_OutSqueeze(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b")).OpOut(Pure(1), func(a any) func(any) any { return func(b any) any { return a } })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("OpOut empty")
	}
	p2 := FastCat(Pure("a"), Pure("b")).OpSqueeze(Pure(2), func(a any) func(any) any { return func(b any) any { return a } })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("OpSqueeze empty")
	}
}

func TestMJS_PureFmapLog(t *testing.T) {
	p := Pure(3).Fmap(func(v any) any { return v.(int) + 1 })
	if v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(int); v != 4 {
		t.Fatalf("fmap 3+1 expected 4 got %v", v)
	}
	// Log-like via Fmap for debugging
	p2 := Pure("a").Fmap(func(v any) any { return v })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fmap id empty")
	}
}
