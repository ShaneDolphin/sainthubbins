package core

import "testing"

func TestMJS_Port142_BinaryOpAddMulSubDivMod(t *testing.T) {
	p := FastCat(Pure(5), Pure(10))
	a := p.Add(Pure(2))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 2")
	}
	m := FastCat(Pure(2), Pure(4)).Mul(Pure(3))
	hm := m.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if hm[0].Value.(float64) < 5.9 {
		t.Fatalf("Mul 3 got %v", hm[0].Value)
	}
	s := Pure(20).Sub(Pure(5))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 15 {
		t.Fatalf("Sub 5 ->15")
	}
	d := Pure(20).Div(Pure(4))
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 5 {
		t.Fatalf("Div 4 ->5")
	}
	mo := Pure(10).Mod(Pure(3))
	if mo.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mod 3")
	}
}

func TestMJS_Port142_StructureWithSignalRange(t *testing.T) {
	s := Pure("a").Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct t f t f")
	}
	sig := Sine().Range(0, 1).Struct(Pure(true))
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Struct true")
	}
	m := Pure("bd").Mask(FastCat(Pure(true), Pure(true), Pure(false)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tt f")
	}
}

func TestMJS_Port142_EveryWhenOffSometimesDegrade(t *testing.T) {
	e := Pure("bd").Every(3, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 3")
	}
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	o := Pure("a b c").Off(0.25, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25")
	}
	s := Pure("bd").SometimesBy(0.5, func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.5")
	}
}
