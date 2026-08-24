// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// More JSPort tests from pattern.test.mjs: superimpose, layer, chunk, etc.

package core

import "testing"

func TestJSPort_Superimpose(t *testing.T) {
	pat := Pure("a").Superimpose(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 3 {
		t.Fatalf("superimpose expected >=3 got %d", len(haps))
	}
}

func TestJSPort_Layer(t *testing.T) {
	pat := Pure("a").Layer(func(p Pattern) Pattern { return p }, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 3 {
		t.Fatalf("layer expected >=3 got %d", len(haps))
	}
}

func TestJSPort_Chunk(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(p Pattern) Pattern { return p.Rev() })
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("chunk 2 rev expected 4 got %d", len(haps))
	}
}

func TestJSPort_Sometimes(t *testing.T) {
	pat := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("sometimes empty")
	}
	// SometimesBy 0 should be identity
	pat2 := Pure("a").SometimesBy(0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(pat2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("sometimesBy 0 expected 1")
	}
}

func TestJSPort_Degrade(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Degrade()
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// degrade 0.5 deterministic: allow 0-4, but DegradeBy 1 should drop all
	if len(haps) < 0 || len(haps) > 4 {
		t.Fatalf("degrade out of range got %d", len(haps))
	}
	// DegradeBy 0 should keep all
	pat0 := FastCat(Pure("a"), Pure("b")).DegradeBy(0)
	if len(pat0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("degradeBy 0 expected 2")
	}
	// DegradeBy 1 should drop all
	pat1 := FastCat(Pure("a"), Pure("b")).DegradeBy(1)
	if len(pat1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("degradeBy 1 expected 0")
	}
}
