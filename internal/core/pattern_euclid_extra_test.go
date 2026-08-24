package core

import "testing"

func TestEuclidVariants(t *testing.T) {
	p := Pure("x").Euclid(2, 5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Euclid 2,5 expected 2 got %d", len(haps))
	}
	q := Pure("x").Euclid(3, 8)
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 3 {
		t.Fatalf("Euclid 3,8 expected 3 got %d", len(haps2))
	}
}

func TestEuclidRot(t *testing.T) {
	p := Pure("x").EuclidRot(3, 8, 2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("EuclidRot expected 3 got %d", len(haps))
	}
}

func TestBjorklundAlias(t *testing.T) {
	// Bjorklund is alias for Euclid via mini term
	p := Pure("bd").Euclid(3, 8)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Bjorklund alias expected 3 got %d", len(haps))
	}
}

func TestStructBasic(t *testing.T) {
	// Struct: boolean pattern determines where original is played
	boolPat := FastCat(Pure(true), Pure(false), Pure(true))
	p := Pure("a").Struct(boolPat)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Struct expected >=1 got 0")
	}
}

func TestStepCatBasic(t *testing.T) {
	p := StepCat(Pure("a"), Pure("b"), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("StepCat expected 3 got %d", len(haps))
	}
}

func TestTimeCatBasic(t *testing.T) {
	p := TimeCat(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("TimeCat expected >=1")
	}
}
