// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"strings"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestEvaluateCodeTriesJSFirst(t *testing.T) {
	pat, err := EvaluateCode(`s("bd sd")`)
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
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

// TestEvaluateCodeFallsBackToMini is the case the plan must not break:
// "bd sd" is two bare identifiers, not valid JS, but is exactly the
// mini-notation every existing template and eval gate depends on.
func TestEvaluateCodeFallsBackToMini(t *testing.T) {
	pat, err := EvaluateCode("bd sd")
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2 — mini-notation fallback must still work", len(haps))
	}
	if haps[0].Value != "bd" || haps[1].Value != "sd" {
		t.Errorf("values = %v, %v, want bd, sd", haps[0].Value, haps[1].Value)
	}
}

// TestEvaluateCodeReportsJSErrorWhenMiniAlsoFails is the critical half: a
// broken JS expression that mini-notation cannot rescue must surface the
// JS error, not silence and not a literal-string hap.
func TestEvaluateCodeReportsJSErrorWhenMiniAlsoFails(t *testing.T) {
	_, err := EvaluateCode(`s("bd" +`)
	if err == nil {
		t.Fatal("want an error for unparseable JS that mini-notation cannot rescue, got nil")
	}
}

func TestEvaluateCodeReportsErrorForBadMethod(t *testing.T) {
	_, err := EvaluateCode(`s("bd").nope()`)
	if err == nil {
		t.Fatal("want an error for a nonexistent method, got nil")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error %q should mention the offending method", err.Error())
	}
}
