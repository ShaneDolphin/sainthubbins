package core

import "testing"

func TestMJS_Port123_DegradeByChainSignal(t *testing.T) {
	p := Pure("bd").DegradeBy(0.2).Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.2 Sometimes")
	}
	q := Pure("a b c").Degrade()
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Degrade")
	}
	r := Sine().DegradeBy(0.5)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine DegradeBy 0.5")
	}
}

func TestMJS_Port123_SometimesByOftenRarely(t *testing.T) {
	s := Pure("bd").SometimesBy(0.5, func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.5 Rev")
	}
	o := Pure("bd").Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes")
	}
	// Often is SometimesBy 0.75 in JS, but we test via SometimesBy directly
	often := Pure("bd").SometimesBy(0.75, func(q Pattern) Pattern { return q.Rev() })
	if often.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.75")
	}
}

func TestMJS_Port123_ArpWithSlowFast(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up").Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp up Slow 2")
	}
	q := Pure("c3 e3 g3").Arp("down").FastF(FractionFromInt(2))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp down Fast 2")
	}
	r := Pure("c3 e3 g3 b3").Arp("converge")
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge")
	}
}
