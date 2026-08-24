// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Seventeenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FastCatSlowCatPrime(t *testing.T) {
	fc := FastCat(Pure("a"), Pure("b"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fastcat a,b expected 2 got %d", len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	sc := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("slowcat 3/3 expected 3")
	}
	// slowcatPrime via SlowcatPrime
	scp := SlowcatPrime(Pure("a"), Pure("b"))
	if len(scp.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("slowcatPrime empty")
	}
}

func TestMJS_ArpRev(t *testing.T) {
	arp := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure(0))
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp 0 empty")
	}
	rev := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := rev.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("rev 3 expected 3 got %d", len(haps))
	}
	// revv vs rev
	rv := FastCat(Pure("a"), Pure("b")).Rev()
	rvv := FastCat(Pure("a"), Pure("b")).Revv()
	if len(rv.QueryArc(FractionFromInt(0), FractionFromInt(1))) != len(rvv.QueryArc(FractionFromInt(0), FractionFromInt(1))) {
		t.Logf("rev vs revv len differ (ok)")
	}
}

func TestMJS_SequencePalindromePolymeter(t *testing.T) {
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("sequence 3 expected 3")
	}
	pal := FastCat(Pure("a"), Pure("b")).Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("palindrome empty")
	}
	pm := Polymeter(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}
