// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestNextSteps_PolyJoinResetJoinDefragment(t *testing.T) {
	p := Pure(Pure("a")).PolyJoin()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolyJoin expected haps")
	}
	p2 := Pure(Pure("a")).ResetJoin()
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ResetJoin expected haps")
	}
	p3 := Pure("a").Defragment()
	// Defragment merges adjacent same-value haps — for single hap, should remain 1
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Defragment expected haps")
	}
}

func TestNextSteps_TimeCatPolymeter(t *testing.T) {
	p := TimeCat(FractionFromInt(1), Pure("a"), FractionFromInt(1), Pure("b"), FractionFromInt(2), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("TimeCat 1,1,2 expected haps")
	}
	pm := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Polymeter expected haps")
	}
}

func TestNextSteps_PickChooseRandom(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b", "c", "d"})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose expected haps")
	}
	p2 := Pure("a").DegradeBy(0)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0 expected haps")
	}
	p3 := Pure("a").SometimesBy(1.0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 1.0 expected haps")
	}
}

func TestNextSteps_OpMatrix(t *testing.T) {
	a := Pure(2.0)
	if v := a.Add(Pure(3.0)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v != 5 {
		t.Fatalf("Add 2+3 expected 5 got %v", v)
	}
	if v := a.Mul(Pure(4.0)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v != 8 {
		t.Fatalf("Mul 2*4 expected 8 got %v", v)
	}
	if v := a.Sub(Pure(1.0)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v != 1 {
		t.Fatalf("Sub 2-1 expected 1 got %v", v)
	}
	if v := a.Div(Pure(2.0)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v != 1 {
		t.Fatalf("Div 2/2 expected 1 got %v", v)
	}
	if v := a.Mod(Pure(1.5)).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64); v < 0.49 || v > 0.51 {
		t.Fatalf("Mod 2 mod 1.5 expected 0.5 got %v", v)
	}
}

func TestNextSteps_RegisterInterop(t *testing.T) {
	// SessionScope register interop: Evaluate via scope
	RegisterScope("testPat", Pure("hello"))
	p, _, err := Evaluate("testPat", nil)
	if err != nil {
		t.Fatalf("Evaluate testPat err %v", err)
	}
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value != "hello" {
		t.Fatalf("Evaluate testPat expected hello got %v", haps)
	}
	ClearScope()
}
