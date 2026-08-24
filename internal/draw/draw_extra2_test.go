// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package draw

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestDrawPianorollSilence(t *testing.T) {
	p := core.Silence()
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	pr := Pianoroll(haps)
	if len(pr) != 0 {
		t.Fatalf("silence pianoroll expected 0 got %d", len(pr))
	}
	js := ToJSON(haps)
	if js != "[]" {
		t.Fatalf("silence json expected [] got %q", js)
	}
}

func TestDrawSpiralSingle(t *testing.T) {
	p := core.Pure(map[string]any{"note": "c4"})
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	sp := Spiral(haps)
	if len(sp) != 1 {
		t.Fatalf("spiral single expected 1 got %d", len(sp))
	}
}
