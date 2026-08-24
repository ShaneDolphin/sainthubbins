// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Tenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_PickModSqueeze(t *testing.T) {
	// pick via free func already tested; test pickMod
	p := PickMod("a", Pure(1))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("pickMod empty")
	}
	p2 := Pure(Pure("a")).SqueezeJoin()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("squeezeJoin empty")
	}
	p3 := Pure(Pure("a")).InnerJoin()
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("innerJoin empty")
	}
}

func TestMJS_CatSequence(t *testing.T) {
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("cat 3/3 expected 3 got %d", len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))))
	}
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("sequence 2 expected 2 got %d", len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}

func TestMJS_SetWithValue(t *testing.T) {
	p := Pure(map[string]any{"a": 1}).Set(Pure(map[string]any{"b": 2}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("set empty")
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok || m["b"] == nil {
		t.Fatalf("set b missing %v", haps[0].Value)
	}
	p2 := Pure(1).WithValue(func(v any) any { return v.(int) * 2 })
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("withValue empty")
	}
}
