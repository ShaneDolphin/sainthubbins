package core

import "testing"

func TestMJS_Port131_BinaryOpAddMul(t *testing.T) {
	p := FastCat(Pure(1), Pure(2))
	q := p.Add(Pure(10))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 10")
	}
	m := FastCat(Pure(2), Pure(3)).Mul(Pure(4))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mul 4")
	}
	// Float tolerant check like earlier ports
	add := Pure(1.5).Add(Pure(2.5))
	if add.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) < 3.9 {
		t.Fatalf("Add float 1.5+2.5")
	}
}

func TestMJS_Port131_StructureWithSignal(t *testing.T) {
	p := Pure("a").Struct(FastCat(Pure(true), Pure(false), Pure(true)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf t")
	}
	q := Pure("bd").Struct(Pure(true))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true")
	}
	s := Sine().Struct(Pure(true))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Struct true")
	}
}

func TestMJS_Port131_EveryWhenOff(t *testing.T) {
	e := Pure("bd").Every(2, func(q Pattern) Pattern { return q.Rev() })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2 Rev")
	}
	w := Pure("bd sd").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true")
	}
	o := Pure("a b c").Off(0.125, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125")
	}
}
