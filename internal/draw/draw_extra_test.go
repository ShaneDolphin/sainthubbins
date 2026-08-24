// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package draw

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestDrawPianorollStack(t *testing.T) {
	p := core.Stack(core.Pure(map[string]any{"note": "c4"}), core.Pure(map[string]any{"note": "e4"}))
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	pr := Pianoroll(haps)
	if len(pr) == 0 {
		t.Fatalf("pianoroll empty")
	}
	// Check time within 0-2
	for _, ev := range pr {
		if ev.Time < 0 || ev.Time >= 1 {
			t.Fatalf("time %v out of 0-1", ev.Time)
		}
	}
}

func TestDrawSpiral(t *testing.T) {
	p := core.FastCat(core.Pure(map[string]any{"n": 60}), core.Pure(map[string]any{"n": 64}))
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	sp := Spiral(haps)
	if len(sp) == 0 {
		t.Fatalf("spiral empty")
	}
}
