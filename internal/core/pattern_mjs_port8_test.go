// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Eighth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_HushGapSilence(t *testing.T) {
	p := Pure("a").Hush()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("hush expected 0 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	if len(Gap(2).QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("gap 2 expected 0")
	}
	if len(Silence().QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("silence expected 0")
	}
}

func TestMJS_DefragResetPoly(t *testing.T) {
	p := Pure(Pure("a")).PolyJoin()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polyJoin empty")
	}
	p2 := Pure(Pure("a")).ResetJoin()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("resetJoin empty")
	}
	p3 := FastCat(Pure("a"), Pure("a")).Defragment()
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("defrag empty")
	}
}

func TestMJS_BinaryOps(t *testing.T) {
	// binary ops via Add etc already tested, add bitwise
	p := Pure(3).Fmap(func(v any) any { return v.(int) & 1 })
	if v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(int); v != 1 {
		t.Fatalf("band 3&1 expected 1 got %v", v)
	}
	p2 := Pure(2.0).Range(0, 10)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("range empty")
	}
}
