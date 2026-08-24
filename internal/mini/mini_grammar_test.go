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

// cycleValue returns the single value sounding in cycle n of p, or nil if
// none. Angle-bracket alternation is one value per whole cycle, so this is
// enough to check what each cycle's slot resolved to.
func cycleValue(p core.Pattern, n int64) any {
	haps := p.QueryArc(core.FractionFromInt(n), core.FractionFromInt(n+1))
	if len(haps) == 0 {
		return nil
	}
	return haps[0].Value
}

// TestMiniAngleBracketReplicateExpandsAlternatives guards the direct
// parseToken('<...>' branch) call site, which does not go through
// parseSequence and so never saw splitReplicate/splitWeight until this
// fix. Before the fix, "bd!3" reached parseToken still carrying "!3", and
// since parseToken no longer understands "!" the raw suffix leaked into
// the value instead of expanding into three alternative cycles.
func TestMiniAngleBracketReplicateExpandsAlternatives(t *testing.T) {
	p := Mini("<bd!3 sd>")
	want := []string{"bd", "bd", "bd", "sd"}
	for i, w := range want {
		got := cycleValue(p, int64(i))
		if got != w {
			t.Errorf("cycle %d = %v, want %q", i, got, w)
		}
	}
}

// TestMiniAngleBracketWeightIsStrippedNotApplied guards the same call site
// for "@": a lone alternation slot is always exactly one cycle, so the
// weight has nothing to be relative to and must be discarded — but the
// value still needs to come out clean, not carrying "@3".
func TestMiniAngleBracketWeightIsStrippedNotApplied(t *testing.T) {
	p := Mini("<bd@3 sd>")
	want := []string{"bd", "sd"}
	for i, w := range want {
		got := cycleValue(p, int64(i))
		if got != w {
			t.Errorf("cycle %d = %v, want %q", i, got, w)
		}
	}
}

// TestMiniAngleBracketPlainAlternationUnchanged is the regression guard:
// every existing "<...>" pattern goes through the same branch this fix
// touched, so a plain alternation with no suffixes must still alternate
// exactly as before.
func TestMiniAngleBracketPlainAlternationUnchanged(t *testing.T) {
	p := Mini("<bd sd>")
	want := []string{"bd", "sd", "bd", "sd"}
	for i, w := range want {
		got := cycleValue(p, int64(i))
		if got != w {
			t.Errorf("cycle %d = %v, want %q", i, got, w)
		}
	}
}

// wantValues asserts haps' values, in order, against want.
func wantValues(t *testing.T, haps []core.Hap, want []any) {
	t.Helper()
	if len(haps) != len(want) {
		t.Fatalf("got %d haps, want %d: %v", len(haps), len(want), haps)
	}
	for i, w := range want {
		if haps[i].Value != w {
			t.Errorf("hap %d value = %v, want %v", i, haps[i].Value, w)
		}
	}
}

// TestMiniRangeReplicateBeforeExpands guards the ".." range operator's
// surrounding-token handling in parseSequence: a token before the range is
// a sequence step like any other, so "!" must still expand into siblings
// rather than leaking "bd!2" as a literal value.
func TestMiniRangeReplicateBeforeExpands(t *testing.T) {
	haps := Mini("bd!2 0 .. 3").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	wantValues(t, haps, []any{"bd", "bd", 0, 1, 2, 3})
}

// TestMiniRangeWeightAfterIsStripped covers the other ordering — a token
// after the range — and the other suffix: "@" must be stripped from the
// value even though (as before this fix) this branch has never applied its
// weight to the timing.
func TestMiniRangeWeightAfterIsStripped(t *testing.T) {
	haps := Mini("0 .. 3 sd@2").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	wantValues(t, haps, []any{0, 1, 2, 3, "sd"})
}

// TestMiniRangePlainUnchanged is the regression guard for the range branch:
// a range with no surrounding suffixes must behave exactly as before.
func TestMiniRangePlainUnchanged(t *testing.T) {
	haps := Mini("0 .. 3").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	wantValues(t, haps, []any{0, 1, 2, 3})
}

