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

func TestDirtArgsFromBareInt64(t *testing.T) {
	// Bare int64 values should produce n parameter, avoiding silent data loss.
	m := pairs(t, DirtArgs(hapOf(int64(42)), 0.5, 0.25))
	if m["n"] != int64(42) {
		t.Errorf("a bare int64 should become n, got %v (type %T)", m["n"], m["n"])
	}
}

func TestDirtArgsSliceOrderIsDeterministic(t *testing.T) {
	// Bags with multiple keys must produce sorted output so the wire format
	// is stable and reproducible. Assert on the slice itself, not through pairs().
	h := hapOf(map[string]any{"z": 1, "a": 2, "m": 3})
	args := DirtArgs(h, 0.5, 0.25)

	// Expected order: sorted keys (a, m, z), then cps and delta.
	// args = []any{"a", 2, "m", 3, "z", 1, "cps", 0.5, "delta", 0.25}
	if len(args) != 10 {
		t.Fatalf("expected 10 elements (3 keys + cps/delta), got %d: %v", len(args), args)
	}

	if args[0] != "a" || args[1] != 2 {
		t.Errorf("first pair should be a, 2; got %v, %v", args[0], args[1])
	}
	if args[2] != "m" || args[3] != 3 {
		t.Errorf("second pair should be m, 3; got %v, %v", args[2], args[3])
	}
	if args[4] != "z" || args[5] != 1 {
		t.Errorf("third pair should be z, 1; got %v, %v", args[4], args[5])
	}
	if args[6] != "cps" || args[7] != 0.5 {
		t.Errorf("cps pair should be cps, 0.5; got %v, %v", args[6], args[7])
	}
	if args[8] != "delta" || args[9] != 0.25 {
		t.Errorf("delta pair should be delta, 0.25; got %v, %v", args[8], args[9])
	}
}

func TestDirtArgsBareNumericStringIsASoundName(t *testing.T) {
	// Mini.go stores raw numeric tokens as strings (e.g., "3" from pattern "0 1 2 3").
	// These arrive as bare string values, not Go numeric types, so they become s
	// (sound names), not n (note numbers). This is the documented contract: numeric
	// identity requires either a Go numeric value or a control such as core.N.
	m := pairs(t, DirtArgs(hapOf("3"), 0.5, 0.25))
	if m["s"] != "3" {
		t.Errorf("a bare numeric string should become sound name s, got %v", m["s"])
	}
	if _, hasN := m["n"]; hasN {
		t.Error("a bare numeric string should not produce n")
	}
}

func TestDirtArgsOutputsCanBeEncoded(t *testing.T) {
	// Every value type that DirtArgs emits must be encodable by EncodeMessage.
	// This is the contract between these layers: DirtArgs produces only types
	// the encoder can handle. This test would catch any regression where a new
	// type is emitted but the encoder cannot encode it (e.g., the uint64 case
	// before it was converted to int64).
	tests := []struct {
		name  string
		value any
	}{
		{"bare string", "bd"},
		{"bare int", 42},
		{"bare int64", int64(99)},
		{"bare float32", float32(3.14)},
		{"bare float64", 2.71},
		{"bare uint64", uint64(100)},
		{"control bag with multiple types", map[string]any{
			"s":    "kick",
			"gain": 0.8,
			"n":    5,
			"pan":  float32(0.5),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hapOf(tt.value)
			args := DirtArgs(h, 0.5, 0.25)

			// Attempt to encode all args produced by DirtArgs.
			_, err := EncodeMessage(DirtAddress, args...)
			if err != nil {
				t.Errorf("EncodeMessage failed for %s: %v", tt.name, err)
			}
		})
	}
}
