package core

import "testing"

func TestMJS_Port114_CompressZoomWithinSegment(t *testing.T) {
	p := Pure("bd sd").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0.25-0.75")
	}
	z := Pure("a b c").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0-0.5")
	}
	w := Pure("bd sd").Within(0.25, 0.75, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Within 0.25-0.75")
	}
	seg := Pure("a b c d").Segment(2)
	if seg.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Segment 2")
	}
}

func TestMJS_Port114_EveryWhenOffSometimes(t *testing.T) {
	p := Pure("bd").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2")
	}
	whenTrue := Pure("bd").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(whenTrue.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true should fast 2")
	}
	whenFalse := Pure("bd").When(false, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(whenFalse.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false should be 1")
	}
	off := Pure("bd").Off(0.25, func(q Pattern) Pattern { return q.Add(Pure(10)) })
	if len(off.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 should duplicate")
	}
	som := Pure("bd").Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if som.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes nil")
	}
}

func TestMJS_Port114_FastSlowEarlyLate(t *testing.T) {
	f := Pure("bd").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Fast 2")
	}
	s := Pure("bd sd").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2 nil")
	}
	e := Pure("bd").Early(FractionFromFloat(0.25))
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Early 0.25")
	}
	l := Pure("bd").Late(FractionFromFloat(0.25))
	if l.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Late 0.25")
	}
}
