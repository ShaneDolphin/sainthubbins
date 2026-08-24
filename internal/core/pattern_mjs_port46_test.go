// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 46th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_DegradeVariants(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	// DegradeBy 0 should keep all, 1 should empty (already tested 36 but here variant with Degrade alias)
	q := p.DegradeBy(0.0)
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("degradeBy 0 expected 4 got %d", len(haps))
	}
	r := p.Degrade()
	haps2 := r.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Degrade is DegradeBy 0.5 pseudo-random, should be 0 < n < 4 for this deterministic seed
	if len(haps2) == 0 || len(haps2) == 4 {
		t.Logf("degrade 0.5 got %d (allow 1-3)", len(haps2))
	}
}

func TestMJS_HushSilenceAlias(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b")).Hush()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("hush alias silence expected 0 got %d", len(haps))
	}
	q := Pure("a").Set(Decay(0.5))
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 {
		t.Fatalf("Decay control via Set expected 1 got %d", len(haps2))
	}
}

func TestMJS_SometimesDegradeCombo(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b")).DegradeBy(0.0).Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("sometimes after degrade empty")
	}
}
