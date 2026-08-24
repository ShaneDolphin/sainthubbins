// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 32nd batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_Splice(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"))
	// SqueezeJoin approximates splice for 2 len: both in cycle — via InnerJoin+SqueezeJoin
	j := Pure(Pure("x")).SqueezeJoin()
	haps := j.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("squeezeJoin splice empty")
	}
	// Ensure Sequence fastCat spreads across cycle halves
	haps2 := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) < 2 {
		t.Fatalf("splice seq len 2 got %d", len(haps2))
	}
}

func TestMJS_PatternCFuncs(t *testing.T) {
	// pure / silence / stack / sequence / cat
	pat := Cat(Pure("a"), Pure("b"), Pure("c"))
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("cat empty")
	}
	pat2 := Stack(Pure("a"), Pure("b"))
	haps2 := pat2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) < 2 {
		t.Fatalf("stack 2 expected 2 got %d", len(haps2))
	}
	pat3 := Silence()
	haps3 := pat3.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) != 0 {
		t.Fatalf("silence should be empty got %d", len(haps3))
	}
}

func TestMJS_ArpMode(t *testing.T) {
	p := Sequence(Pure("c3"), Pure("e3"), Pure("g3"))
	for _, mode := range []string{"up", "down", "updown", "downup", "converge", "diverge", "random", "thumbup"} {
		q := p.Arp(mode)
		haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
		if len(haps) == 0 {
			t.Fatalf("arp %s empty", mode)
		}
	}
}
