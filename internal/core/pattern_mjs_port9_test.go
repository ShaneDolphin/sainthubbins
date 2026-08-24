// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Ninth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_StructAllMask(t *testing.T) {
	p := Pure("a").Struct(Pure(true))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("struct true empty")
	}
	p2 := Pure("a").Mask(Pure(true))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("mask true empty")
	}
	p3 := Pure("a").StructAll(Pure(true))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("structAll empty")
	}
}

func TestMJS_SometimesOften(t *testing.T) {
	p := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("sometimes empty")
	}
	p2 := Pure("a").SometimesBy(0.75, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("often (0.75) empty")
	}
	p3 := Pure("a").SometimesBy(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("rarely (0.25) maybe empty but ok")
	}
}

func TestMJS_PlayCloned(t *testing.T) {
	p := Pure("a").Off(0.1, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("off expected >=2 got %d", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := Pure("a").Echo(2, 0.25, 0.5)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("echo empty")
	}
}
