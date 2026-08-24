// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Thirtieth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_HapContext(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "n": 1})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("map pure empty")
	}
	m := haps[0].Value.(map[string]any)
	if m["s"] != "bd" || m["n"] != 1 {
		t.Fatalf("map s bd n1 got %v", m)
	}
	// Context via Set
	p2 := Pure("a").SetContext(map[string]any{"orbit": 1})
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps2[0].Context["orbit"] != 1 {
		t.Fatalf("context orbit 1 got %v", haps2[0].Context)
	}
}

func TestMJS_PatternTypes(t *testing.T) {
	if !IsPattern(Pure("a")) {
		t.Fatalf("IsPattern pure")
	}
	if IsPattern("a") {
		t.Fatalf("IsPattern string should be false")
	}
	p := Reify("bd")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("reify string empty")
	}
	p2 := Reify(Pure("a"))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("reify pattern empty")
	}
}

func TestMJS_CalculatSteps(t *testing.T) {
	CalculateSteps(true)
	p := Pure("a").WithSteps(func(f Fraction) Fraction { return f })
	if p.Steps == nil {
		t.Fatalf("CalculateSteps true withSteps expected steps")
	}
	CalculateSteps(false)
	p2 := Pure("a").WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	// when disabled, WithSteps returns original (steps still 1 from Pure)
	if p2.Steps == nil {
		t.Fatalf("Pure steps should still be 1 even when disabled")
	}
	CalculateSteps(true)
}
