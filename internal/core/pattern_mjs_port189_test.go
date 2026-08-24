package core

import "testing"

func TestMJS_Port189_CompressWithinSlowFastFourth(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Whole == nil {
		t.Fatalf("Compress whole nil")
	}
	q := Pure("b c").Within(0.5, 1, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Within 0.5-1 FastF2")
	}
	r := Pure("x y").Slow(FractionFromInt(2))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Slow 2")
	}
	s := Pure("z").FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("FastF2 Slow2 cancel 1")
	}
}

func TestMJS_Port189_EuclidBjorklundLegatoFourth(t *testing.T) {
	e := Pure("bd").EuclidLegato(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidLegato 3,8")
	}
	b := Bjorklund(5, 8)
	cnt := 0
	for _, v := range b {
		if v != 0 {
			cnt++
		}
	}
	if cnt != 5 {
		t.Fatalf("Bjorklund 5 !=5 %d", cnt)
	}
	r := Pure("bd").EuclidRot(3, 8, 1)
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidRot 3,8,1")
	}
}

func TestMJS_Port189_ScaleQuantizeWithOffFourth(t *testing.T) {
	s := Pure("c3").Scale("major")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale major")
	}
	t2 := Pure("c3 d3 e3").Scale("minor")
	if len(t2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	o := Pure("bd").Euclid(3, 8).Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Euclid Off 0.25 <2")
	}
}
