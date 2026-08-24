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

// TestMiniSingleTokenWeightStripsValue guards against a shortcut that used
// to bypass splitWeight entirely for a lone token: parseSequence special
// cased len(tokens) == 1 and called parseToken(tokens[0]) directly, so
// "bd@3" reached parseToken still carrying "@3", and since parseToken no
// longer understands "@" the raw suffix leaked into the hap value.
func TestMiniSingleTokenWeightStripsValue(t *testing.T) {
	haps := Mini("bd@3").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("got %d haps, want 1: %v", len(haps), haps)
	}
	if haps[0].Value != "bd" {
		t.Errorf("value = %v, want %q — a lone step's @n must not leak into the value", haps[0].Value, "bd")
	}
	if haps[0].Part.Begin.Float64() != 0 || haps[0].Part.End.Float64() != 1 {
		t.Errorf("span = %v..%v, want 0..1 — a single step has no sibling to weight against",
			haps[0].Part.Begin.Float64(), haps[0].Part.End.Float64())
	}
}

// TestMiniBracketedSingleTokenWeightStripsValue covers the same shortcut bug
// reached through a bracket group: parseToken unwraps "[bd@3]" and calls
// parseSequence("bd@3"), landing on the same single-token path.
func TestMiniBracketedSingleTokenWeightStripsValue(t *testing.T) {
	haps := Mini("[bd@3] sd").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2: %v", len(haps), haps)
	}
	var sawBD bool
	for _, h := range haps {
		if h.Value == "bd@3" {
			t.Fatalf("value = %v, the @3 suffix must not leak into the value", h.Value)
		}
		if h.Value == "bd" {
			sawBD = true
			if h.Part.Begin.Float64() != 0 || h.Part.End.Float64() != 0.5 {
				t.Errorf("bd span = %v..%v, want 0..0.5", h.Part.Begin.Float64(), h.Part.End.Float64())
			}
		}
	}
	if !sawBD {
		t.Errorf("no clean %q hap found: %v", "bd", haps)
	}
}

func TestMiniReplicateAddsEqualSteps(t *testing.T) {
	p := Mini("bd!3 sd")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("got %d haps, want 4", len(haps))
	}
	want := []float64{0, 0.25, 0.5, 0.75}
	for i, w := range want {
		if got := haps[i].Part.Begin.Float64(); got != w {
			t.Errorf("hap %d begins at %v, want %v", i, got, w)
		}
	}
	if v, _ := haps[3].Value.(string); v != "sd" {
		t.Errorf("last hap is %v, want sd", haps[3].Value)
	}
}

func TestMiniReplicateMatchesWritingItOut(t *testing.T) {
	a := Mini("bd!3 sd").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	b := Mini("bd bd bd sd").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(a) != len(b) {
		t.Fatalf("!3 gave %d haps, writing it out gave %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Part.Begin.Float64() != b[i].Part.Begin.Float64() {
			t.Errorf("hap %d: !3 at %v, written out at %v",
				i, a[i].Part.Begin.Float64(), b[i].Part.Begin.Float64())
		}
	}
}
