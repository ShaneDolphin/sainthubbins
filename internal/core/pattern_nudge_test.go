package core

import "testing"

func TestEarlyBasic2(t *testing.T) {
	p := Pure("a").Early(0.25)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Early expected >=1")
	}
}

func TestMaskBasic(t *testing.T) {
	maskPat := FastCat(Pure(true), Pure(false), Pure(true))
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Mask(maskPat)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Mask expected >=1")
	}
}

func TestBypassBasic(t *testing.T) {
	p := Pure("a").Bypass(true)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("Bypass true expected 0 got %d", len(haps))
	}
	q := Pure("a").Bypass(false)
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 {
		t.Fatalf("Bypass false expected 1 got %d", len(haps2))
	}
}
