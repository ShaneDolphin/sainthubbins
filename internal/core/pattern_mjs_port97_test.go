package core

import "testing"

func TestMJS_ArpeggioChordStack3(t *testing.T) {
	p := Pure("c3")
	arp := p.Arp("up")
	stacked := Stack(arp, arp.Slow(FractionFromInt(2)))
	haps := stacked.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Arpeggio Stack expected haps")
	}
}

func TestMJS_DegradeByChainSignal3(t *testing.T) {
	p := Stack(Pure("bd"), Pure("sd"))
	d := p.DegradeBy(0.5)
	slo := d.Slow(FractionFromInt(2))
	j := slo.Jux(func(q Pattern) Pattern { return q.Add(Pure(1)) })
	haps := j.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("DegradeBy Slow Jux expected haps")
	}
}

func TestMJS_SometimesByOften2(t *testing.T) {
	p := Pure("bd")
	s := p.SometimesBy(0.5, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	o := s.Slow(FractionFromInt(2))
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy Often expected non-nil")
	}
}
