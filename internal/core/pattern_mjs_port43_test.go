// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 43rd batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_HushAndSilence(t *testing.T) {
	p := Pure("a").Hush()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("hush should be silence got %d", len(haps))
	}
	// IsSilence not a Pattern method; check Silence query empty vs Pure non-empty
	if len(Silence().QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence should be empty")
	}
	if len(Pure("a").QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Pure not empty")
	}
}

func TestMJS_ZoomCompress(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Zoom(FractionFromFloat(0.5), FractionFromInt(1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("zoom 0.5-1 empty")
	}
	// Zoom 0.5,1 is like compress to second half + fast? Check not empty
	p2 := Pure("a").Zoom(FractionFromInt(0), FractionFromInt(1))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 {
		t.Fatalf("zoom 0-1 pure 1 got %d", len(haps2))
	}
}

func TestMJS_PlainValues(t *testing.T) {
	p := Pure(42)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value != 42 {
		t.Fatalf("pure 42 got %v", haps)
	}
	p2 := Pure(3.14)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("pure 3.14 empty")
	}
	if toFloat(haps2[0].Value) < 3.13 || toFloat(haps2[0].Value) > 3.15 {
		t.Fatalf("pure 3.14 got %v", haps2[0].Value)
	}
}
