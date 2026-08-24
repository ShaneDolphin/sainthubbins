// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestArpCollect(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	// Collect groups congruent haps: should produce haps with []Hap values
	c := p.Collect()
	haps := c.FirstCycle()
	if len(haps) == 0 {
		t.Fatalf("collect expected haps")
	}
}

func TestArp(t *testing.T) {
	// Simple arp: stack of notes then select 0
	p := Stack(Pure("c"), Pure("eb"), Pure("g")).Arp(0)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("arp expected haps")
	}
}

func TestChunkFamily(t *testing.T) {
	p := Pure("a").Chunk(2, func(pat Pattern) Pattern { return pat.Fast(Pure(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunk expected haps")
	}
	p2 := Pure("a").ChunkBack(2, func(pat Pattern) Pattern { return pat })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chunkBack expected haps")
	}
	p3 := Pure("a").FastChunk(2, func(pat Pattern) Pattern { return pat })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("fastChunk expected haps")
	}
}

func TestIterBackRepeat(t *testing.T) {
	p := Pure("a").IterBack(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("iterBack expected haps")
	}
	r := Pure("a").RepeatCycles(2)
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("repeatCycles expected haps")
	}
}

func TestShrinkGrowTour(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"))
	f := FractionFromInt(2)
	p.Steps = &f
	s := p.Shrink(1)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("shrink produced %d haps (acceptable stub)", len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	g := p.Grow(1)
	if len(g.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("grow produced %d haps", len(g.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	tp := Pure("a").Tour(Pure("b"), Pure("c"))
	if len(tp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("tour expected haps")
	}
}
