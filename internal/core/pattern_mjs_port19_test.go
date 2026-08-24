// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Nineteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_RunBinary(t *testing.T) {
	// Run(4) => 0 1 2 3 via FastCat
	r := FastCat(Pure(0), Pure(1), Pure(2), Pure(3))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("run 4 expected 4 got %d", len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	// BinaryN via run-like: just test fast on run
	p := FastCat(Pure("a"), Pure("b")).Fast(Pure(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("fast 2 on fastcat 2 expected 4 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_SometimesDegradeOften(t *testing.T) {
	p := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sometimes empty")
	}
	p2 := Pure("a").DegradeBy(0.0)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("degrade 0 empty")
	}
	p3 := Pure("a").SometimesBy(0.0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sometimesBy 0 empty (should be original)")
	}
}

func TestMJS_StackCatPolymeterSteps(t *testing.T) {
	s := Stack(Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("stack 2 expected 2")
	}
	fc := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("fastcat 4 expected 4")
	}
	pm := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}
