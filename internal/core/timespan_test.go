// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestTimeSpanSpanCycles(t *testing.T) {
	ts := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if len(ts.SpanCycles()) != 2 {
		t.Fatalf("expected 2 cycles, got %d", len(ts.SpanCycles()))
	}
	ts = NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	if len(ts.SpanCycles()) != 1 {
		t.Fatalf("expected 1 cycle")
	}
	// zero-width
	ts = NewTimeSpan(FractionFromInt(1), FractionFromInt(1))
	if len(ts.SpanCycles()) != 1 {
		t.Fatalf("zero-width should be 1")
	}
}

func TestTimeSpanIntersection(t *testing.T) {
	a := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	b := NewTimeSpan(FractionFromInt(1), FractionFromInt(3))
	inter := a.Intersection(b)
	if inter == nil || !inter.Equals(NewTimeSpan(FractionFromInt(1), FractionFromInt(2))) {
		t.Fatalf("intersection failed")
	}
	// no intersect
	c := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	d := NewTimeSpan(FractionFromInt(1), FractionFromInt(2))
	if c.Intersection(d) != nil {
		t.Fatalf("expected nil for touching at 1")
	}
	// zero-width edge
	e := NewTimeSpan(MustParseFraction("1/2"), MustParseFraction("1/2"))
	f := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	if e.Intersection(f) == nil {
		t.Fatalf("zero-width inside should intersect")
	}
}

func TestTimeSpanDuration(t *testing.T) {
	ts := NewTimeSpan(FractionFromInt(0), MustParseFraction("3/2"))
	if !ts.Duration().Equals(MustParseFraction("3/2")) {
		t.Fatalf("duration failed")
	}
	if !ts.Midpoint().Equals(MustParseFraction("3/4")) {
		t.Fatalf("midpoint failed")
	}
}
