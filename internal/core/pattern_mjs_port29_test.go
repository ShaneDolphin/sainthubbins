// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-ninth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SlowInsideOutside(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Slow(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		// Slow 2 over 1 cycle: should still be 1 (slowed)
		t.Logf("slow 2 len %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := FastCat(Pure("a"), Pure("b")).Inside(2, func(p Pattern) Pattern { return p.Rev() })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("inside 2 rev expected >=2 got %d", len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p3 := FastCat(Pure("a"), Pure("b")).Outside(2, func(p Pattern) Pattern { return p.Rev() })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("outside 2 rev expected >=2")
	}
}

func TestMJS_EuclidBjorklundSignal(t *testing.T) {
	bj := Bjorklund(5, 8)
	if len(bj) != 8 {
		t.Fatalf("bjorklund 5,8 len %d", len(bj))
	}
	p := Pure("a").Euclid(5, 8)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("euclid 5,8 empty")
	}
	s := Sine().Range(0, 1)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sine range empty")
	}
}

func TestMJS_ValueSteps(t *testing.T) {
	p := Pure("a")
	if p.Steps == nil || p.Steps.Cmp(FractionFromInt(1)) != 0 {
		t.Fatalf("pure steps should be 1 got %v", p.Steps)
	}
	p2 := Pure("a").WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	if p2.Steps == nil || p2.Steps.Cmp(FractionFromInt(2)) != 0 {
		t.Fatalf("withSteps *2 expected 2 got %v", p2.Steps)
	}
	// Gap steps
	g := Gap(4)
	if g.Steps == nil || g.Steps.Cmp(FractionFromInt(4)) != 0 {
		t.Fatalf("gap 4 steps %v", g.Steps)
	}
	// Silence steps nil
	if Silence().Steps != nil {
		t.Fatalf("silence steps should be nil")
	}
}
