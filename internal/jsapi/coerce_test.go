// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
//
// Proof that the four JS->domain conversion paths this file used to
// implement independently — register()'s control constructors,
// attachMethods()'s control setters, patternFromJSValue, and unwrap — now
// agree, because they all route through coerceJSValue (via toPattern,
// toPatternResult and toControlValue). Before this consolidation, the
// same null that raised a clean TypeError inside stack(...) was silently
// accepted by s(...), and s(function(){}) embedded a raw Go pointer
// address as a control value.

package jsapi

import (
	"strings"
	"testing"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// TestControlConstructorRejectsInvalidArguments is the register()-side half
// of the reproduction table from review: s(null), s({}), s(function(){})
// and s(true) all used to succeed silently (or, for the function case,
// embed a raw Go pointer address in the control bag). Each must now raise
// a TypeError naming both the control and a JS-recognizable description of
// the offending value.
func TestControlConstructorRejectsInvalidArguments(t *testing.T) {
	cases := map[string]struct {
		code string
		want string // substring the error must contain
	}{
		"null":         {`s(null)`, "s: argument must be a pattern, string or number, got null"},
		"undefined":    {`s(undefined)`, "s: argument must be a pattern, string or number, got undefined"},
		"plain object": {`s({})`, "s: argument must be a pattern, string or number, got a plain object"},
		"function":     {`s(function(){})`, "s: argument must be a pattern, string or number, got a function"},
		"boolean":      {`s(true)`, "s: argument must be a pattern, string or number, got the boolean true"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Evaluate(c.code)
			if err == nil {
				t.Fatalf("Evaluate(%q): want an error, got nil", c.code)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("Evaluate(%q) error = %q, want it to contain %q", c.code, err.Error(), c.want)
			}
		})
	}
}

// TestControlSetterRejectsInvalidArguments is the attachMethods()-side half
// of the same table: .gain(null), .gain({}), .gain(function(){}) and
// .gain(true) used to be silently accepted the same way.
func TestControlSetterRejectsInvalidArguments(t *testing.T) {
	cases := map[string]struct {
		code string
		want string
	}{
		"null":         {`s("bd").gain(null)`, "gain: argument must be a pattern, string or number, got null"},
		"undefined":    {`s("bd").gain(undefined)`, "gain: argument must be a pattern, string or number, got undefined"},
		"plain object": {`s("bd").gain({})`, "gain: argument must be a pattern, string or number, got a plain object"},
		"function":     {`s("bd").gain(function(){})`, "gain: argument must be a pattern, string or number, got a function"},
		"boolean":      {`s("bd").gain(true)`, "gain: argument must be a pattern, string or number, got the boolean true"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Evaluate(c.code)
			if err == nil {
				t.Fatalf("Evaluate(%q): want an error, got nil", c.code)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("Evaluate(%q) error = %q, want it to contain %q", c.code, err.Error(), c.want)
			}
		})
	}
}

