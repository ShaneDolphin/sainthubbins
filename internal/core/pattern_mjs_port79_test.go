package core

import "testing"

func TestMJS_RunBinary2(t *testing.T) {
	p := FastCat(Pure(1), Pure(2), Pure(3), Pure(4))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		t.Fatalf("FastCat 4 expected 4 got %d", len(haps))
	}
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2 expected 2")
	}
}

func TestMJS_SometimesDegradeOften2(t *testing.T) {
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes")
	}
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DegradeBy 0")
	}
	sb := Pure("a").SometimesBy(0, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(sb.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0")
	}
}

func TestMJS_StackCatPolymeterSteps2(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	fc := FastCat(Pure("a"), Pure("b"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat")
	}
}
