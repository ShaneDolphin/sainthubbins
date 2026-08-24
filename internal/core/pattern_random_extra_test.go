package core

import "testing"

func TestRandBasic(t *testing.T) {
	p := Rand()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Rand expected >=1")
	}
	// Rand should be in [0,1)
	for _, h := range haps {
		if v, ok := h.Value.(float64); ok {
			if v < 0 || v >= 1 {
				t.Fatalf("Rand value out of range %v", v)
			}
		}
	}
}

func TestChooseBasic(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b", "c"})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Choose expected >=1")
	}
}

func TestDegradeBasic(t *testing.T) {
	p := Pure("a").Degrade()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Degrade may remove some haps but not panic
	if len(haps) > 1 {
		t.Fatalf("Degrade single expected 0 or 1 got %d", len(haps))
	}
}

func TestSometimesBasic(t *testing.T) {
	p := Pure("a").Sometimes(func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Sometimes expected 1 got %d", len(haps))
	}
}

func TestSometimesByBasic(t *testing.T) {
	p := Pure("a").SometimesBy(0.5, func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("SometimesBy expected 1")
	}
}
