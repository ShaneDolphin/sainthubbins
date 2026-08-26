// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// TestSessionEvaluateAcceptsJS is the live `play` path — the one call site
// the task-4 brief's original two-file list missed. If this were still
// wired to the pre-jsapi fallback, `s("bd sd")` would silently become a
// one-hap literal-string pattern instead of two real control-bag haps.
func TestSessionEvaluateAcceptsJS(t *testing.T) {
	s := NewSession()
	pat, err := s.Evaluate(`s("bd sd")`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2", len(haps))
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Errorf("haps[0].Value = %#v, want a control bag carrying s:bd", haps[0].Value)
	}
}

// TestSessionEvaluateStillAcceptsMini is what the pre-jsapi fallback used to
// handle unconditionally, and must keep working.
func TestSessionEvaluateStillAcceptsMini(t *testing.T) {
	s := NewSession()
	pat, err := s.Evaluate("bd*4")
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("got %d haps, want 4", len(haps))
	}
}

// TestSessionEvaluateReportsJSError is the critical half: a genuine JS
// error (something mini-notation cannot rescue) must be returned to the
// caller, not swallowed. `play`'s caller (cmd/saint-hubbins/runPlay) relies
// on this error to abort before ever touching a socket.
//
// The property that actually matters on stage is stronger than "err !=
// nil" — asserting only that was named in this test's own original comment
// but never checked, the same shape of gap that has bitten this repo
// before. A live-coder who submits a typo mid-set needs the *previous*
// pattern to keep playing, not the room to go silent because of a mistake
// in what they just typed. So this also drives the trigger path (the same
// way session_output_test.go's TestSessionSendsHapsToItsSink does) and
// confirms the sink still receives the pattern from before the failed
// Evaluate call — s.Pattern (and what Cyclist is actually driving) must be
// untouched by a failed Evaluate, not reset to Silence.
func TestSessionEvaluateReportsJSError(t *testing.T) {
	s := NewSession()
	rec := &recorder{}
	s.SetSink(rec)

	if _, err := s.Evaluate("bd*4"); err != nil {
		t.Fatalf("Evaluate(good pattern): %v", err)
	}

	if _, err := s.Evaluate(`s("bd" +`); err == nil {
		t.Fatal("want an error for unparseable JS, got nil")
	}

	// The session's own live pattern must still be the one from before the
	// failed Evaluate call, not Silence.
	haps := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("s.Pattern produced %d haps after a failed Evaluate, want 4 (the previous pattern)", len(haps))
	}

	// And the Cyclist that's actually driving playback must be on the same
	// pattern.
	//
	// The OnTrigger call below does NOT catch a fix that updates s.Pattern
	// correctly but still calls s.Cyclist.SetPattern(core.Silence()):
	// OnTrigger is a plain callback field, invoked by the clock loop after
	// the Cyclist has already queried its own pattern, so calling it here
	// with a hap we pulled off s.Pattern never reads Cyclist.Pattern at all.
	// It exercises sink forwarding, not which pattern the Cyclist holds.
	// Assert that separately, or the Cyclist half of this test is a claim
	// with nothing behind it.
	if cp := s.Cyclist.Pattern; cp == nil ||
		len((*cp).QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) != 4 {
		t.Fatalf("Cyclist pattern was clobbered by a failed Evaluate — the room goes silent on a typo")
	}

	h := haps[0]
	s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))
	if rec.count() != 1 {
		t.Fatalf("sink received %d haps after a failed Evaluate, want 1 from the previous pattern", rec.count())
	}
}

func TestSessionEvaluateReportsBadMethod(t *testing.T) {
	s := NewSession()
	if _, err := s.Evaluate(`s("bd").nope()`); err == nil {
		t.Fatal("want an error for a nonexistent method, got nil")
	}
}
