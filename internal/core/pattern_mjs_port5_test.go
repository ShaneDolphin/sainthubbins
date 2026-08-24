// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Fifth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_FocusSplice(t *testing.T) {
	p := Pure("a").Focus(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("focus empty")
	}
	// splice: splice 2 "a b" "c" style -> Splice free?
	// Use generic Slice/Splice via pattern_samples
	s := Slice(2, Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("slice 2 a b empty")
	}
}

func TestMJS_WeaveMorph(t *testing.T) {
	p := Pure("a").Weave(1, Pure("b"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("weave empty")
	}
	p2 := S("bd").WeaveWith(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) }, func(p Pattern) Pattern { return p })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("weaveWith empty")
	}
	// morph 3 args: from, to, by
	m := Pure("a").Morph(Pure("a"), Pure("b"), FractionFromFloat(0.5))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("morph empty")
	}
}

func TestMJS_VlpfXFade(t *testing.T) {
	p := S("bd").Vlpf(Pure(800))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("vlpf empty")
	}
	p2 := S("bd").XFade(Pure("sd"), Pure(0.5))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("xfade empty")
	}
}

func TestMJS_ArpChunk(t *testing.T) {
	a := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure("up"))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp up empty")
	}
	c := Pure("a").Chunk(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunk empty")
	}
	ib := Pure("a").IterBack(2)
	if len(ib.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("iterBack empty")
	}
}
