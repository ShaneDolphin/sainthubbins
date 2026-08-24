package core

import "testing"

func TestMJS_StackWithRest2(t *testing.T) {
	s := Stack(Pure("a"), Silence())
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Stack silence+a expected 1 got %d", len(haps))
	}
}

func TestMJS_CatWithSilence2(t *testing.T) {
	c := Cat(Pure("a"), Silence())
	haps := c.QueryArc(FractionFromInt(0), FractionFromInt(2))
	// Cat = SlowCat 1 per cycle, so 2 cycles: first a, second silence => 1 hap
	if len(haps) != 1 {
		t.Fatalf("Cat a silence 2 cycles expected 1 got %d", len(haps))
	}
}

func TestMJS_SuperimposeWithOff2(t *testing.T) {
	p := Pure("a").Superimpose(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Superimpose FastF2 expected >=2 got %d", len(haps))
	}
	o := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 expected >=2")
	}
}
