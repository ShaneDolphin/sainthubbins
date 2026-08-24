// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Second batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SetKeep(t *testing.T) {
	// set can set things in objects
	p := Pure(map[string]any{"a": 4, "b": 6}).Set(Pure(map[string]any{"c": 7}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("set empty")
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok {
		t.Fatalf("set not map %T", haps[0].Value)
	}
	if m["c"] != 7 && m["c"] != float64(7) {
		t.Fatalf("set c missing %v", m)
	}
	// keep
	p2 := FastCat(Pure(1), Pure(2)).Keep(Pure(1))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("keep empty")
	}
}

func TestMJS_SlowFastCatPalindrome(t *testing.T) {
	sc := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("slowcat 3 over 3 cycles expected 3 got %d", len(sc.QueryArc(FractionFromInt(0), FractionFromInt(3))))
	}
	// palindrome
	pal := FastCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("palindrome empty")
	}
}

func TestMJS_EuclidArp(t *testing.T) {
	p := Pure("a").Euclid(3, 8)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("euclid empty")
	}
	// arp wraps around with positive and negative indices
	arpPat := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure(0))
	if len(arpPat.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp 0 empty")
	}
	arpNeg := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure(-1))
	if len(arpNeg.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp -1 empty")
	}
}

func TestMJS_RandPick(t *testing.T) {
	// pick via free func
	p := Pick("a", Pure(1))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("pick empty")
	}
	// choose via Pure(0).Choose
	p2 := Pure(0).Choose([]any{"a", "b"})
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("choose empty")
	}
	// rand
	r := Rand().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(r) == 0 {
		t.Fatalf("rand empty")
	}
}

func TestMJS_StackLeftRightCentre(t *testing.T) {
	if len(StackLeft(Pure("a"), Pure("b")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stackLeft empty")
	}
	if len(StackRight(Pure("a"), Pure("b")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stackRight empty")
	}
	if len(StackCentre(Pure("a"), Pure("b")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stackCentre empty")
	}
}
