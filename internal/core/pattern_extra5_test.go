// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestStackCatArrange(t *testing.T) {
	s := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3 expected 3")
	}
	// Cat is SlowCat: each pat per cycle, so 3 cycles = 3 haps
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	haps := c.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps) != 3 {
		t.Fatalf("Cat 3 over 3 cycles expected 3 got %d", len(haps))
	}
	// FastCat gives Sequence: all pats within one cycle
	fc := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3 expected 3 got %d", len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	// Fast cat 2x should give 4 over 1 cycle (2 pats *2)
	c2 := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2))
	if len(c2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastCat*2 expected 4 got %d", len(c2.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestRevWhenOff(t *testing.T) {
	p := Pure("a").Rev()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rev expected haps")
	}
	p2 := Pure("a").When(func(b bool) bool { return b }, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When expected haps")
	}
	p3 := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Off expected haps")
	}
}

func TestStructEuclidBjorklund(t *testing.T) {
	p := Pure("a").Struct(Pure(true))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true expected haps")
	}
	p2 := Pure("a").Euclid(3, 8)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid expected haps")
	}
	// Bjorklund pure
	bj := Bjorklund(3, 8)
	if len(bj) != 8 {
		t.Fatalf("Bjorklund len %d", len(bj))
	}
}
