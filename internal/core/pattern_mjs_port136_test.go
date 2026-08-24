package core

import "testing"

func TestMJS_Port136_StackWithRestControl(t *testing.T) {
	s := Stack(Pure("bd"), Silence())
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Stack bd Silence 1")
	}
	c := Cat(Pure("bd"), Silence())
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(2))) != 1 {
		t.Fatalf("Cat bd Silence 2 cycles 1")
	}
	p := Stack(S("bd"), Gain(0.5))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack S bd Gain 0.5 2")
	}
}

func TestMJS_Port136_PatternSlowFastWhen(t *testing.T) {
	s := Pure("a b c").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(3))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastF 3")
	}
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	w2 := Pure("bd").When(false, func(q Pattern) Pattern { return q.Rev() })
	if len(w2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
}

func TestMJS_Port136_ArpWithMasksSignal(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	masked := arp.Mask(FastCat(Pure(true), Pure(false), Pure(true)))
	if masked.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Mask")
	}
	sig := Saw().Range(0, 5)
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,5")
	}
}
