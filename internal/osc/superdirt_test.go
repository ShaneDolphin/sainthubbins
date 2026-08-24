// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// pairs turns the flat argument list into a map so tests can assert on keys
// without depending on ordering.
func pairs(t *testing.T, args []any) map[string]any {
	t.Helper()
	if len(args)%2 != 0 {
		t.Fatalf("argument list must be key/value pairs, got %d items: %v", len(args), args)
	}
	m := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			t.Fatalf("key at %d is %T, want string", i, args[i])
		}
		m[k] = args[i+1]
	}
	return m
}

func hapOf(v any) core.Hap {
	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	return core.Hap{Whole: &span, Part: span, Value: v}
}

func TestDirtArgsFromControlBag(t *testing.T) {
	h := hapOf(map[string]any{"s": "bd", "gain": 0.8, "n": 2})
	m := pairs(t, DirtArgs(h, 0.5, 0.5))

	if m["s"] != "bd" {
		t.Errorf("s = %v, want bd", m["s"])
	}
	if m["gain"] != 0.8 {
		t.Errorf("gain = %v, want 0.8", m["gain"])
	}
	if m["n"] != 2 {
		t.Errorf("n = %v, want 2", m["n"])
	}
	if _, ok := m["cps"]; !ok {
		t.Error("every message must carry cps")
	}
	if m["delta"] != 0.5 {
		t.Errorf("delta = %v, want the event duration in seconds", m["delta"])
	}
}

func TestDirtArgsFromBareValues(t *testing.T) {
	// Raw mini-notation produces bare values rather than control bags.
	m := pairs(t, DirtArgs(hapOf("bd"), 0.5, 0.25))
	if m["s"] != "bd" {
		t.Errorf("a bare string should become the sound name, got %v", m["s"])
	}

	m = pairs(t, DirtArgs(hapOf(3), 0.5, 0.25))
	if m["n"] != 3 {
		t.Errorf("a bare number should become n, got %v", m["n"])
	}
}

func TestDirtArgsSkipsInternalKeys(t *testing.T) {
	h := hapOf(map[string]any{"s": "bd", "_cps": 0.5})
	m := pairs(t, DirtArgs(h, 0.5, 0.25))
	if _, ok := m["_cps"]; ok {
		t.Error("underscore-prefixed keys are engine internals and must not be sent")
	}
}
