// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEvalCodeEvaluatesJS(t *testing.T) {
	out, n, err := evalCode(`s("bd sd")`)
	if err != nil {
		t.Fatalf("evalCode: %v", err)
	}
	if n != 2 {
		t.Fatalf("got %d haps, want 2", n)
	}
	if !strings.Contains(out, `"bd"`) || !strings.Contains(out, `"sd"`) {
		t.Errorf("output missing bd/sd control values: %s", out)
	}
}

func TestEvalCodeFallsBackToMini(t *testing.T) {
	out, n, err := evalCode("bd sd")
	if err != nil {
		t.Fatalf("evalCode: %v", err)
	}
	if n != 2 {
		t.Fatalf("got %d haps, want 2 — mini-notation must keep working", n)
	}
	if !strings.Contains(out, "bd") || !strings.Contains(out, "sd") {
		t.Errorf("output missing bd/sd: %s", out)
	}
}

// TestEvalCodeReportsJSError is the critical half: unparseable JS that
// mini-notation cannot rescue must be reported as an error, not silently
// rendered as a literal-string hap.
func TestEvalCodeReportsJSError(t *testing.T) {
	_, _, err := evalCode(`s("bd" +`)
	if err == nil {
		t.Fatal("want an error for unparseable JS, got nil")
	}
}

func TestEvalCodeReportsBadMethod(t *testing.T) {
	_, _, err := evalCode(`s("bd").nope()`)
	if err == nil {
		t.Fatal("want an error for a nonexistent method, got nil")
	}
}

func TestRenderPatternEvaluatesJS(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.wav")
	n, err := renderPattern(out, `s("bd sd")`)
	if err != nil {
		t.Fatalf("renderPattern: %v", err)
	}
	if n == 0 {
		t.Fatal("renderPattern wrote zero samples")
	}
}

func TestRenderPatternReportsJSError(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.wav")
	if _, err := renderPattern(out, `s("bd" +`); err == nil {
		t.Fatal("want an error for unparseable JS, got nil")
	}
}

func TestRunMIDIEvaluatesJS(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.mid")
	if err := runMIDI(`s("bd sd")`, out, 1); err != nil {
		t.Fatalf("runMIDI: %v", err)
	}
}

func TestRunMIDIReportsJSError(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.mid")
	if err := runMIDI(`s("bd" +`, out, 1); err == nil {
		t.Fatal("want an error for unparseable JS, got nil")
	}
}
