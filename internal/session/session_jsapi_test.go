// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"testing"

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
func TestSessionEvaluateReportsJSError(t *testing.T) {
	s := NewSession()
	if _, err := s.Evaluate(`s("bd" +`); err == nil {
		t.Fatal("want an error for unparseable JS, got nil")
	}
}

func TestSessionEvaluateReportsBadMethod(t *testing.T) {
	s := NewSession()
	if _, err := s.Evaluate(`s("bd").nope()`); err == nil {
		t.Fatal("want an error for a nonexistent method, got nil")
	}
}