func TestMiniPolymeterUsesFirstLayerStepCount(t *testing.T) {
	// Two steps per cycle, taken from the first layer's length.
	p := Mini("{bd sd, hh hh hh}")
	got := len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if got != 4 {
		t.Fatalf("got %d haps, want 4 (2 steps per cycle in each of 2 layers)", got)
	}
}

func TestMiniPolymeterExplicitStepCount(t *testing.T) {
	p := Mini("{bd sd, hh hh hh}%4")
	got := len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if got != 8 {
		t.Errorf("got %d haps, want 8 (4 steps per cycle in each of 2 layers)", got)
	}
}

// The point of a polymeter is that a layer whose length does not divide the
// step count lands on different elements from one cycle to the next.
func TestMiniPolymeterLayersDriftAcrossCycles(t *testing.T) {
	p := Mini("{bd sd, a b c}")
	seen := map[string]bool{}
	for c := int64(0); c < 3; c++ {
		for _, h := range p.QueryArc(core.FractionFromInt(c), core.FractionFromInt(c+1)) {
			if s, ok := h.Value.(string); ok {
				seen[s] = true
			}
		}
	}
	for _, want := range []string{"a", "b", "c"} {
		if !seen[want] {
			t.Errorf("the three-element layer never played %q across three cycles: %v", want, seen)
		}
	}
}

// TestMiniPolymeterStepReplicateExpands guards a step inside a "{...}"
// layer against the same suffix-leak class Tasks 2 and 3 fixed for plain
// sequence steps and "<...>" alternatives: "bd!2" is one token that must
// expand into two sibling steps in the layer's own list, not leak "bd!2" as
// a literal value. "{bd!2 sd, hh}" makes the first layer 3 steps long
// (bd, bd, sd), so with no explicit %n that becomes the steps-per-cycle
// rate and "hh" (a 1-step layer) repeats 3 times to match.
func TestMiniPolymeterStepReplicateExpands(t *testing.T) {
	haps := Mini("{bd!2 sd, hh}").QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	wantValues(t, haps, []any{"bd", "bd", "sd", "hh", "hh", "hh"})
}

// TestMiniPolymeterStepsAndFastCombine covers a leftover *n riding along
// after %n on the same "{...}" token. "%4" alone gives 4 steps per cycle in
// each of the 2 layers (8 haps); the trailing "*2" then doubles the event
// rate on top of that polymeter, for 16. This also guards against a
// digit-extraction regression: parsing "%4*2" as a whole integer fails, and
// silently discarding the unparsed "*2" instead of applying it would leave
// only 8 haps.
func TestMiniPolymeterStepsAndFastCombine(t *testing.T) {
	p := Mini("{bd sd, hh hh hh}%4*2")
	got := len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if got != 16 {
		t.Fatalf("got %d haps, want 16 (4 steps/cycle * 2 layers, doubled by *2)", got)
	}
}

// TestMiniPolymeterLayersDriftInOrder pins the exact per-cycle selection of
// the drifting three-element layer, not just that all three appear
// somewhere across three cycles: at 2 steps/cycle the layer's own clock
// only advances by 2 steps per outer cycle, so cycle 0 sees "a","b", cycle 1
// wraps to "c","a" and cycle 2 lands on "b","c". A test that only checked
// membership (as TestMiniPolymeterLayersDriftAcrossCycles does) would not
// catch a layer that cycled in the wrong order or at the wrong rate.
func TestMiniPolymeterLayersDriftInOrder(t *testing.T) {
	p := Mini("{bd sd, a b c}")
	want := [][]any{
		{"a", "b"},
		{"c", "a"},
		{"b", "c"},
	}
	for c := int64(0); c < 3; c++ {
		haps := p.QueryArc(core.FractionFromInt(c), core.FractionFromInt(c+1))
		var got []any
		for _, h := range haps {
			if s, ok := h.Value.(string); ok && s != "bd" && s != "sd" {
				got = append(got, s)
			}
		}
		if len(got) != len(want[c]) {
			t.Fatalf("cycle %d: got %v, want %v", c, got, want[c])
		}
		for i := range got {
			if got[i] != want[c][i] {
				t.Errorf("cycle %d: got %v, want %v", c, got, want[c])
			}
		}
	}
}
