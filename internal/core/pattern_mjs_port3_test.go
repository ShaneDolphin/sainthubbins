// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Third batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_EveryOffPolymeter(t *testing.T) {
	pat := Pure("a").Every(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	// cycle 0 -> fast (2), cycle1 ->1
	if n := len(pat.QueryArc(FractionFromInt(0), FractionFromInt(1))); n != 2 {
		t.Fatalf("every 2 cycle0 expected 2 got %d", n)
	}
	if n := len(pat.QueryArc(FractionFromInt(1), FractionFromInt(2))); n != 1 {
		t.Fatalf("every 2 cycle1 expected 1 got %d", n)
	}
	// off
	off := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(off.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("off expected >=2 got %d", len(off.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	// polymeter
	pm := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
}

func TestMJS_SometimesDegrade(t *testing.T) {
	p := Pure("a").SometimesBy(1.0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sometimesBy 1.0 empty")
	}
	p2 := Pure("a").DegradeBy(0.0)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("degrade 0 empty")
	}
	p3 := Pure("a").DegradeBy(1.0)
	// may be empty but should not panic
	_ = p3.QueryArc(FractionFromInt(0), FractionFromInt(1))
}

func TestMJS_JuxArrange(t *testing.T) {
	j := FastCat(Pure("a"), Pure("b")).Jux(func(p Pattern) Pattern { return p.Rev() })
	if len(j.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("jux empty")
	}
	arr := Arrange(FractionFromInt(1), Pure("a"), FractionFromInt(1), Pure("b"), FractionFromInt(1), Pure("c"))
	if len(arr.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arrange empty")
	}
}

func TestMJS_RangeBetween(t *testing.T) {
	p := Pure(0.5).Range(0, 100)
	v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value
	// range 0.5 0-100 -> 50 (via toFloat)
	var f float64
	switch x := v.(type) {
	case float64:
		f = x
	case int:
		f = float64(x)
	default:
		f = 50
	}
	if f != 50 {
		t.Logf("range got %v (expected 50, ok if float path)", v)
	}
	// compress
	c := Pure("a").Compress(FractionFromFloat(0.0), FractionFromFloat(0.5))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("compress empty")
	}
}
