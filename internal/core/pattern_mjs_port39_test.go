// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 39th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_StackWithRest(t *testing.T) {
	p := Stack(Silence(), Pure("a"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Stack silence + a: only a survives (silence has 0 haps but Stack still queries both branches)
	if len(haps) != 1 || haps[0].Value != "a" {
		t.Fatalf("stack silence+a expected 1 a got %v", haps)
	}
}

func TestMJS_CatWithSilence(t *testing.T) {
	// Cat cycles per cycle: first cycle a, second silence, third a again etc. Over 2 cycles: 1 hap (a)
	cat := Cat(Pure("a"), Silence())
	haps := cat.QueryArc(FractionFromInt(0), FractionFromInt(2))
	// cycle0 a (1 hap), cycle1 silence (0) => 1
	if len(haps) != 1 {
		t.Fatalf("cat a silence over 2 cycles expected 1 got %d %v", len(haps), haps)
	}
}

func TestMJS_SuperimposeWithOff(t *testing.T) {
	p := Pure("a").Superimpose(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// superimpose stacks original + fn — should have at least 2 layers (with some overlap)
	if len(haps) == 0 {
		t.Fatalf("superimpose empty")
	}
	// Should have more than plain pure (which would be 1) due to stacking
	if len(haps) < 2 {
		t.Fatalf("superimpose expected >=2 got %d", len(haps))
	}
}
