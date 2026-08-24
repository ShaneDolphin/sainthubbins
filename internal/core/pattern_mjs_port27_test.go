// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-seventh batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_CatPrime2(t *testing.T) {
	scp := SlowcatPrime(Pure("a"), Pure("b"), Pure("c"))
	if len(scp.QueryArc(FractionFromInt(0), FractionFromInt(3))) == 0 {
		t.Fatalf("slowcatPrime empty")
	}
	fc := FastCat(Pure("a"), Pure("b"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fastcat 2 expected 2")
	}
}

func TestMJS_RevBrak(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("rev 3 expected 3 got %d", len(haps))
	}
	b := Pure("a").Brak()
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("brak empty")
	}
}

func TestMJS_PalindromePolyrhythm(t *testing.T) {
	pal := FastCat(Pure("a"), Pure("b")).Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("palindrome empty")
	}
	pm := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}
