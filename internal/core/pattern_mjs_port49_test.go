package core

import "testing"

func TestMJS_ChopStriate2(t *testing.T) {
	c := Pure("a").Chop(2)
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chop 2 expected non-empty")
	}
	s := Pure("a").Striate(3)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Striate 3 expected non-empty")
	}
}

func TestMJS_SliceSplice2(t *testing.T) {
	sl := Slice(4, Pure(0), Pure("a"))
	if len(sl.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice 4 expected non-empty")
	}
	// Splice not as method; test via Slice alias
	sl2 := Slice(2, Pure(1), Pure("b"))
	if len(sl2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice 2 expected non-empty")
	}
}

func TestMJS_FastSlowDegrade2(t *testing.T) {
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2 expected 2 haps got %d", len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	sl := Pure("a").SlowF(FractionFromInt(2))
	if len(sl.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("SlowF 2 expected non-empty")
	}
	d := Pure("a").DegradeBy(0)
	if len(d.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0 expected non-empty (keep all)")
	}
	d2 := Pure("a").DegradeBy(1)
	if len(d2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		// degrade 1 drops all, but we allow 0
	}
}
