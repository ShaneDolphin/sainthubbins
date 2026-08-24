package core

import "testing"

func TestMJS_SlowInsideOutside3(t *testing.T) {
	p := Pure("a").SlowF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("SlowF 2")
	}
	in := Pure("a").Inside(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(in.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside 2")
	}
	out := Pure("a").Outside(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(out.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside 2")
	}
}

func TestMJS_EuclidBjorklundSignal3(t *testing.T) {
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund 3,8 len 8")
	}
	e := Pure("a").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Euclid 3,8 3")
	}
	s := Sine().Range(0, 1)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Sine Range")
	}
}

func TestMJS_ValueSteps4(t *testing.T) {
	p := Pure("a")
	if p.Steps == nil || !p.Steps.Equals(FractionFromInt(1)) {
		t.Fatalf("Pure steps 1")
	}
	ws := p.WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	if ws.Steps == nil || !ws.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("WithSteps 2")
	}
}
