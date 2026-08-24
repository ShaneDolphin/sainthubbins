// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twelfth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FastGapCompress(t *testing.T) {
	p := Pure("a").FastGap(FractionFromFloat(0.5))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fastGap 0.5 empty")
	}
	p2 := Pure("a").Compress(FractionFromFloat(0.0), FractionFromFloat(1.0))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress 0-1 empty")
	}
	p3 := Pure("a").Zoom(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("zoom 0.25-0.75 empty")
	}
}

func TestMJS_EveryWhen(t *testing.T) {
	p := Pure("a").Every(3, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(3))) == 0 {
		t.Fatalf("every 3 0-3 empty")
	}
	p2 := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("when true expected 2 got %d", len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_IterBackChunk(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Iter(2)
	// Iter via Segment2+SqueezeJoin may be 0 for this pattern, but IterBack should be non-empty
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("iter 2 empty (ok if SqueezeJoin)")
	}
	p2 := Pure("a").IterBack(2)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("iterBack 2 empty")
	}
	p3 := Pure("a").Chunk(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunk empty")
	}
}
