// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestHushSilenceGapExtra(t *testing.T) {
	p := Pure("a").Hush()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush expected 0")
	}
	s := Silence()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence expected 0")
	}
	g := Gap(2)
	if g.Steps == nil || g.Steps.Cmp(FractionFromInt(2)) != 0 {
		t.Fatalf("Gap steps expected 2")
	}
}

func TestIterChunkExtra(t *testing.T) {
	p := Pure("a").Chunk(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk expected haps")
	}
	p2 := Pure("a").IterBack(2)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("IterBack expected haps")
	}
}

func TestPickAliasExtra(t *testing.T) {
	_ = Pure(0).Choose([]any{"a"})
}
