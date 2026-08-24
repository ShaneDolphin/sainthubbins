package core

import "testing"

func TestMJS_Port172_BinaryOpAddMulThirdFourth(t *testing.T) {
	p := FastCat(Pure(1), Pure(2), Pure(3))
	q := p.Add(Pure(10))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Add 10 len 3")
	}
	m := FastCat(Pure(2), Pure(3)).Mul(Pure(4))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mul 4 len 2")
	}
	s := Pure(10).Sub(Pure(3))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 7 {
		t.Fatalf("Sub 3 ->7")
	}
	d := Pure(12).Div(Pure(3))
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 4 {
		t.Fatalf("Div 3 ->4")
	}
}

func TestMJS_Port172_StructureWithSignalFourth(t *testing.T) {
	p := Pure("a").Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf tf")
	}
	q := Sine().Range(0, 1).Mask(FastCat(Pure(true), Pure(true)))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Mask tt")
	}
	r := Saw().Range(0, 10).Struct(Pure(true))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Struct true")
	}
}

func TestMJS_Port172_EveryWhenOffFourth(t *testing.T) {
	e := Pure("bd").Every(3, func(q Pattern) Pattern { return q.Rev() })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 3 Rev")
	}
	w := Pure("a").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true Fast 2")
	}
	o := Pure("bd sd").Off(0.125, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125 Fast 2")
	}
}
