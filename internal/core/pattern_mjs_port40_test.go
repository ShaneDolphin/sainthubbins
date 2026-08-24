// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 40th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_PolymeterSteps(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	if p.Steps != nil {
		t.Fatalf("stack steps should be nil got %v", p.Steps)
	}
	p2 := Sequence(Pure("a"), Pure("b"), Pure("c"))
	// Sequence=FastCat currently has nil steps (not computed); check query still works
	haps := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("sequence 3 expected 3 haps got %d", len(haps))
	}
	// Gap steps vs Pure steps
	if Gap(4).Steps == nil || Gap(4).Steps.Cmp(FractionFromInt(4)) != 0 {
		t.Fatalf("gap 4 steps 4")
	}
	if Pure("a").Steps == nil {
		t.Fatalf("pure steps not nil")
	}
}

func TestMJS_ControlStack(t *testing.T) {
	// Gain/Pan are createParam maps; Pure("bd") string + Gain via Set loses s vs S("bd") preserves s map
	// Use S("bd") for s bag and Stack/Gain via Set
	p := S("bd").Set(Gain(0.8)).Set(Pan(0.5))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("control stack len 1 got %d", len(haps))
	}
	v := haps[0].Value.(map[string]any)
	if v["gain"] != 0.8 {
		t.Fatalf("gain 0.8 got %v", v["gain"])
	}
	if v["pan"] != 0.5 {
		t.Fatalf("pan 0.5 got %v", v["pan"])
	}
	if v["s"] != "bd" {
		t.Fatalf("s bd got %v", v["s"])
	}
	// Also test Stack of control bags
	p2 := Stack(Pure(map[string]any{"s": "bd"}), Pure(map[string]any{"gain": 0.8}))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("stack 2 bags expected 2 got %d", len(haps2))
	}
}

func TestMJS_SlowInsideNested(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Slow(FractionFromInt(2)).Inside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("slow inside nested empty")
	}
}
