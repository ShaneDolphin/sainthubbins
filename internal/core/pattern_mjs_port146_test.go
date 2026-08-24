package core

import "testing"

func TestMJS_Port146_BinaryOpAddMulSub(t *testing.T) {
	p := FastCat(Pure(10), Pure(20))
	q := p.Add(Pure(5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 5")
	}
	m := Pure(10).Mul(Pure(2))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 20 {
		t.Fatalf("Mul 2 ->20")
	}
	s := Pure(20).Sub(Pure(7))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 13 {
		t.Fatalf("Sub 7 ->13")
	}
	d := Pure(15).Div(Pure(3))
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 5 {
		t.Fatalf("Div 3 ->5")
	}
}

func TestMJS_Port146_StructureWithSignalSine(t *testing.T) {
	p := Pure("bd").Struct(FastCat(Pure(true), Pure(false)))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct tf")
	}
	s := Sine().Slow(FractionFromInt(2)).Struct(Pure(true))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Struct true")
	}
	r := Saw().Range(0, 10).Struct(Pure(true))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range Struct")
	}
}

func TestMJS_Port146_EveryWhenOffSlowFast(t *testing.T) {
	e := Pure("bd sd").Every(2, func(q Pattern) Pattern { return q.Rev() })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2 Rev")
	}
	w := Pure("a").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true Fast 2")
	}
	o := Pure("a b c d").Off(0.5, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.5")
	}
	f := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Fast Slow 2 2 =>2")
	}
}
