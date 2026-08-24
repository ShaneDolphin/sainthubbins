package core

import "testing"

func TestMJS_Port215_RandChooseDegradeShuffleFourth(t *testing.T) {
	r := Rand().Range(0, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand")
	}
	ch := Pure(1).Choose([]any{"x", "y", "z"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose x y z")
	}
	d := Pure("a").DegradeBy(0.5)
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.5 nil")
	}
	sh := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Shuffle(4)
	if len(sh.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Shuffle 4")
	}
}

func TestMJS_Port215_MultiMapWithValueSFourth(t *testing.T) {
	p := Pure(2).WithValue(func(v any) any { return v.(int) * 3 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 6 {
		t.Fatalf("WithValue 6 got %v", p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value)
	}
	q := S("bd").Set(Gain(0.8))
	v := q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["s"] != "bd" || v["gain"] != 0.8 {
		t.Fatalf("S bd gain 0.8 got %v", v)
	}
	r := Stack(S("bd").Set(Gain(0.5)), S("sd").Set(Gain(0.7)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack gain 2")
	}
}

func TestMJS_Port215_HapSpanTimeValueFourth(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Cmp(FractionFromInt(2)) != 0 {
		t.Fatalf("Duration 2")
	}
	spans := span.SpanCycles()
	if len(spans) != 2 {
		t.Fatalf("SpanCycles 2 got %d", len(spans))
	}
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 1})
	if h.Value != "bd" || h.Context["orbit"] != 1 {
		t.Fatalf("Hap bd orbit 1")
	}
	f := FractionFromFloat(1.5)
	if f.Float64() < 1.4 || f.Float64() > 1.6 {
		t.Fatalf("Fraction 1.5")
	}
}
