// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 31st batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FilterEvents(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c")).Filter(func(h Hap) bool {
		if s, ok := h.Value.(string); ok {
			return s != "b"
		}
		return true
	})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Should have filtered b, leaving a and c onset pieces.
	for _, h := range haps {
		if h.Value == "b" {
			t.Fatalf("filter should remove b got %v", h)
		}
	}
	if len(haps) == 0 {
		t.Fatalf("filter left empty")
	}
}

func TestMJS_OnsetsOnly(t *testing.T) {
	// Logic query-like: only haps with onset in arc
	p := Pure("a").Slow(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	onsets := 0
	for _, h := range haps {
		if h.HasOnset() {
			onsets++
		}
	}
	if onsets == 0 {
		t.Fatalf("onsets expected >0 got %d haps=%v", onsets, haps)
	}
}

func TestMJS_SpanConversions(t *testing.T) {
	s := TimeSpan{Begin: FractionFromInt(0), End: FractionFromInt(1)}
	if s.Duration().Cmp(FractionFromInt(1)) != 0 {
		t.Fatalf("duration 1 got %v", s.Duration())
	}
	s2 := TimeSpan{Begin: FractionFromFloat(0.25), End: FractionFromFloat(0.75)}
	if s2.Duration().Cmp(FractionFromFloat(0.5)) != 0 {
		t.Fatalf("duration 0.5 got %v", s2.Duration())
	}
	// Span equals
	a := TimeSpan{Begin: FractionFromInt(0), End: FractionFromInt(1)}
	b := TimeSpan{Begin: FractionFromInt(0), End: FractionFromInt(1)}
	if !a.Equals(b) {
		t.Fatalf("span equals true")
	}
	c := TimeSpan{Begin: FractionFromInt(0), End: FractionFromFloat(0.5)}
	if a.Equals(c) {
		t.Fatalf("span equals false")
	}
}
