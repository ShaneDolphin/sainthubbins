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

// TestEveryCallbackInfiniteLoopIsInterrupted covers a hang Evaluate's own
// timer cannot see: an every()/off()-style callback runs at query time,
// long after Evaluate has returned and its timer has already stopped, so
// an infinite loop written inside the callback itself needs its own
// interrupt armed independently. Before that, this exact input ran past
// 15s in a real process (measured with a bounded runner) while a bare
// top-level `while(true){}` already returned by 6s — confirming the two
// hangs were not caught by the same mechanism.
func TestEveryCallbackInfiniteLoopIsInterrupted(t *testing.T) {
	// The callback hangs only on its first invocation (cycle 0) and behaves
	// normally on its second (cycle 2 — every(2, ...) matches Mod(cycle,2)
	// == 0). This is what lets the test distinguish "the stale interrupt
	// from cycle 0's timeout bled into cycle 2" from "it didn't": without
	// vm.ClearInterrupt() (armInterruptTimeout, runtime.go), goja queues an
	// Interrupt() called while not running and fires it on the very next
	// JS call into this vm — so cycle 2's otherwise-normal invocation would
	// itself immediately report an interrupted error and fall back to
	// Silence, indistinguishable from cycle 0's real timeout except by
	// timing. Checking cycle 2's haps, not just its speed, is what actually
	// proves the interrupt was cleared rather than never having reached it.
	pat, err := Evaluate(`
		let n = 0;
		s("bd*4").every(2, x => {
			n++;
			if (n === 1) { while (true) {} }
			return x;
		})
	`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	start := time.Now()
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	elapsed := time.Since(start)

	if elapsed < evaluateTimeout-500*time.Millisecond {
		t.Fatalf("query returned after %s, want it to wait close to the %s per-callback timeout", elapsed, evaluateTimeout)
	}
	if elapsed > evaluateTimeout+5*time.Second {
		t.Fatalf("query took %s, want it to return promptly after the %s per-callback timeout", elapsed, evaluateTimeout)
	}
	// The interrupted callback falls back to Silence — the same fallback a
	// callback error already uses (registry.go) — so this cycle simply
	// produces no haps rather than hanging or crashing the process.
	if len(haps) != 0 {
		t.Errorf("got %d haps from an interrupted callback cycle, want 0 (Silence fallback)", len(haps))
	}

	// Cycle 2 re-enters the SAME callback on the SAME vm, this time
	// returning normally and promptly. A stale, uncleared interrupt would
	// make this cycle fail exactly like cycle 0 did — near-instantly,
	// rather than hanging, but still zero haps.
	start = time.Now()
	haps = pat.QueryArc(core.FractionFromInt(2), core.FractionFromInt(3))
	elapsed = time.Since(start)
	if elapsed > time.Second {
		t.Fatalf("cycle 2 query took %s, want near-instant — its callback does not loop", elapsed)
	}
	if len(haps) != 4 {
		t.Fatalf("got %d haps on cycle 2, want 4 — a stale interrupt from cycle 0's timeout would zero this out too", len(haps))
	}
}
