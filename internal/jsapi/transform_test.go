// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func countHaps(t *testing.T, code string) int {
	t.Helper()
	p, err := Evaluate(code)
	if err != nil {
		t.Fatalf("Evaluate(%q): %v", code, err)
	}
	return len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
}

func TestFastDoublesEvents(t *testing.T) {
	if got := countHaps(t, `s("bd sd").fast(2)`); got != 4 {
		t.Errorf("got %d haps, want 4", got)
	}
}

func TestChainedCalls(t *testing.T) {
	p, err := Evaluate(`s("bd").gain(0.5).cutoff(800)`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	m := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0].Value.(map[string]any)
	if m["s"] != "bd" || m["gain"] != 0.5 || m["cutoff"] != 800.0 {
		t.Errorf("got %v, want s:bd gain:0.5 cutoff:800", m)
	}
}

func TestUnaryTransform(t *testing.T) {
	p, err := Evaluate(`s("bd sd").rev()`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	first := haps[0].Value.(map[string]any)
	if first["s"] != "bd" || haps[0].Part.Begin.Float64() != 0.5 {
		t.Errorf("rev did not move bd to the second half: %v @ %v",
			first, haps[0].Part.Begin.Float64())
	}
}

func TestEuclidFromJS(t *testing.T) {
	if got := countHaps(t, `s("bd").euclid(3, 8)`); got != 3 {
		t.Errorf("got %d haps, want 3", got)
	}
}

// Error surfacing is the point of this whole plan: a bad chained call must
// become a Go error, never a silently-unchanged or silently-wrong pattern.
// Each case below is a way attachMethods could plausibly fall back to a
// no-op instead of failing loud.
func TestChainErrorSurfacing(t *testing.T) {
	cases := map[string]string{
		"fast with no argument":                     `s("bd").fast()`,
		"fast with a non-numeric argument":          `s("bd").fast("banana")`,
		"unrecognized method":                       `s("bd").nosuchmethod()`,
		"euclid with one argument":                  `s("bd").euclid(3)`,
		"control setter with no argument":           `s("bd").gain()`,
		"every with too few arguments":              `s("bd").every(3)`,
		"every with a non-function second argument": `s("bd").every(3, "notAFunction")`,
	}
	for name, code := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Evaluate(code); err == nil {
				t.Errorf("Evaluate(%q): want an error, got nil", code)
			}
		})
	}
}

// hapKey compares haps by (start, end, value) rather than struct identity —
// Whole/Context can legitimately differ in representation between a wide
// query and a stitched sequence of narrow ones without the pattern being
// wrong; start, end and value are the property that actually matters.
type hapKey struct {
	begin, end float64
	value      string
}

func hapKeys(haps []core.Hap) []hapKey {
	out := make([]hapKey, len(haps))
	for i, h := range haps {
		out[i] = hapKey{h.Part.Begin.Float64(), h.Part.End.Float64(), fmt.Sprintf("%v", h.Value)}
	}
	return out
}

// TestEveryWideQueryMatchesStitchedCycles is the property CLAUDE.md calls
// out explicitly: internal/audio/webaudio.go renders with exactly one wide
// QueryArc(0, cycles) call, so a cycle-dependent combinator that forgets
// SplitQueries is correct only by accident, when a test happens to query
// one cycle at a time. core.Pattern.Every already calls SplitQueries
// (pattern_time.go:198), but that guarantee is only worth anything if the
// jsapi wiring around it — re-entering the JS callback per split query —
// doesn't accidentally undo it. every(3, fn) only actually applies fn on
// one cycle in three, so a naive (non-splitting) implementation would show
// up here as a mismatch between the wide query and the stitched cycles,
// not as a crash.
func TestEveryWideQueryMatchesStitchedCycles(t *testing.T) {
	p, err := Evaluate(`s("bd").every(3, function(p) { return p.fast(2); })`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	const cycles = 6
	wide := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles))

	var stitched []core.Hap
	for i := int64(0); i < cycles; i++ {
		stitched = append(stitched, p.QueryArc(core.FractionFromInt(i), core.FractionFromInt(i+1))...)
	}

	wantWide, wantStitched := hapKeys(wide), hapKeys(stitched)
	if !reflect.DeepEqual(wantWide, wantStitched) {
		t.Fatalf("wide query over %d cycles disagrees with %d stitched one-cycle queries:\nwide:     %v\nstitched: %v",
			cycles, cycles, wantWide, wantStitched)
	}

	// Guard against a vacuous pass: every(3, fast(2)) must actually vary
	// hap count cycle to cycle (2 on the applied cycles, 1 elsewhere), or
	// this test can't tell "correctly split" from "never applied at all".
	counts := map[int]int{}
	for i := int64(0); i < cycles; i++ {
		n := len(p.QueryArc(core.FractionFromInt(i), core.FractionFromInt(i+1)))
		counts[n]++
	}
	if counts[1] == 0 || counts[2] == 0 {
		t.Fatalf("expected a mix of 1-hap and 2-hap cycles, got counts %v (haps per cycle)", counts)
	}
}

// TestChainDoesNotMutateReceiver confirms method chaining builds a new
// wrapped pattern rather than mutating jp.pat in place. This is a
// white-box test (same package) so it can hold the *jsPattern the JS
// object wraps and check it directly after invoking a method through JS —
// a mutation bug here would otherwise only show up as action-at-a-distance
// on an unrelated reference to the same pattern.
func TestChainDoesNotMutateReceiver(t *testing.T) {
	jp := &jsPattern{pat: mini.Mini("bd sd")}
	vm := goja.New()
	if err := register(vm); err != nil {
		t.Fatalf("register: %v", err)
	}
	obj := vm.ToValue(jp).(*goja.Object)
	attachMethods(vm, obj, jp)
	if err := vm.Set("p", obj); err != nil {
		t.Fatalf("vm.Set: %v", err)
	}

	before := hapKeys(jp.pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
	if _, err := vm.RunString(`p.fast(2)`); err != nil {
		t.Fatalf("RunString: %v", err)
	}
	after := hapKeys(jp.pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))

	if !reflect.DeepEqual(before, after) {
		t.Errorf("jp.pat was mutated by chaining .fast(2): before %v, after %v", before, after)
	}
}
