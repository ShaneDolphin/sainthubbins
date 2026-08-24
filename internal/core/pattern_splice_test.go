package core

import "testing"

func TestSpliceBasic(t *testing.T) {
	p := Splice(Pure("a"), Pure(1), Pure(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Splice expected >=1")
	}
}

func TestSliceBasic(t *testing.T) {
	p := Slice(2, Pure("a"), Pure(0))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Slice expected >=1")
	}
}

func TestChopBasic(t *testing.T) {
	p := Pure("a").Chop(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Chop expected >=1")
	}
}
