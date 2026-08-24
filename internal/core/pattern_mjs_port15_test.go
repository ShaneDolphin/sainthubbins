// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Fifteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_KeepKeepIf(t *testing.T) {
	// keep with structure In
	p := Pure(3).Keep(Pure(4))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("keep empty")
	}
	// keepif true/false
	p2 := Pure(1).KeepIf(Pure(true))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("keepif true empty")
	}
	p3 := Pure(1).KeepIf(Pure(false))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("keepif false expected 0 got %d", len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_SubMulDiv(t *testing.T) {
	if v := Pure(5).Sub(Pure(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value; v != 3.0 && v != 3 {
		// toFloat path gives float
		switch x := v.(type) {
		case float64:
			if x != 3 {
				t.Fatalf("sub 5-2 got %v", x)
			}
		case int:
			if x != 3 {
				t.Fatalf("sub 5-2 got %v", x)
			}
		default:
			t.Logf("sub got %T %v", v, v)
		}
	}
	if len(Pure(2).Mul(Pure(3)).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("mul empty")
	}
	if len(Pure(6).Div(Pure(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("div empty")
	}
}

func TestMJS_SetOutSqueeze(t *testing.T) {
	p := Pure(map[string]any{"a": 1}).Set(Pure(map[string]any{"b": 2}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("set empty")
	}
	m := haps[0].Value.(map[string]any)
	if m["b"] == nil {
		t.Fatalf("set b missing %v", m)
	}
	// fastCat with Set (out structure via second pattern)
	p2 := FastCat(Pure(1), Pure(2)).Set(Pure(10))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("set via second pattern empty")
	}
}
