// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Fourth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_TimeCatStepCat(t *testing.T) {
	// timeCat with durations
	tc := TimeCat(FractionFromInt(1), Pure("a"), FractionFromInt(2), Pure("b"))
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("timeCat empty")
	}
	// stepCat
	sc := StepCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stepCat empty")
	}
	// polymeter
	pm := Polymeter(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}

func TestMJS_ChopStriate(t *testing.T) {
	p := Pure("a").Chop(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chop empty")
	}
	p2 := Pure("a").Striate(3)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("striate empty")
	}
	p3 := Slice(2, Pure("a"), Pure(0))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("slice empty")
	}
}

func TestMJS_WithinCompress(t *testing.T) {
	p := Pure("a").Within(0.5, 1.0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("within empty")
	}
	p2 := Pure("a").Compress(FractionFromFloat(0.0), FractionFromFloat(0.5))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress empty")
	}
	p3 := Pure("a").Zoom(FractionFromFloat(0.5), FractionFromFloat(1.0))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("zoom empty")
	}
}

func TestMJS_SegmentBrak(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b")).Segment(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("segment empty")
	}
	p2 := Pure("a").Brak()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("brak empty")
	}
	p3 := Pure("a").Invert()
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("invert empty")
	}
}
