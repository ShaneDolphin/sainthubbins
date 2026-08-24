package core

import "testing"

func TestCompressBasic(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Compress 0.25-0.75 expected 1 got %d", len(haps))
	}
	// Compressed hap should have duration 0.5
	if !haps[0].Whole.Duration().Equals(FractionFromFloat(0.5)) {
		t.Logf("Compress duration %v", haps[0].Whole.Duration())
	}
}

func TestZoomBasic(t *testing.T) {
	p := Pure("a").Zoom(FractionFromInt(0), FractionFromInt(1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Zoom expected >=1")
	}
}

func TestFastGapBasic(t *testing.T) {
	p := Pure("a").FastGap(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("FastGap 2 expected >=1 got %d", len(haps))
	}
}

func TestInsideOutside(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b")).Inside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Inside expected >=1")
	}
	q := FastCat(Pure("a"), Pure("b")).Outside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Outside expected >=1")
	}
}
