// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-second batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_StackCatSteps(t *testing.T) {
	s := Stack(Pure("a"), Pure("b"))
	if s.Steps != nil {
		t.Fatalf("stack steps should be nil")
	}
	fc := FastCat(Pure("a"), Pure("b"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fastcat 2 expected 2")
	}
	// stack with steps via WithSteps
	p := Pure("a").WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	if p.Steps == nil {
		t.Fatalf("withSteps expected steps")
	}
}

func TestMJS_TimeCatPolymeterSteps(t *testing.T) {
	tc := TimeCat(FractionFromInt(1), Pure("a"), FractionFromInt(1), Pure("b"))
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("timeCat empty")
	}
	pm := Polymeter(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}

func TestMJS_RandDegradeSometimes(t *testing.T) {
	r := Rand().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(r) == 0 {
		t.Fatalf("rand empty")
	}
	p := Pure("a").DegradeBy(0.0)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("degrade 0 empty")
	}
	p2 := Pure("a").SometimesBy(0.5, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	// may be original or fast, but non-empty
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sometimesBy 0.5 empty")
	}
}
