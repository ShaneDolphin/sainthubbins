// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Sixteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FastSlowCompress(t *testing.T) {
	p := Pure("a").FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fast 2 expected 2 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := Pure("a").Slow(FractionFromInt(2))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("slow 2 expected 1 got %d", len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p3 := Pure("a").Compress(FractionFromFloat(0.0), FractionFromFloat(0.5))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress 0-0.5 empty")
	}
}

func TestMJS_InsideOutside(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Inside(2, func(p Pattern) Pattern { return p.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 3 {
		t.Fatalf("inside rev >=3 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := FastCat(Pure("a"), Pure("b")).Outside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("outside empty")
	}
}

func TestMJS_FilterValues(t *testing.T) {
	p := FastCat(Pure(true), Pure(false), Pure(true)).FilterValues(func(v any) bool {
		if b, ok := v.(bool); ok {
			return b
		}
		return false
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("filterValues empty")
	}
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	for _, h := range haps {
		if h.Value != true {
			t.Fatalf("filterValues should be true got %v", h.Value)
		}
	}
}
