// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"sync"
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func hapValues(t *testing.T, p core.Pattern) []any {
	t.Helper()
	var out []any
	for _, h := range p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)) {
		out = append(out, h.Value)
	}
	return out
}

func TestEvaluateSoundConstructor(t *testing.T) {
	p, err := Evaluate(`s("bd sd")`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	vals := hapValues(t, p)
	if len(vals) != 2 {
		t.Fatalf("got %d haps, want 2", len(vals))
	}
	m0, ok := vals[0].(map[string]any)
	if !ok {
		t.Fatalf("vals[0] is %T, want a control bag", vals[0])
	}
	if m0["s"] != "bd" {
		t.Errorf("vals[0][\"s\"] = %v, want bd", m0["s"])
	}
	// The claim is about both haps of the sequence, not just the first — a
	// pattern that duplicated "bd" instead of sequencing in "sd" must fail
	// here, not slip through on a presence-only check.
	m1, ok := vals[1].(map[string]any)
	if !ok {
		t.Fatalf("vals[1] is %T, want a control bag", vals[1])
	}
	if m1["s"] != "sd" {
		t.Errorf("vals[1][\"s\"] = %v, want sd", m1["s"])
	}
}

// A syntax error must be reported, not silently turned into a literal hap.
func TestEvaluateReportsErrors(t *testing.T) {
	if _, err := Evaluate(`s("bd" +`); err == nil {
		t.Fatal("want a syntax error, got nil")
	}
	if _, err := Evaluate(`notAFunction("x")`); err == nil {
		t.Fatal("want a reference error, got nil")
	}
}

func TestEvaluateRejectsNonPatternResult(t *testing.T) {
	if _, err := Evaluate(`42`); err == nil {
		t.Fatal("want an error when the result is not a pattern")
	}
}

// Every Evaluate call registers the mini-notation string-parser hook
// (core.SetStringParser), a package-level var in internal/core shared by
// every caller. Run enough concurrent evaluations, with a mix of code paths
// that both write and read that hook, to give the race detector a real
// chance to see an unsynchronized access if one exists.
func TestEvaluateConcurrent(t *testing.T) {
	const n = 20
	var wg sync.WaitGroup
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			code := `s("bd sd")`
			if i%2 == 0 {
				// Exercises Reify's read of the string-parser hook directly,
				// not just through the s() constructor.
				code = `"bd sd"`
			}
			_, err := Evaluate(code)
			errs[i] = err
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: Evaluate: %v", i, err)
		}
	}
}

// A script that never returns must not hang the calling goroutine forever —
// it must be interrupted and reported as an error.
func TestEvaluateInterruptsRunawayScript(t *testing.T) {
	start := time.Now()
	_, err := Evaluate(`while (true) {}`)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("want an error from a runaway script, got nil")
	}
	// A near-instant error would mean something other than the interrupt
	// path fired (e.g. a parse error) — confirm we actually waited out
	// something close to the timeout rather than failing for a different
	// reason entirely.
	if elapsed < evaluateTimeout-500*time.Millisecond {
		t.Fatalf("Evaluate returned after %s, want it to wait close to the %s timeout", elapsed, evaluateTimeout)
	}
	if elapsed > evaluateTimeout+5*time.Second {
		t.Fatalf("Evaluate took %s, want it to return promptly after the %s timeout", elapsed, evaluateTimeout)
	}
	// The interrupt value is already unprefixed and the error branch adds
	// "jsapi:" once. Assert the whole string so a doubled prefix — which is
	// what this looked like before — fails here rather than reaching a user.
	want := "jsapi: evaluation exceeded " + evaluateTimeout.String()
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}
