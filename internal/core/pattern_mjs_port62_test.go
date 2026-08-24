package core

import "testing"

func TestMJS_FastSlowEdge2(t *testing.T) {
	f := Pure("a").FastF(FractionFromInt(2))
	haps := f.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("FastF 2 expected 2 got %d", len(haps))
	}
	s := Pure("a").SlowF(FractionFromInt(2))
	haps2 := s.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps2) == 0 {
		t.Fatalf("SlowF 2 expected non-empty over 2 cycles")
	}
}

func TestMJS_BinaryOp2(t *testing.T) {
	a := Pure(3).Add(Pure(4))
	haps := a.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || toFloat(haps[0].Value) != 7 {
		t.Fatalf("Add 3+4=7 got %v", haps[0].Value)
	}
	sub := Pure(7).Sub(Pure(4))
	haps2 := sub.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps2[0].Value) != 3 {
		t.Fatalf("Sub 7-4=3 got %v", haps2[0].Value)
	}
	mul := Pure(3).Mul(Pure(4))
	if toFloat(mul.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value) != 12 {
		t.Fatalf("Mul 3*4=12")
	}
}

func TestMJS_CompressWithin2(t *testing.T) {
	c := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := c.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Compress 0.25-0.75 expected non-empty")
	}
	w := Pure("a").Within(0.5, 1.0, func(p Pattern) Pattern {
		return p.FastF(FractionFromInt(2))
	})
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Within 0.5-1 FastF2 expected non-empty")
	}
}
