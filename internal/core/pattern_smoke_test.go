// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestPureFirstCycle(t *testing.T) {
	p := Pure("hello")
	haps := p.FirstCycle()
	if len(haps) != 1 {
		t.Fatalf("pure first cycle expected 1 hap, got %d", len(haps))
	}
	if haps[0].Value != "hello" {
		t.Fatalf("expected hello, got %v", haps[0].Value)
	}
	if !haps[0].Whole.Equals(NewTimeSpan(FractionFromInt(0), FractionFromInt(1))) {
		t.Fatalf("whole mismatch %v", haps[0].Whole)
	}
}

func TestPureQueryArc(t *testing.T) {
	p := Pure("x")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps) != 2 {
		t.Fatalf("expected 2 haps for 0-2, got %d", len(haps))
	}
}

func TestStack(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	haps := p.FirstCycle()
	if len(haps) != 2 {
		t.Fatalf("stack expected 2, got %d: %v", len(haps), haps)
	}
}

func TestSlowCat(t *testing.T) {
	a := Pure("a")
	b := Pure("b")
	p := SlowCat(a, b)
	// cycle 0 should be a, cycle 1 should be b
	haps0 := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps0) != 1 || haps0[0].Value != "a" {
		t.Fatalf("slowcat cycle0 expected a, got %v", haps0)
	}
	haps1 := p.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(haps1) != 1 || haps1[0].Value != "b" {
		t.Fatalf("slowcat cycle1 expected b, got %v", haps1)
	}
}

func TestWithValue(t *testing.T) {
	p := Pure(3).Fmap(func(v any) any { return v.(int) + 4 })
	haps := p.FirstCycle()
	if haps[0].Value != 7 {
		t.Fatalf("fmap expected 7, got %v", haps[0].Value)
	}
}

func TestFastF(t *testing.T) {
	p := Pure("x").FastF(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Pure fast 2: each cycle now contains 2 haps (compressed whole cycles)
	if len(haps) != 2 {
		t.Fatalf("fast pure expected 2 haps, got %d: %v", len(haps), haps)
	}
	// Check they cover 0-0.5 and 0.5-1
	if !haps[0].Part.Equals(NewTimeSpan(FractionFromInt(0), MustParseFraction("1/2"))) {
		t.Fatalf("first hap part wrong: %v", haps[0].Part)
	}
}
