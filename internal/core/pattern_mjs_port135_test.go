package core

import "testing"

func TestMJS_Port135_WithValueControls(t *testing.T) {
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["n"] = 3; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["n"] != 3 {
		t.Fatalf("WithValue n 3")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.5; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.5 {
		t.Fatalf("gain 0.5")
	}
	r := Stack(S("bd"), S("sd")).WithValue(func(v any) any { return v })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue 2")
	}
}

func TestMJS_Port135_HapDurationSpan(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Float64() != 2 {
		t.Fatalf("Duration 2")
	}
	cycles := span.SpanCycles()
	if len(cycles) != 2 {
		t.Fatalf("SpanCycles 2")
	}
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 1})
	if h.Part.Duration().Float64() != 2 {
		t.Fatalf("Hap Duration 2")
	}
}

func TestMJS_Port135_PatternFmapWithValue(t *testing.T) {
	p := Pure(5).Fmap(func(v any) any { return v.(int) * 3 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 15 {
		t.Fatalf("Fmap 15")
	}
	q := FastCat(Pure(1), Pure(2)).Fmap(func(v any) any {
		switch x := v.(type) {
		case int:
			return x * 10
		case float64:
			return x * 10
		default:
			return v
		}
	})
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Fmap Cat 2")
	}
}
