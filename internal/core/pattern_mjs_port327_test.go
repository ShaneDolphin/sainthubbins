package core

import "testing"

func TestMJS_Port327_PolymeterTimeCatStepCatFourth(t *testing.T) {
	p := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat a,b")
	}
	q := TimeCatWeighted(1, Pure("a"), 2, Pure("b"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("TimeCatWeighted 1:a 2:b")
	}
	r := SlowCat(Pure("x"), Pure("y"), Pure("z"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat xyz 3")
	}
}

func TestMJS_Port327_FilterWhenStructMaskFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FilterValues(func(v any) bool { return v != "b" })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FilterValues !=b 2")
	}
	q := FastCat(Pure("a"), Pure("b")).KeepIf(Pure(true))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("KeepIf true 2")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Struct true false 2")
	}
	s := FastCat(Pure("a"), Pure("b")).Mask(Pure(true))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mask true 2")
	}
}

func TestMJS_Port327_SuperimposeWithSlowDegradeFourth(t *testing.T) {
	p := Pure("bd").Superimpose(func(q Pattern) Pattern { return q.Slow(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Slow2 <2")
	}
	q := Pure("a").DegradeBy(0.3)
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.3 nil")
	}
	r := Pure("b").SometimesBy(0.5, func(q Pattern) Pattern { return q.Rev() })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.5")
	}
}
