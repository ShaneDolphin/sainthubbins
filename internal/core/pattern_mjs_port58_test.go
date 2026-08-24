package core

import "testing"

func TestMJS_ArpeggioChain2(t *testing.T) {
	p := Pure("a")
	// Slow 2 then FastF 2 should cancel
	sf := p.SlowF(FractionFromInt(2)).FastF(FractionFromInt(2))
	haps := sf.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Slow2 FastF2 cancel expected non-empty")
	}
}

func TestMJS_DegradeByChain2(t *testing.T) {
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DegradeBy 0 expected 1")
	}
	d1 := Pure("a").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("DegradeBy 1 expected 0")
	}
}

func TestMJS_PatternValues2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "gain": 0.5})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Pure map s bd gain 0.5 expected 1")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["s"] != "bd" || m["gain"] != 0.5 {
		t.Fatalf("map s bd gain 0.5 got %v", haps[0].Value)
	}
}
