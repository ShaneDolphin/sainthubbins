// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Regression test: a note held across a cycle boundary must be rendered once.

package audio

import (
	"math"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func peakOf(t *testing.T, p core.Pattern, cycles int) float64 {
	t.Helper()
	samples, err := RenderPatternAudio(p, cycles, 8000)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	var peak float64
	for _, s := range samples {
		if v := math.Abs(float64(s)); v > peak {
			peak = v
		}
	}
	return peak
}

// Combinators that split a query per cycle (Every, LastOf, SometimesBy) cut a
// note held across a cycle boundary into one fragment per cycle, each carrying
// the same Whole. Rendering every fragment summed the note into itself and
// multiplied its amplitude.
func TestHeldNoteRendersOnceUnderCycleSplitting(t *testing.T) {
	held := core.Note(core.Pure("c3")).Slow(core.FractionFromInt(2))
	identity := func(p core.Pattern) core.Pattern { return p }

	bare := peakOf(t, held, 4)
	if bare == 0 {
		t.Fatal("held note rendered silence")
	}

	for _, tc := range []struct {
		name string
		pat  core.Pattern
	}{
		{"Every", held.Every(3, identity)},
		{"LastOf", held.LastOf(3, identity)},
	} {
		got := peakOf(t, tc.pat, 4)
		if math.Abs(got-bare) > 1e-6 {
			t.Errorf("%s: peak %.4f, want %.4f — the held note is being rendered once per cycle fragment",
				tc.name, got, bare)
		}
	}
}

// Fragments must be skipped, not the notes that own them.
func TestFragmentsSkippedButOnsetsKept(t *testing.T) {
	held := core.Note(core.Pure("c3")).Slow(core.FractionFromInt(2))
	identity := func(p core.Pattern) core.Pattern { return p }
	haps := held.Every(3, identity).
		QueryArc(core.FractionFromInt(0), core.FractionFromInt(4))

	var onsets, fragments int
	for _, h := range haps {
		if h.HasOnset() {
			onsets++
		} else {
			fragments++
		}
	}
	if onsets == 0 || fragments == 0 {
		t.Skipf("expected a mix of onsets and fragments, got %d/%d", onsets, fragments)
	}
	if peakOf(t, held.Every(3, identity), 4) == 0 {
		t.Error("all haps were skipped; the onset-carrying fragment must still render")
	}
}
