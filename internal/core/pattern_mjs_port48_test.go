package core

import "testing"

func TestMJS_SegmentSlice(t *testing.T) {
	// Segment 2 should split into 2 per cycle
	p := Pure("a").Segment(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Segment 2 expected non-empty")
	}
	// Slice free function
	sl := Slice(8, Pure(0), Pure("a"))
	if len(sl.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice expected non-empty")
	}
}

func TestMJS_FitAndLoopAt(t *testing.T) {
	// Fit should not be empty
	f := Pure("a").Fit()
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Fit expected non-empty")
	}
	la := Pure("a").LoopAt(2)
	if len(la.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("LoopAt 2 expected non-empty")
	}
	// loopAtCps lower alias
	lac := Pure("a").Loopatcps(2, 1.0)
	if len(lac.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Loopatcps expected non-empty")
	}
}

func TestMJS_RollAndStut(t *testing.T) {
	// Stut already but test stut alias
	s := Pure("a").Stut(2, 0.5, 0.25)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Stut expected non-empty")
	}
	// EchoWith vs StutWith alias equivalence (both stack)
	ew := Pure("a").EchoWith(2, FractionFromFloat(0.25), func(p Pattern, n int) Pattern { return p })
	sw := Pure("a").StutWith(2, FractionFromFloat(0.25), func(p Pattern, n int) Pattern { return p })
	if len(ew.QueryArc(FractionFromInt(0), FractionFromInt(1))) != len(sw.QueryArc(FractionFromInt(0), FractionFromInt(1))) {
		t.Fatalf("EchoWith vs StutWith length mismatch %d vs %d", len(ew.QueryArc(FractionFromInt(0), FractionFromInt(1))), len(sw.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	// ApplyN 0 should be identity
	an0 := Pure(5).ApplyN(0, func(p Pattern) Pattern { return p.Fmap(func(v any) any { return v.(int) + 1 }) })
	haps := an0.FirstCycle()
	if len(haps) == 0 || haps[0].Value.(int) != 5 {
		t.Fatalf("ApplyN 0 expected 5 got %v", haps)
	}
}
