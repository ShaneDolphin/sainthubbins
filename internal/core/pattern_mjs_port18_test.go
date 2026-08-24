// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Eighteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_PlyOffWhen(t *testing.T) {
	p := Pure("a").Ply(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ply 2 empty")
	}
	p2 := Pure("a").Off(0.5, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("off 0.5 expected >=2 got %d", len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p3 := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("when true expected 2 got %d", len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_BrakHurry(t *testing.T) {
	p := Pure("a").Brak()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("brak empty")
	}
	p2 := Pure("a").Hurry(2)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("hurry 2 empty")
	}
}

func TestMJS_SegmentChop(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Segment(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("segment 2 empty")
	}
	p2 := Pure("a").Chop(4)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chop 4 empty")
	}
	p3 := Pure("a").Striate(2)
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("striate 2 empty")
	}
}
