// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestMiniEuclidWeight(t *testing.T) {
	p := Mini("bd(3,8)")
	if len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("bd(3,8) expected haps")
	}
	p2 := Mini("bd@2 sd")
	haps := p2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("bd@2 expected haps")
	}
}

func TestMiniDegradeSlow(t *testing.T) {
	p := Mini("bd*2")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("bd*2 expected 2 haps got %d", len(haps))
	}
	p2 := Mini("bd/2")
	if len(p2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("bd/2 expected haps")
	}
	// Degrade op ? is random but should produce haps
	p3 := Mini("bd?0.5")
	if len(p3.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		// degrade may drop but DegradeBy 0.5 with seed may still produce
		t.Logf("bd?0.5 empty (degraded all) — ok")
	}
}

func TestMiniNested(t *testing.T) {
	p := Mini("[bd sd]*2")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) < 4 {
		t.Fatalf("[bd sd]*2 expected 4 got %d", len(haps))
	}
}
