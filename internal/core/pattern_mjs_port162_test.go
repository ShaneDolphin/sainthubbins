package core

import "testing"

func TestMJS_Port162_BinaryOpAddMulSecondThird(t *testing.T) {
	p := FastCat(Pure(2), Pure(3))
	q := p.Add(Pure(5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 5 len 2")
	}
	m := FastCat(Pure(4), Pure(5)).Mul(Pure(2))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mul 2")
	}
	s := Pure(10).Sub(Pure(4))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 6 {
		t.Fatalf("Sub 4 ->6")
	}
}

func TestMJS_Port162_StructureWithSignalThird(t *testing.T) {
	p := Pure("bd").Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(true)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf tt")
	}
	q := Saw().Range(0, 5).Struct(Pure(true))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Struct true")
	}
	r := Tri().Range(-1, 1).Mask(FastCat(Pure(true), Pure(false)))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Mask tf")
	}
}

func TestMJS_Port162_EveryWhenOffThird(t *testing.T) {
	e := Pure("bd").Every(2, func(q Pattern) Pattern { return q.Rev() })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2 Rev")
	}
	w := Pure("a").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(3)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("When true Fast 3")
	}
	o := Pure("bd sd hh").Off(0.5, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.5 Add 1")
	}
}
