package core

import "testing"

func TestMJS_Port193_FastSlowWithStructureFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 6 {
		t.Fatalf("FastCat3 FastF2 6 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	q := FastCat(Pure("a"), Pure("b")).Slow(FractionFromInt(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Slow 2")
	}
	r := Pure("z").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0.25-0.75")
	}
}

func TestMJS_Port193_InsideOutsideWithSignalFourth(t *testing.T) {
	p := Sine().Inside(FractionFromInt(2), func(q Pattern) Pattern { return q.Rev() })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Inside 2 Rev")
	}
	q := Saw().Outside(FractionFromInt(2), func(x Pattern) Pattern { return x.FastF(FractionFromInt(3)) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Outside 2 FastF3")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Inside(FractionFromInt(4), func(q Pattern) Pattern { return q.Rev() })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Inside 4 Rev <2")
	}
}

func TestMJS_Port193_SegmentChopWithFastFourth(t *testing.T) {
	p := Pure("a b c d").Segment(4)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Segment 4")
	}
	q := Pure("a b c d").Chop(2)
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chop 2")
	}
	r := Pure("a b").Striate(2)
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Striate 2")
	}
	s := Slice(2, 0, 1)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice 2 0 1")
	}
}
