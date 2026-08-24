package core

import "testing"

func TestMJS_Port257_CompressZoomWithinSegmentFourth(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Whole == nil {
		t.Fatalf("Compress whole nil")
	}
	q := Pure("b").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Zoom 0-0.5")
	}
	r := Pure("c d e f").Within(0.5, 1, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Within 0.5-1 FastF2")
	}
	s := Pure("x y z").Segment(2)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Segment 2")
	}
}

func TestMJS_Port257_EveryWhenOffSometimesFourth(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Every 2")
	}
	q := Pure("bd").Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 >=2")
	}
	r := Pure("a").When(true, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("When true")
	}
	s := Pure("a").Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes Rev")
	}
}

func TestMJS_Port257_FastSlowEarlyLateFourth(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastF2 4 got %d", len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	s := FastCat(Pure("a"), Pure("b")).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Slow 2")
	}
	e := Pure("a").Early(FractionFromFloat(0.25))
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Early 0.25")
	}
	l := Pure("a").Late(FractionFromFloat(0.25))
	if l.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Late 0.25")
	}
}
