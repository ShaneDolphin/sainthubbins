package core

import "testing"

func TestMJS_ArpeggioChainSlow2(t *testing.T) {
	p := Pure("c3").Slow(FractionFromInt(2))
	a := p.Arp("up")
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Arp slow") }
}

func TestMJS_DegradeBySignalRev2(t *testing.T) {
	p := Stack(Pure("bd"), Pure("sd"))
	d := p.DegradeBy(0.3).Rev()
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("DegradeBy Rev") }
}

func TestMJS_PatternControlStack2(t *testing.T) {
	p := Stack(Pure(map[string]any{"s": "bd"}), Pure(map[string]any{"s": "sd"}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)!=2 { t.Fatalf("Control Stack 2 got %d", len(haps)) }
}
