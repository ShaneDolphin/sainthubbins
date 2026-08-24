package core

import "testing"

func TestMJS_ArpeggioChain3(t *testing.T) {
	p := Pure("a")
	sf := p.SlowF(FractionFromInt(2)).FastF(FractionFromInt(2))
	if len(sf.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slow2 FastF2 expected non-empty")
	}
}

func TestMJS_DegradeByChain3(t *testing.T) {
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DegradeBy 0 expected 1")
	}
	d1 := Pure("a").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("DegradeBy 1 expected 0")
	}
}

func TestMJS_PatternValues3(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "gain": 0.5})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Pure map s bd gain 0.5 expected 1")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["s"] != "bd" {
		t.Fatalf("map s bd got %v", haps[0].Value)
	}
}
