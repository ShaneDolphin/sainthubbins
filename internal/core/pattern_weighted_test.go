// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestStepCat(t *testing.T) {
	p1 := Pure("a")
	p2 := Pure("b")
	pat := StepCat([]any{FractionFromInt(3), p1}, []any{FractionFromInt(1), p2})
	if pat.Steps == nil || !pat.Steps.Equals(FractionFromInt(4)) {
		t.Fatalf("StepCat steps !=4 got %v", pat.Steps)
	}
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	haps := pat.Query(NewState(span, nil))
	if len(haps) == 0 {
		t.Fatalf("StepCat no haps")
	}
}

func TestPolymeterLCM(t *testing.T) {
	// two patterns with 2 and 3 steps -> lcm 6
	p1 := Pure("a").SetSteps(FractionFromInt(2))
	p2 := Pure("b").SetSteps(FractionFromInt(3))
	pm := Polymeter(p1, p2)
	if pm.Steps == nil || !pm.Steps.Equals(FractionFromInt(6)) {
		t.Fatalf("Polymeter lcm !=6 got %v", pm.Steps)
	}
}
