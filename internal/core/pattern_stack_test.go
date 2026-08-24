package core

import "testing"

func TestStackWithSilence(t *testing.T) {
	p := Stack(Pure("a"), Silence())
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Stack with Silence expected 1 got %d", len(haps))
	}
}

func TestFastCatEmpty(t *testing.T) {
	p := FastCat()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("FastCat empty expected 0 got %d", len(haps))
	}
}

func TestSlowCatWithSilence(t *testing.T) {
	p := SlowCat(Pure("a"), Silence(), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps) < 2 {
		t.Fatalf("SlowCat with Silence expected >=2 over 3 cycles got %d", len(haps))
	}
}

func TestSequenceAlias(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Sequence alias expected 2 got %d", len(haps))
	}
}

func TestCatAlias(t *testing.T) {
	p := Cat(Pure("a"), Pure("b"), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps) != 3 {
		t.Fatalf("Cat alias expected 3 over 3 cycles got %d", len(haps))
	}
}
