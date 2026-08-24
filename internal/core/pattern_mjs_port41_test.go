// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 41st batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ReifyAndIsPattern(t *testing.T) {
	if !IsPattern(Pure("a")) {
		t.Fatalf("IsPattern pure true")
	}
	if IsPattern("a") {
		t.Fatalf("IsPattern string false")
	}
	if IsPattern(42) {
		t.Fatalf("IsPattern int false")
	}
	p := Reify("bd")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value != "bd" {
		t.Fatalf("Reify string bd got %v", haps)
	}
	p2 := Reify(Pure("a"))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 || haps2[0].Value != "a" {
		t.Fatalf("Reify pattern got %v", haps2)
	}
}

func TestMJS_SilenceGap(t *testing.T) {
	s := Silence()
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("silence 0 haps got %d", len(haps))
	}
	g := Gap(2)
	haps2 := g.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 0 {
		t.Fatalf("gap 2 still 0 haps got %d", len(haps2))
	}
	if g.Steps == nil || g.Steps.Cmp(FractionFromInt(2)) != 0 {
		t.Fatalf("gap 2 steps 2 got %v", g.Steps)
	}
}

func TestMJS_PatternWithLoc(t *testing.T) {
	p := Pure("a").WithLoc(1, 2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value != "a" {
		t.Fatalf("WithLoc value a got %v", haps)
	}
	if _, ok := haps[0].Context["locations"]; !ok {
		t.Fatalf("WithLoc Context locations missing got %v", haps[0].Context)
	}
}
