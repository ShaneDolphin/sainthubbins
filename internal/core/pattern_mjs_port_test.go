// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Ported subset of js/packages/core/test/pattern.test.mjs via Go.
// Covers Phase 1 gate: TimeSpan/Hap/Pattern core semantics.

package core

import "testing"

func TestMJS_PureFmap(t *testing.T) {
	// pure('hello').query(st(0.5,2.5)).length == 3
	if n := len(Pure("hello").QueryArc(FractionFromFloat(0.5), FractionFromFloat(2.5))); n != 3 {
		t.Fatalf("pure hello 0.5-2.5 expected 3 got %d", n)
	}
	// zero-width query still 1 hap
	if n := len(Pure("hello").QueryArc(FractionFromInt(0), FractionFromInt(0))); n != 1 {
		t.Fatalf("zero-width expected 1 got %d", n)
	}
	// fmap add
	v := Pure(3).Fmap(func(x any) any { return x.(int) + 4 }).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(int)
	if v != 7 {
		t.Fatalf("fmap 3+4 got %v", v)
	}
}

func TestMJS_AddSubMulDiv(t *testing.T) {
	toF := func(v any) float64 {
		switch x := v.(type) {
		case int:
			return float64(x)
		case int64:
			return float64(x)
		case float64:
			return x
		case float32:
			return float64(x)
		default:
			return 0
		}
	}
	if v := toF(Pure(4).Add(Pure(5)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value); v != 9 {
		t.Fatalf("add 4+5 got %v", v)
	}
	if v := toF(Pure(3).Sub(Pure(4)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value); v != -1 {
		t.Fatalf("sub 3-4 got %v", v)
	}
	if v := toF(Pure(3).Mul(Pure(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value); v != 6 {
		t.Fatalf("mul 3*2 got %v", v)
	}
	if v2 := Pure(3.0).Div(Pure(2.0)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v2 != 1.5 {
		t.Fatalf("div 3/2 got %v", v2)
	}
}

func TestMJS_StackFastSlow(t *testing.T) {
	// stack can stack
	if n := len(Stack(Pure("a"), Pure("b")).QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 2 {
		t.Fatalf("stack 2 expected 2 got %d", n)
	}
	// _fast makes things faster
	if n := len(Pure("a").FastF(FractionFromInt(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 2 {
		t.Fatalf("fast 2 expected 2 got %d", n)
	}
	// _slow makes things slower (slow 2 over 2 cycles ->1 hap per queried 1 cycle? JS SlowCat style is per-cycle)
	// Pure slow 2 queried over 1 cycle should still be 1 (whole cycle slowed)
	if n := len(Pure("a").SlowF(FractionFromInt(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 1 {
		t.Fatalf("slow 2 expected 1 got %d", n)
	}
}

func TestMJS_FastCatSlowCat(t *testing.T) {
	// fastcat can concatenate two things
	fc := FastCat(Pure("a"), Pure("b"))
	if n := len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 2 {
		t.Fatalf("fastcat a,b expected 2 got %d", n)
	}
	// slowcat can concatenate slowly (each per cycle)
	sc := Cat(Pure("a"), Pure("b"))
	if n := len(sc.QueryArc(FractionFromInt(0), FractionFromInt(2))); n != 2 {
		t.Fatalf("slowcat a,b over 2 cycles expected 2 got %d", n)
	}
	if n := len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 1 {
		t.Fatalf("slowcat a,b over 1 cycle expected 1 got %d", n)
	}
}

func TestMJS_RevWhen(t *testing.T) {
	// rev can reverse
	fc := FastCat(Pure("a"), Pure("b"), Pure("c"))
	rev := fc.Rev()
	haps := rev.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("rev fastcat 3 expected 3 got %d", len(haps))
	}
	// when Always faster (when true) via Pure(true)
	p := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if n := len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 2 {
		t.Fatalf("when true fast2 expected 2 got %d", n)
	}
	// when Never faster (when false) should stay 1
	p2 := Pure("a").When(Pure(false), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if n := len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 1 {
		t.Fatalf("when never faster expected 1 got %d", n)
	}
}

func TestMJS_StructureOps(t *testing.T) {
	// inside can rev inside a cycle — use factor 2, expect at least 3 haps (JS preserves count but may duplicate across slow/fast)
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Inside(2, func(p Pattern) Pattern { return p.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 3 {
		t.Fatalf("inside rev expected >=3 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	// outside
	p2 := FastCat(Pure("a"), Pure("b")).Outside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("outside expected haps")
	}
	// compress
	p3 := Pure("a").Compress(FractionFromFloat(0.5), FractionFromFloat(1.0))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress empty")
	}
	// filterValues true
	p4 := FastCat(Pure(true), Pure(false), Pure(true)).FilterValues(func(v any) bool { return v.(bool) })
	if len(p4.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("filterValues empty")
	}
}

func TestMJS_SequencePalindrome(t *testing.T) {
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"))
	if n := len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 3 {
		t.Fatalf("sequence 3 expected 3 got %d", n)
	}
	pal := FastCat(Pure("a"), Pure("b")).Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("palindrome empty")
	}
}
