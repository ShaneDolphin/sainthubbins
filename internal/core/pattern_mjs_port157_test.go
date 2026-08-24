package core

import "testing"

func TestMJS_Port157_StackWithRestControlThird(t *testing.T) {
	s := Stack(S("bd"), Silence(), S("sd"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack bd Silence sd 2")
	}
	p := Stack(Pure("a"), Pure("b")).WithValue(func(v any) any { return v })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue 2")
	}
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
}

func TestMJS_Port157_PatternSlowFastWhenThird(t *testing.T) {
	s := Pure("a b c d").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2")
	}
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	e := Pure("bd").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2")
	}
}

func TestMJS_Port157_ArpWithMasksSignalThird(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	m := arp.Mask(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Mask tf tf")
	}
	sig := Sine().Range(0, 5).FastF(FractionFromInt(2))
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,5 Fast 2")
	}
}
