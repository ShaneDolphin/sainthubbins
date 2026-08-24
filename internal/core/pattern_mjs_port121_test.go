package core

import "testing"

func TestMJS_Port121_PolymeterTimeCatStepCat(t *testing.T) {
	// PolymeterSlowcat variant
	pm := PolymeterSlowcat(Pure("a"), Pure("b"), Pure("c"))
	if pm.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("PolymeterSlowcat")
	}
	// TimeCatWeighted
	tc := TimeCatWeighted([]any{1, Pure("a"), 2, Pure("b")})
	if len(tc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("TimeCatWeighted")
	}
	// StepCat via Polymeter? Use SlowCat alias check
	sc := SlowCat(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("SlowCat 1 per cycle")
	}
}

func TestMJS_Port121_FilterWhenStructMask(t *testing.T) {
	f := Pure("a b c").FilterValues(func(v any) bool {
		s, ok := v.(string)
		return ok && s != "b"
	})
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FilterValues not b")
	}
	keep := Pure("a b c d").KeepIf(Pure(true))
	if keep.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("KeepIf")
	}
	st := Pure("bd").Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if st.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct t f t f")
	}
	mask := Pure("a b c").Mask(FastCat(Pure(true), Pure(false), Pure(true)))
	if mask.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask t f t")
	}
}

func TestMJS_Port121_SuperimposeWithSlowDegrade(t *testing.T) {
	s := Pure("bd").Superimpose(func(q Pattern) Pattern { return q.Slow(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Slow 2")
	}
	d := Pure("bd sd").Degrade()
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Degrade nil")
	}
	sb := Pure("a b c").SometimesBy(0.5, func(q Pattern) Pattern { return q.Rev() })
	if sb.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.5")
	}
}
