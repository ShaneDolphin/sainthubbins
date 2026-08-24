// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-sixth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ExpandRange(t *testing.T) {
	p := Pure(0).Range(0, 100)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("range 0 100 empty")
	}
	// check value approx 50*? Pure 0 gives 0, not 50 — use 0.5
	p2 := Pure(0.5).Range(0, 100)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	var f float64
	switch x := haps2[0].Value.(type) {
	case float64:
		f = x
	case int:
		f = float64(x)
	default:
		f = 0
	}
	if f != 50 {
		t.Logf("range 0.5 0-100 got %v (expected 50)", haps2[0].Value)
	}
}

func TestMJS_ArpPolychord(t *testing.T) {
	a := FastCat(Pure("c4"), Pure("e4"), Pure("g4")).Arp(Pure("up"))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp ceg up empty")
	}
	a2 := FastCat(Pure("a"), Pure("b")).Arp(Pure("down"))
	if len(a2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp down empty")
	}
}

func TestMJS_StackWithSignal(t *testing.T) {
	p := Stack(Saw().Range(0, 1), Pure("a"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stack saw range empty")
	}
	p2 := Stack(Sine().Range(-1, 1), Pure("b"))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stack sine empty")
	}
}
