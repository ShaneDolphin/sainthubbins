package core

import "testing"

func TestMJS_Port416_CompressZoomWithinSegmentFourth(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Whole == nil { t.Fatalf("Compress") }
	q := Pure("b").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Zoom") }
	r := Pure("c d e f").Within(0.5, 1, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Within") }
	s := Pure("x y z").Segment(2)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Segment") }
}
func TestMJS_Port416_EveryWhenOffSometimesFourth(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Every") }
	q := Pure("bd").Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 { t.Fatalf("Off") }
	r := Pure("a").When(true, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("When") }
	s := Pure("a").Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Sometimes") }
}
func TestMJS_Port416_FastSlowEarlyLateFourth(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 { t.Fatalf("FastF2 4") }
	s := FastCat(Pure("a"), Pure("b")).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 { t.Fatalf("Slow") }
	e := Pure("a").Early(FractionFromFloat(0.25))
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Early") }
	l := Pure("a").Late(FractionFromFloat(0.25))
	if l.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Late") }
}
