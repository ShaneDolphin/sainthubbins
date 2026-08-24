package core

import "testing"

func TestMJS_SteadySignal2(t *testing.T) {
	p := Pure(0.5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Pure 0.5")
	}
	s := Sine().Range(0, 1)
	haps2 := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Sine Range")
	}
}

func TestMJS_ExpandArp2(t *testing.T) {
	p := Pure(2).Fmap(func(v any) any { return v.(int) * 2 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(int) != 4 {
		t.Fatalf("Fmap 2*2=4")
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"))
	arp := seq.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
}

func TestMJS_IdPly2(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Fmap id")
	}
	ply := Pure("a").Ply(3)
	haps := ply.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Ply 3")
	}
}
