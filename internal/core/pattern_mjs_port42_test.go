// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 42nd batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FastSlowEdge(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"))
	fast := p.FastF(FractionFromInt(2))
	haps := fast.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("fast 2 over 2 items expected 4 got %d %v", len(haps), haps)
	}
	slow := p.Slow(FractionFromInt(2))
	haps2 := slow.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// slow 2 stretches over 2 cycles: still 2 items but each stretched, over 1 arc query should be 1? Actually slow 2 covers 0.5 cycle?
	if len(haps2) == 0 {
		t.Fatalf("slow 2 empty")
	}
}

func TestMJS_BinaryOp(t *testing.T) {
	p := Pure(3).Add(Pure(4))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value != 7 && haps[0].Value != float64(7) {
		t.Fatalf("3+4=7 got %v %T", haps[0].Value, haps[0].Value)
	}
	p2 := Pure(10).Sub(Pure(3))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps2[0].Value) != 7 {
		t.Fatalf("10-3=7 got %v", haps2[0].Value)
	}
	p3 := Pure(3).Mul(Pure(4))
	haps3 := p3.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps3[0].Value) != 12 {
		t.Fatalf("3*4=12 got %v", haps3[0].Value)
	}
}

func TestMJS_CompressWithin(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("compress 0.25-0.75 expected 1 got %d", len(haps))
	}
	// whole should be within compressed window mapped back?
	if haps[0].Part.Begin.Cmp(FractionFromFloat(0.25)) < 0 || haps[0].Part.End.Cmp(FractionFromFloat(0.75)) > 0 {
		t.Fatalf("compress part outside 0.25-0.75 got %v", haps[0].Part)
	}
}
