package core

import "testing"

func TestMJS_Port191_BinaryOpAddMulSubFourth(t *testing.T) {
	p := Pure(5).Add(Pure(3))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 8 {
		t.Fatalf("Add 5+3=8")
	}
	q := Pure(10).Mul(Pure(2))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 20 {
		t.Fatalf("Mul 10*2=20")
	}
	r := Pure(15).Sub(Pure(5))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 10 {
		t.Fatalf("Sub 15-5=10")
	}
	s := Pure(20).Div(Pure(4))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 5 {
		t.Fatalf("Div 20/4=5")
	}
}

func TestMJS_Port191_StructureWithSignalSineFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b")).Struct(Pure(true))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Struct true 2")
	}
	q := FastCat(Pure("a"), Pure("b"), Pure("c")).Struct(Pure(false))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Struct false 0")
	}
	r := Sine().Range(0, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,1")
	}
	s := Saw().Range(-1, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range -1,1")
	}
}

func TestMJS_Port191_EveryWhenOffSlowFastFourth(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Every 2 Rev")
	}
	q := Pure("x").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true FastF2 2")
	}
	r := Pure("bd").Off(0.5, func(pat Pattern) Pattern { return pat.Rev() })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.5 <2")
	}
	s := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Fast Slow 2")
	}
}
