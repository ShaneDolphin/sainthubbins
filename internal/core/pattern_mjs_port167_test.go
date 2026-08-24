package core

import "testing"

func TestMJS_Port167_StackWithRestControlFourth(t *testing.T) {
	s := Stack(Pure("bd"), Silence(), S("sd"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack bd Silence sd 2")
	}
	p := Stack(Pure("hh"), Gain(0.5))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack hh Gain 2")
	}
	c := Cat(Pure("a"), Pure("b"), Silence())
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 2 {
		t.Fatalf("Cat a b Silence 3 cycles 2")
	}
}

func TestMJS_Port167_PatternSlowFastWhenFourth(t *testing.T) {
	s := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Slow 2 =>2")
	}
	f := Pure("a").FastF(FractionFromInt(3))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastF 3")
	}
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
}

func TestMJS_Port167_ArpWithMasksSignalFourth(t *testing.T) {
	arp := Pure("c3 e3 g3 b3").Arp("updown")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
	m := arp.Mask(FastCat(Pure(true), Pure(false), Pure(true)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Mask tf t")
	}
	sig := Saw().Range(0, 10).Slow(FractionFromInt(2))
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,10 Slow 2")
	}
}
