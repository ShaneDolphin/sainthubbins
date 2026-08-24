package core

import "testing"

func TestMJS_EveryOffWhen3(t *testing.T) {
	e := Pure("a").Every(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Every 2")
	}
	o := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25")
	}
	w := Pure("a").When(Pure(false), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
}

func TestMJS_FilterHapsValues2(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	f := p.FilterHaps(func(h Hap) bool { return h.Value == "a" })
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("FilterHaps a 1")
	}
	fv := p.FilterValues(func(v any) bool { return v == "b" })
	if len(fv.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("FilterValues b 1")
	}
}

func TestMJS_SpliceFit2(t *testing.T) {
	sl := Slice(2, Pure(0), Pure("a"))
	if len(sl.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice 2")
	}
	f := Pure("a").Fit()
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Fit")
	}
}
