// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Seventh batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SelectCarry(t *testing.T) {
	// select-like via FilterValues / When
	p := FastCat(Pure(1), Pure(2), Pure(3)).FilterValues(func(v any) bool {
		if n, ok := v.(int); ok {
			return n%2 == 1
		}
		return false
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("filter odd empty")
	}
	// carry-like: Stack + Filter?
	p2 := Pure(map[string]any{"a": 1, "b": 2}).Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m["c"] = 3
			return m
		}
		return v
	})
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fmap carry empty")
	}
}

func TestMJS_FastSlowEvery(t *testing.T) {
	// fast with pattern of factors (pattern.test mjs: fast can take pattern)
	p := Pure("a").Fast(Pure(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fast pattern empty")
	}
	p2 := Pure("a").Slow(FractionFromInt(2))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("slow pattern empty")
	}
	// slowAny with pattern
	p2b := Pure("a").SlowAny(Pure(2))
	if len(p2b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("slowAny pattern empty")
	}
	// every 3
	p3 := Pure("a").Every(3, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("every 3 empty")
	}
}

func TestMJS_ZoomCompressFocus(t *testing.T) {
	p := Pure("a").Zoom(FractionFromFloat(0.0), FractionFromFloat(0.5))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("zoom 0-0.5 empty")
	}
	p2 := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress 0.25-0.75 empty")
	}
	p3 := Pure("a").FastGap(FractionFromFloat(0.5))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fastGap 0.5 empty")
	}
}
