// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 36th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ArpeggioChain(t *testing.T) {
	p := Sequence(Pure("c3"), Pure("e3"), Pure("g3")).Slow(FractionFromInt(2)).FastF(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("arpeggio chain Slow/Fast empty")
	}
	orig := Sequence(Pure("c3"), Pure("e3"), Pure("g3")).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != len(orig) {
		t.Fatalf("Slow2+FastF2 should equal orig len %d vs %d", len(haps), len(orig))
	}
}

func TestMJS_DegradeByChain(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).DegradeBy(0.0).DegradeBy(1.0)
	// DegradeBy 0 = no degrade, 1 = all degrade (empty)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("degradeBy 1 after 0 should be empty got %d", len(haps))
	}
	p2 := Sequence(Pure("a"), Pure("b")).DegradeBy(0.0)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("degradeBy 0 expected 2 got %d", len(haps2))
	}
}

func TestMJS_PatternValues(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "gain": 0.5})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("map pattern len 1 got %d", len(haps))
	}
	v := haps[0].Value.(map[string]any)
	if v["s"] != "bd" || v["gain"] != 0.5 {
		t.Fatalf("map values s bd gain 0.5 got %v", v)
	}
}
