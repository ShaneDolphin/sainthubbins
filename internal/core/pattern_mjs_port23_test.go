// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-third batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ControlsValue(t *testing.T) {
	p := S("bd")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("s bd empty")
	}
	// value test: pure with controls
	p2 := Pure(map[string]any{"s": "bd", "n": 2})
	haps := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("map pure empty")
	}
	m := haps[0].Value.(map[string]any)
	if m["n"] != 2 {
		t.Fatalf("n 2 expected got %v", m["n"])
	}
}

func TestMJS_UtilsMod(t *testing.T) {
	if Mod(5, 3) != 2 {
		t.Fatalf("mod 5,3 expected 2 got %d", Mod(5, 3))
	}
	if Mod(-1, 3) != 2 {
		t.Fatalf("mod -1,3 expected 2 got %d", Mod(-1, 3))
	}
}

func TestMJS_LogValues(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v.(string) + "b" })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value != "ab" {
		t.Fatalf("fmap ab got %v", haps[0].Value)
	}
	// stack log-like
	s := Stack(Pure("a"), Pure("b")).Fmap(func(v any) any { return v })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("stack fmap 2 expected 2")
	}
}
