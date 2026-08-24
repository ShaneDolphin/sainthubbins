// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestVlpfXFade(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"}).Vlpf(800)
	if len(p.FirstCycle()) != 1 {
		t.Fatalf("vlpf expected 1")
	}
	if v := p.FirstCycle()[0].Value.(map[string]any)["lpf"]; v != 800.0 {
		t.Fatalf("vlpf expected 800 got %v", v)
	}
	x := XFade2(Pure("a"), 0.8, Pure("b"))
	if len(x.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("xfade expected haps")
	}
}

func TestMorphWeave(t *testing.T) {
	m := MorphList([]int{1, 0, 1}, []int{1, 1, 0}, 0.3)
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("morph expected haps")
	}
	w := Pure("a").Weave2(2, Pure("b"), Pure("c"))
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("weave2 expected haps")
	}
}
