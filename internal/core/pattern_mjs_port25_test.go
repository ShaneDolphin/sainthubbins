// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twenty-fifth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_BinaryNCases(t *testing.T) {
	// binaryN like patterns already approximated via FastCat
	p := FastCat(Pure(0), Pure(1), Pure(0), Pure(1))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("binary 0101 4 expected 4")
	}
}

func TestMJS_PickCases(t *testing.T) {
	p := Pick("a", Pure(0))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("pick a 0 empty")
	}
	p2 := Pick(Pure("a"), Pure(0))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("pick pure a 0 empty")
	}
}

func TestMJS_TimeCatPolymeterVariations(t *testing.T) {
	tc := TimeCat(FractionFromInt(2), Pure("a"), FractionFromInt(1), Pure("b"))
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("timeCat 2:a 1:b empty")
	}
	pm := Polymeter(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter 4 empty")
	}
	// stepCat variation
	sc := StepCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stepCat empty")
	}
}