// TestControlConstructorAcceptsMiniNotationAndPattern is half of the
// must-keep-working list: s("bd sd") (bare string, mini-notation) and
// s(<a pattern>) (a wrapped pattern reaching createParam's Pattern branch,
// not its raw-value branch) must still work, with actual hap contents
// checked rather than just "no error".
func TestControlConstructorAcceptsMiniNotationAndPattern(t *testing.T) {
	p, err := Evaluate(`s("bd sd")`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2", len(haps))
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["s"] != "bd" {
		t.Errorf("haps[0] = %v, want s:bd", haps[0].Value)
	}
	if m, ok := haps[1].Value.(map[string]any); !ok || m["s"] != "sd" {
		t.Errorf("haps[1] = %v, want s:sd", haps[1].Value)
	}

	// s(mini("bd sd")) passes a *jsPattern (a raw two-hap pattern, not yet
	// control-wrapped) into s()'s constructor — this only produces a
	// correct two-hap control pattern if the argument reaches createParam's
	// Pattern branch (Fmap), not its raw-value branch (which would instead
	// wrap the whole Pattern object as a single opaque value).
	p2, err := Evaluate(`s(mini("bd sd"))`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps2 := p2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("s(mini(\"bd sd\")): got %d haps, want 2", len(haps2))
	}
	if m, ok := haps2[0].Value.(map[string]any); !ok || m["s"] != "bd" {
		t.Errorf("haps2[0] = %v, want s:bd", haps2[0].Value)
	}
	if m, ok := haps2[1].Value.(map[string]any); !ok || m["s"] != "sd" {
		t.Errorf("haps2[1] = %v, want s:sd", haps2[1].Value)
	}
}

// TestControlSetterAcceptsModulatedPattern is the .gain(sine)-shaped half
// of the must-keep-working list. This jsapi doesn't expose an oscillator
// global (sine) yet, so this white-box test builds the same shape by hand:
// a raw numeric Pattern (standing in for what sine would produce) injected
// as a JS global, then used as a control-setter argument. This only
// produces the expected per-hap gain value if the argument reaches
// createParam's Pattern branch (Fmap) rather than being rejected or
// embedded as an opaque value.
func TestControlSetterAcceptsModulatedPattern(t *testing.T) {
	vm := goja.New()
	if err := register(vm); err != nil {
		t.Fatalf("register: %v", err)
	}
	modSource := newJSPattern(vm, core.Pure(0.75))
	if err := vm.Set("modSource", modSource); err != nil {
		t.Fatalf("vm.Set: %v", err)
	}
	v, err := vm.RunString(`s("bd").gain(modSource)`)
	if err != nil {
		t.Fatalf("RunString: %v", err)
	}
	p, err := unwrap(vm, v)
	if err != nil {
		t.Fatalf("unwrap: %v", err)
	}
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("got %d haps, want 1", len(haps))
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok {
		t.Fatalf("hap value is %T, want a control bag", haps[0].Value)
	}
	if m["s"] != "bd" {
		t.Errorf("s = %v, want bd", m["s"])
	}
	if m["gain"] != 0.75 {
		t.Errorf("gain = %v, want 0.75 (modulated via a Pattern argument)", m["gain"])
	}
}

// TestControlAcceptsZeroAndNegativeNumbers is the numbers half of the
// must-keep-working list: controls without a domain restriction (unlike
// fast/slow/segment/ply) must keep accepting zero and negative numbers —
// gain(0) and gain(-1) are unusual but coherent (silence, or an inverted
// signal downstream), not caller mistakes.
func TestControlAcceptsZeroAndNegativeNumbers(t *testing.T) {
	cases := map[string]struct {
		code string
		want float64
	}{
		"zero":     {`s("bd").gain(0)`, 0},
		"negative": {`s("bd").gain(-1)`, -1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p, err := Evaluate(c.code)
			if err != nil {
				t.Fatalf("Evaluate(%q): %v", c.code, err)
			}
			haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
			if len(haps) != 1 {
				t.Fatalf("got %d haps, want 1", len(haps))
			}
			m := haps[0].Value.(map[string]any)
			if m["gain"] != c.want {
				t.Errorf("gain = %v, want %v", m["gain"], c.want)
			}
		})
	}
}

// TestIntegerControlValueNormalizedToFloat64 locks in normalizeNumber's
// existing behavior (goja exports a whole-number JS literal as int64) still
// applies after routing through toControlValue: `.cutoff(800)` must produce
// a float64(800), not an int64(800), so a control bag's numeric values stay
// consistently typed regardless of whether the JS literal had a fractional
// part.
func TestIntegerControlValueNormalizedToFloat64(t *testing.T) {
	p, err := Evaluate(`s("bd").cutoff(800)`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	m := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0].Value.(map[string]any)
	v, ok := m["cutoff"].(float64)
	if !ok {
		t.Fatalf("cutoff is %T, want float64", m["cutoff"])
	}
	if v != 800 {
		t.Errorf("cutoff = %v, want 800", v)
	}
}
