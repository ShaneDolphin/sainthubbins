package core

import "testing"

func TestMJS_Port140_WithValueControlsDegrade(t *testing.T) {
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["cutoff"] = 800; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["cutoff"] != 800 {
		t.Fatalf("cutoff 800")
	}
	q := Pure("sd").DegradeBy(0)
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0")
	}
	r := Stack(S("bd"), S("sd")).WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 1.0; return m
	})
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue gain 1.0 2")
	}
}

func TestMJS_Port140_HapDurationSpanCycles(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(1), FractionFromInt(4))
	if span.Duration().Float64() != 3 {
		t.Fatalf("Duration 3")
	}
	cycles := span.SpanCycles()
	if len(cycles) != 3 {
		t.Fatalf("SpanCycles 3 got %d", len(cycles))
	}
	h := NewHap(&span, span, 42, nil)
	if h.Value != 42 {
		t.Fatalf("Hap 42")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["pan"] = 0.5; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["pan"] != 0.5 {
		t.Fatalf("pan 0.5")
	}
}

func TestMJS_Port140_PatternFmapWithValueStack(t *testing.T) {
	p := Pure(10).Fmap(func(v any) any { return v.(int) + 7 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 17 {
		t.Fatalf("Fmap 17")
	}
	q := FastCat(Pure("a"), Pure("b")).Fmap(func(v any) any { return v.(string) + v.(string) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat Fmap 2")
	}
	r := Stack(Pure("x"), Pure("y"), Pure("z")).WithValue(func(v any) any { return v })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack WithValue 3")
	}
}
