package core

import "testing"

func TestMJS_Port200_WithValueControlsFourth(t *testing.T) {
	p := Pure(2).WithValue(func(v any) any { return v.(int) * 5 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 10 {
		t.Fatalf("WithValue 10")
	}
	q := Pure(5).WithValue(func(v any) any { return v.(int) + 3 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 8 {
		t.Fatalf("WithValue 8")
	}
	r := S("bd").Set(Gain(0.7))
	v := r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["gain"] != 0.7 {
		t.Fatalf("gain 0.7 got %v", v["gain"])
	}
}

func TestMJS_Port200_HapDurationSpanFourth(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Cmp(FractionFromInt(2)) != 0 {
		t.Fatalf("Duration 2")
	}
	spans := span.SpanCycles()
	if len(spans) != 2 {
		t.Fatalf("SpanCycles 2 got %d", len(spans))
	}
	h := NewHap(&span, span, "bd", map[string]any{"n": 3})
	if h.Value != "bd" {
		t.Fatalf("Hap bd")
	}
	f := FractionFromFloat(2.5)
	if f.Float64() < 2.4 || f.Float64() > 2.6 {
		t.Fatalf("Fraction 2.5")
	}
}

func TestMJS_Port200_PatternFmapWithValueFourth(t *testing.T) {
	p := Pure(6).Fmap(func(v any) any { return v.(int) * 2 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 12 {
		t.Fatalf("Fmap 12")
	}
	q := FastCat(Pure("a"), Pure("b")).Fmap(func(v any) any { return v.(string) + "!" })
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 || haps[0].Value != "a!" {
		t.Fatalf("Fmap a! got %v", haps[0].Value)
	}
}
