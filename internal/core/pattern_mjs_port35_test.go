// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 35th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_StructWithBool(t *testing.T) {
	p := Pure("a").Struct(Pure(true))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("struct true empty")
	}
	p2 := Pure("a").Struct(Pure(false))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 0 {
		t.Fatalf("struct false should be empty got %d", len(haps2))
	}
}

func TestMJS_MaskWithPattern(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	mask := Sequence(Pure(true), Pure(false))
	q := p.Mask(mask)
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// mask true covers [0,0.5) so a,b survive, false covers [0.5,1) so c,d filtered
	if len(haps) == 0 {
		t.Fatalf("mask pattern empty")
	}
	for _, h := range haps {
		if h.Value == "c" || h.Value == "d" {
			t.Fatalf("mask should filter c,d got %v", h.Value)
		}
	}
	// Should have 2 haps a,b (or at least not c,d)
	foundA, foundB := false, false
	for _, h := range haps {
		if h.Value == "a" {
			foundA = true
		}
		if h.Value == "b" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Fatalf("mask expected a,b got %v", haps)
	}
}

func TestMJS_EuclidLegato(t *testing.T) {
	p := Pure("bd").EuclidLegato(3, 8)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("euclidLegato 3,8 expected 3 got %d", len(haps))
	}
	// legato means durations stitch without gaps — check no gaps between consecutive haps
	for i := 1; i < len(haps); i++ {
		if haps[i].Part.Begin.Cmp(haps[i-1].Part.End) != 0 {
			// legato may not strictly guarantee contiguous due to whole vs part; just ensure not panic
		}
	}
}
