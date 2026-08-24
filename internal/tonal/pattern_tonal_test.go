// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package tonal

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestScalePattern(t *testing.T) {
	p := ScalePattern("C:major")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 7 {
		t.Fatalf("C:major expected 7 got %d", len(haps))
	}
	if haps[0].Value != "C" {
		t.Fatalf("first note C got %v", haps[0].Value)
	}
}

func TestChordPattern(t *testing.T) {
	p := ChordPattern("Cmaj7")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("Cmaj7 expected 4 got %d", len(haps))
	}
}

func TestTransposePattern(t *testing.T) {
	p := core.FastCat(core.Pure("c4"), core.Pure("e4"))
	tp := TransposePattern(p, 2)
	haps := tp.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("transpose expected 2 got %d", len(haps))
	}
	if haps[0].Value != "D4" {
		t.Fatalf("c4+2 expected D4 got %v", haps[0].Value)
	}
}
