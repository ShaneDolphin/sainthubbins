package core

import "testing"

func TestMJS_ArpeggioChain4(t *testing.T) {
	p := Pure("a").SlowF(FractionFromInt(2)).FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slow2 FastF2 expected non-empty")
	}
}

func TestMJS_DegradeByChain4(t *testing.T) {
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DegradeBy 0 expected 1")
	}
	d1 := Pure("a").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("DegradeBy 1 expected 0")
	}
}

func TestMJS_PatternValues4(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(map[string]any)["s"] != "bd" {
		t.Fatalf("Pure map s bd")
	}
}
