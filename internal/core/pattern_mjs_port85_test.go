package core

import "testing"

func TestMJS_BinaryNCases2(t *testing.T) {
	p := FastCat(Pure(0), Pure(1), Pure(0), Pure(1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("FastCat 0101 expected 4 got %d", len(haps))
	}
}

func TestMJS_PickCases2(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b", "c"})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Choose a")
	}
	p2 := Pure("a").Choose([]any{"x", "y"})
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Choose x,y")
	}
}

func TestMJS_TimeCatPolymeterVariations2(t *testing.T) {
	tc := TimeCatWeighted(2, Pure("a"), 1, Pure("b"))
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("TimeCatWeighted 2,1")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"), Pure("c"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat 3")
	}
	sc := StepCat(Pure("a"), Pure("b"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("StepCat 2")
	}
}
