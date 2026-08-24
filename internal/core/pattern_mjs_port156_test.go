package core

import "testing"

func TestMJS_Port156_ControlsWithValueSecond(t *testing.T) {
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.4; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.4 {
		t.Fatalf("gain 0.4")
	}
	q := S("sd").WithValue(func(v any) any {
		m := v.(map[string]any); m["n"] = 2; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["n"] != 2 {
		t.Fatalf("n 2")
	}
	r := Stack(S("bd"), S("sd"), S("hh")).WithValue(func(v any) any { return v })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3 WithValue")
	}
}

func TestMJS_Port156_HapDurationSpanSecondPlus(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(3), FractionFromInt(6))
	if span.Duration().Float64() != 3 {
		t.Fatalf("Duration 3")
	}
	cycles := span.SpanCycles()
	if len(cycles) != 3 {
		t.Fatalf("SpanCycles 3")
	}
	h := NewHap(&span, span, "bd", map[string]any{"velocity": 0.9})
	if h.Context["velocity"] != 0.9 {
		t.Fatalf("velocity 0.9")
	}
}

func TestMJS_Port156_PatternFmapWithValueSecondPlus(t *testing.T) {
	p := Pure("hello").Fmap(func(v any) any { return v.(string) + " world" })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "hello world" {
		t.Fatalf("hello world")
	}
	q := FastCat(Pure(1), Pure(2), Pure(3)).Fmap(func(v any) any {
		switch x := v.(type) {
		case int:
			return x * 2
		case float64:
			return x * 2
		default:
			return v
		}
	})
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat Fmap 3")
	}
}
