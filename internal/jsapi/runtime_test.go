// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"testing"

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
	m, ok := vals[0].(map[string]any)
	if !ok {
		t.Fatalf("value is %T, want a control bag", vals[0])
	}
	if m["s"] != "bd" {
		t.Errorf("s = %v, want bd", m["s"])
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
