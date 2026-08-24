package core

import "testing"

func TestMJS_StackCatSteps2(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	fc := FastCat(Pure("a"), Pure("b"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2")
	}
	ws := Pure("a").WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	if ws.Steps == nil || !ws.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("WithSteps *2")
	}
}

func TestMJS_TimeCatPolymeterSteps2(t *testing.T) {
	tc := TimeCatWeighted(1, Pure("a"), 2, Pure("b"))
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("TimeCatWeighted 1,2")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat")
	}
}

func TestMJS_RandDegradeSometimes2(t *testing.T) {
	r := Rand()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rand")
	}
	d := Pure("a").Degrade()
	haps := d.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 0 || len(haps) > 1 {
		t.Fatalf("Degrade len 0-1")
	}
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes")
	}
}
