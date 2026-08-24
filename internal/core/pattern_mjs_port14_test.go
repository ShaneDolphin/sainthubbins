// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Fourteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_TimeSpanHap(t *testing.T) {
	ts := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if len(ts.SpanCycles()) != 2 {
		t.Fatalf("spanCycles 0-2 expected 2 got %d", len(ts.SpanCycles()))
	}
	wholeA := NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5))
	wholeC := NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.25))
	a := NewHap(&wholeA, NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)), "a", nil)
	b := NewHap(&wholeA, NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)), "a", nil)
	if !a.SpanEquals(b) {
		t.Fatalf("spanEquals true expected")
	}
	c := NewHap(&wholeC, NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)), "a", nil)
	if a.SpanEquals(c) {
		t.Fatalf("spanEquals false expected (whole differs)")
	}
	d := NewHap(nil, NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)), "d", nil)
	e := NewHap(nil, NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)), "e", nil)
	if !d.SpanEquals(e) {
		t.Fatalf("spanEquals nil wholes true expected")
	}
}

func TestMJS_PatternPureFmap(t *testing.T) {
	p := Pure("hello")
	if n := len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 1 {
		t.Fatalf("pure hello 0-1 expected 1 got %d", n)
	}
	p2 := Pure(3).Fmap(func(x any) any { return x.(int) + 4 })
	if v := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(int); v != 7 {
		t.Fatalf("fmap 3+4 expected 7 got %v", v)
	}
}

func TestMJS_AddWithStructure(t *testing.T) {
	p := FastCat(Pure(1), Pure(2)).Add(Pure(10))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("FastCat 1,2 +10 expected 2 got %d", len(haps))
	}
	// check values 11,12 (via float)
	for i, h := range haps {
		var f float64
		switch x := h.Value.(type) {
		case int:
			f = float64(x)
		case float64:
			f = x
		default:
			f = 0
		}
		if f != 11 && f != 12 {
			t.Fatalf("hap %d value %v expected 11 or 12", i, h.Value)
		}
	}
}
