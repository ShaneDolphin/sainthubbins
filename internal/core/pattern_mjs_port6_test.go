// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Sixth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SignalRange(t *testing.T) {
	// signals: saw, sine etc are func()Pattern that vary over time
	s := Saw().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(s) == 0 {
		t.Fatalf("saw empty")
	}
	// range
	r := Saw().Range(0, 100)
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("saw range empty")
	}
	// sine
	si := Sine().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(si) == 0 {
		t.Fatalf("sine empty")
	}
}

func TestMJS_HapState(t *testing.T) {
	// TimeSpan
	ts := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	if ts.Duration().Cmp(FractionFromInt(1)) != 0 {
		t.Fatalf("timespan duration")
	}
	// Hap hasOnset
	whole := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	part := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&whole, part, "a", nil)
	if !h.HasOnset() {
		t.Fatalf("hap hasOnset")
	}
	// State
	st := NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil)
	if st.Span.Duration().Cmp(FractionFromInt(1)) != 0 {
		t.Fatalf("state span")
	}
}

func TestMJS_LogAndFmap(t *testing.T) {
	// logValues via Fmap
	p := Pure(2).Fmap(func(v any) any { return v.(int) * 3 })
	if v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(int); v != 6 {
		t.Fatalf("fmap 2*3 got %v", v)
	}
	// stack with log-like pattern
	s := Stack(Pure("a").Fmap(func(v any) any { return v }), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("stack fmap 2 expected 2")
	}
}
