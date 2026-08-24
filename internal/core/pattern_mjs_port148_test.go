package core

import "testing"

func TestMJS_Port148_FastSlowWithStructure(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 6 {
		t.Fatalf("FastCat3 Fast2 =>6 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	s := Pure("bd sd").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	c := Pure("a b c").Compress(FractionFromFloat(0), FractionFromFloat(1)).Slow(FractionFromInt(1))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0-1 Slow 1")
	}
}

func TestMJS_Port148_InsideOutsideWithSignal(t *testing.T) {
	p := Sine().Inside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Inside 2 Fast 2")
	}
	q := Saw().Outside(3, func(pat Pattern) Pattern { return pat.Range(0, 5) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Outside 3 Range")
	}
	r := Pure("a b c d").Inside(4, func(pat Pattern) Pattern { return pat.Rev() })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Inside 4 Rev")
	}
}

func TestMJS_Port148_SegmentChopWithFast(t *testing.T) {
	s := Pure("a b c d").Segment(4)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Segment 4")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chop(2)
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chop 2")
	}
	fast := Pure("bd").Striate(2)
	if fast.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Striate 2")
	}
	sl := Slice(4, Pure(0), S("bd"))
	if sl.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slice 4 0 bd")
	}
}
