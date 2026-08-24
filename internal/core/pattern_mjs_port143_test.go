package core

import "testing"

func TestMJS_Port143_FastSlowCompressZoom(t *testing.T) {
	f := Pure("a b c d").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastF 2")
	}
	s := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Slow 2 =>2")
	}
	c := Pure("bd sd hh oh").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0.25-0.75")
	}
	z := Pure("a b c").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0-0.5")
	}
}

func TestMJS_Port143_InsideOutsideRevBrak(t *testing.T) {
	p := Pure("bd sd").Inside(2, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside 2 Rev")
	}
	q := Pure("bd sd").Outside(2, func(q2 Pattern) Pattern { return q2.FastF(FractionFromInt(2)) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Outside 2 Fast 2")
	}
	rev := Pure("a b c").Rev()
	if len(rev.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rev")
	}
	if len(Pure("a b c").Brak().QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
}

func TestMJS_Port143_SegmentChopStriate(t *testing.T) {
	s := Pure("a b c d").Segment(2)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Segment 2")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chop(2)
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chop 2")
	}
	st := Pure("a b c d").Striate(3)
	if st.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Striate 3")
	}
	sl := Slice(2, Pure("a"), Pure("a b c"))
	if sl.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slice 2")
	}
}
