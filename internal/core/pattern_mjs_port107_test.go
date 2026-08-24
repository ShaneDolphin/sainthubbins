package core

import "testing"

func TestMJS_StackWithRestControl2(t *testing.T) {
	p := Stack(Pure("bd"), Silence())
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("Stack Rest") }
}

func TestMJS_PatternSlowFastWhen2(t *testing.T) {
	p := Pure("a b c").Slow(FractionFromInt(2))
	w := p.When(func(b bool) bool { return b }, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("When SlowFast") }
}

func TestMJS_ArpWithMasksSignal2(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up")
	m := p.Mask(Pure(true).FastF(FractionFromInt(2)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Arp Mask") }
}
