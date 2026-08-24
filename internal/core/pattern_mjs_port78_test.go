package core

import "testing"

func TestMJS_PlyOffWhen3(t *testing.T) {
	p := Pure("a").Ply(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Ply 2")
	}
	o := Pure("a").Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25")
	}
	w := Pure("a").When(Pure(true), func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true FastF2")
	}
}

func TestMJS_BrakHurry3(t *testing.T) {
	b := Pure("a").Brak()
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
	h := Pure("a").Hurry(2)
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 2")
	}
}

func TestMJS_SegmentChop3(t *testing.T) {
	s := Pure("a").Segment(2)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Segment 2")
	}
	c := Pure("a").Chop(2)
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chop 2")
	}
}
