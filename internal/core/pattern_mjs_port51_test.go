package core

import "testing"

func TestMJS_EveryWhenOff(t *testing.T) {
	e := Pure("a").Every(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Every 2 expected non-empty")
	}
	w := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true FastF 2 expected 2 got %d", len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	o := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 expected >=2")
	}
}

func TestMJS_InsideOutside2(t *testing.T) {
	in := Pure("a").Inside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(in.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside 2 expected non-empty")
	}
	out := Pure("a").Outside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(out.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside 2 expected non-empty")
	}
}

func TestMJS_CompressZoom2(t *testing.T) {
	c := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := c.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Compress 0.25-0.75 expected non-empty")
	}
	z := Pure("a").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	haps2 := z.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Zoom 0-0.5 expected non-empty")
	}
}
