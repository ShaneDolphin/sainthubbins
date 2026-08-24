package core

import "testing"

func TestMJS_Port196_BinaryOpAddMulSubFourth(t *testing.T) {
	p := Pure(10).Add(Pure(5))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 15 {
		t.Fatalf("Add 10+5=15")
	}
	q := Pure(6).Mul(Pure(7))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 42 {
		t.Fatalf("Mul 6*7=42")
	}
	r := Pure(20).Sub(Pure(8))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 12 {
		t.Fatalf("Sub 20-8=12")
	}
	s := Pure(9).Div(Pure(3))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 3 {
		t.Fatalf("Div 9/3=3")
	}
}

func TestMJS_Port196_StructureWithSignalSineFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Struct(Pure(true))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Struct true 3")
	}
	q := Sine().Range(0, 1)
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,1")
	}
	r := Saw().Range(-5, 5)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range -5,5")
	}
}

func TestMJS_Port196_EveryWhenOffSlowFastFourth(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Every 2")
	}
	q := Pure("x").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	r := Pure("bd").Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 <2")
	}
	s := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Fast Slow 2")
	}
}
