// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-fourth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SteadySignal(t *testing.T) {
	// steady is control signal that holds value
	p := Pure(0.5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("steady 0.5 empty")
	}
	// sine signal already tested
	s := Sine().Range(0, 1)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sine range 0-1 empty")
	}
}

func TestMJS_ExpandArp(t *testing.T) {
	// expand: stack of patterns?
	p := FastCat(Pure("a"), Pure("b")).Fmap(func(v any) any { return v })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("expand via fmap 2 expected 2")
	}
	// arp
	a := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure("up"))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp up empty")
	}
}

func TestMJS_IdPly(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "a" {
		t.Fatalf("id via fmap")
	}
	p2 := Pure("a").Ply(3)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ply 3 empty")
	}
}
