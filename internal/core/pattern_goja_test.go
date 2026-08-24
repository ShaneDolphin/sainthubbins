// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Goja bridge smoke test for JS<->Go parity (Phase 1.1 testing strategy).
// This is a temporary shim: goja is only used in tests/transpiler, not production audio.
//go:build goja

package core

import (
	"testing"

	"github.com/dop251/goja"
)

func TestGojaSmoke(t *testing.T) {
	vm := goja.New()
	v, err := vm.RunString(`1+2+3`)
	if err != nil {
		t.Fatalf("goja run: %v", err)
	}
	if v.ToInteger() != 6 {
		t.Fatalf("goja expected 6 got %v", v)
	}
}

func TestGojaPatternParityPlaceholder(t *testing.T) {
	// Placeholder for full pattern.mjs via goja comparison (1394 cases).
	// Phase 1.1 requires: run JS pattern.mjs via goja to generate expected []Hap
	// and compare with Go Query. Full ESM import not yet wired (needs bundling).
	// This test verifies the bridge can evaluate a pure-JS Fraction mock.
	vm := goja.New()
	_, err := vm.RunString(`
		var Fraction = function(n,d){ this.n=n; this.d=d||1; };
		Fraction.prototype.add = function(o){ return new Fraction(this.n*o.d + o.n*this.d, this.d*o.d); };
		var a = new Fraction(1,2);
		var b = new Fraction(1,3);
		var c = a.add(b);
		if (c.n !== 5 || c.d !== 6) throw new Error("fraction add failed");
	`)
	if err != nil {
		t.Fatalf("goja fraction mock: %v", err)
	}
}
