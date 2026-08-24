package core

import "testing"

func TestMJS_Port124_CompressWithinSlowFast(t *testing.T) {
	p := Pure("a b c").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75)).Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress Slow 2")
	}
	q := Pure("a b c d").Within(0.5, 1, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Within 0.5-1 Fast 2")
	}
	r := Pure("a").FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	// Fast then Slow should cancel to 1
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Fast 2 Slow 2 cancel 1 got %d", len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_Port124_EuclidBjorklundLegato(t *testing.T) {
	e := Pure("bd").EuclidLegato(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidLegato 3,8")
	}
	b := Bjorklund(5, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund 5,8 len 8 got %d", len(b))
	}
	er := Pure("bd").EuclidRot(3, 8, 1)
	if len(er.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidRot 3,8,1")
	}
}

func TestMJS_Port124_ScaleQuantizeWithOff(t *testing.T) {
	p := Pure("c3 d3 e3").Scale("major")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale major")
	}
	q := Pure("c3").Scale("minor").Add(Pure(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor Add 2")
	}
	r := Pure("bd").Euclid(3, 8).Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Euclid Off 0.25")
	}
}
