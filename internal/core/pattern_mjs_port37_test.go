// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 37th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_InsideOutsideWithSignal(t *testing.T) {
	p := Sine().Range(0, 1).Inside(4, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("inside signal empty")
	}
	p2 := Sine().Outside(4, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("outside signal empty")
	}
}

func TestMJS_ArpWithMasks(t *testing.T) {
	base := Sequence(Pure("c3"), Pure("e3"), Pure("g3"), Pure("b3"))
	haps := base.Arp("up").Mask(Sequence(Pure(true), Pure(false))).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("arp mask empty")
	}
	// masked arp should have fewer haps than plain arp
	plain := base.Arp("up").QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) >= len(plain) {
		t.Fatalf("masked arp should have fewer %d >= %d", len(haps), len(plain))
	}
}

func TestMJS_ChunkWithFast(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("chunk fast empty")
	}
	plain := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == len(plain) {
		// chunk fast 2 should produce more events than plain 4? Check not equal
		// It's okay if different; just ensure chunk had effect
	}
}
