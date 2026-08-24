package core

import "testing"

func TestMJS_Port155_StackWithRestControlSecond(t *testing.T) {
	s := Stack(Pure("bd"), Silence())
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Stack bd Silence 1")
	}
	p := Stack(S("hh"), Gain(0.3))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack hh Gain 2")
	}
	c := Cat(Pure("a"), Silence(), Pure("b"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 2 {
		t.Fatalf("Cat a Silence b 3 cycles 2")
	}
}

func TestMJS_Port155_PatternSlowFastWhenSecond(t *testing.T) {
	s := Pure("a b c d").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(4))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastF 4")
	}
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true Fast 2")
	}
}

func TestMJS_Port155_ArpWithMasksSignalSecond(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("updown")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
	m := arp.Mask(FastCat(Pure(true), Pure(false)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Mask tf")
	}
	sig := Tri().Range(0, 10).Slow(FractionFromInt(2))
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range 0,10 Slow 2")
	}
}
