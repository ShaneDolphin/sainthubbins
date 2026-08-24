package draw

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestPianorollBasic(t *testing.T) {
	p := core.Stack(core.Pure("a"), core.Pure("b"))
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	evs := Pianoroll(haps)
	if len(evs) != 2 {
		t.Fatalf("Pianoroll expected 2 got %d", len(evs))
	}
	if evs[0].Duration <= 0 {
		t.Fatalf("Pianoroll duration expected >0")
	}
}

func TestSpiralBasic(t *testing.T) {
	p := core.FastCat(core.Pure("a"), core.Pure("b"), core.Pure("c"))
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	evs := Spiral(haps)
	if len(evs) != 3 {
		t.Fatalf("Spiral expected 3 got %d", len(evs))
	}
	// Angle should be in [0, 2pi)
	for _, e := range evs {
		if e.Time < 0 || e.Time >= 6.28318530718+1e-9 {
			t.Fatalf("Spiral angle out of range %v", e.Time)
		}
	}
}

func TestToJSON(t *testing.T) {
	p := core.Pure("a")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	js := ToJSON(haps)
	if len(js) == 0 || js[0] != '[' {
		t.Fatalf("ToJSON expected array got %s", js)
	}
}
