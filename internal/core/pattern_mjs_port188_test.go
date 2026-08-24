package core

import "testing"

func TestMJS_Port188_DegradeByChainSignalFourth(t *testing.T) {
	d0 := Pure("a").DegradeBy(0.2)
	if d0.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.2 nil")
	}
	d1 := Pure("b").Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if d1.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes Rev nil")
	}
	d2 := Pure("c").Degrade()
	if d2.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Degrade nil")
	}
	s := Sine().DegradeBy(0.1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine DegradeBy 0.1")
	}
}

func TestMJS_Port188_SometimesByOftenRarelyFourth(t *testing.T) {
	p := Pure("a").SometimesBy(0.5, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.5")
	}
	q := Pure("b").Sometimes(func(pat Pattern) Pattern { return pat.Rev() })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes Rev")
	}
	r := Pure("c").SometimesBy(0.75, func(q Pattern) Pattern { return q.Rev() })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.75")
	}
}

func TestMJS_Port188_ArpWithSlowFastFourth(t *testing.T) {
	a := Pure("c3 e3 g3").Arp("up").Slow(FractionFromInt(2))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Arp up Slow2")
	}
	b := Pure("c3 e3 g3").Arp("down").FastF(FractionFromInt(2))
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down FastF2")
	}
	c := Pure("c3 e3 g3").Arp("converge")
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge")
	}
}
