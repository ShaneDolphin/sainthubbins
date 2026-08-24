// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// onsets returns each hap's value and start time for one cycle.
func onsets(p core.Pattern) map[string]float64 {
	out := map[string]float64{}
	for _, h := range p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)) {
		if s, ok := h.Value.(string); ok {
			out[s] = h.Part.Begin.Float64()
		}
	}
	return out
}

func TestMiniWeightElongatesAStep(t *testing.T) {
	got := onsets(Mini("bd@3 sd"))
	if len(got) != 2 {
		t.Fatalf("got %d haps, want 2: %v", len(got), got)
	}
	if got["bd"] != 0 {
		t.Errorf("bd starts at %v, want 0", got["bd"])
	}
	if got["sd"] != 0.75 {
		t.Errorf("sd starts at %v, want 0.75 — bd should hold three quarters", got["sd"])
	}
}

func TestMiniWeightIsNotAlwaysThree(t *testing.T) {
	got := onsets(Mini("bd@1 sd@3"))
	if got["sd"] != 0.25 {
		t.Errorf("sd starts at %v, want 0.25", got["sd"])
	}
}

func TestMiniWithoutWeightIsUnchanged(t *testing.T) {
	got := onsets(Mini("bd sd"))
	if got["bd"] != 0 || got["sd"] != 0.5 {
		t.Errorf("unweighted sequence changed: %v", got)
	}
}
