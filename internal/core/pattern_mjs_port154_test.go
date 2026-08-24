package core

import "testing"

func TestMJS_Port154_WithValueControlsSecond(t *testing.T) {
	p := S("hh").WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.7; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.7 {
		t.Fatalf("gain 0.7")
	}
	q := Pure(map[string]any{"s": "bd"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["n"] = 5; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["n"] != 5 {
		t.Fatalf("n 5")
	}
	r := Stack(S("bd"), S("sd")).WithValue(func(v any) any {
		m := v.(map[string]any); m["orbit"] = 1; return m
	})
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue orbit 1 2")
	}
}

func TestMJS_Port154_HapDurationSpanSecond(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(2), FractionFromInt(5))
	if span.Duration().Float64() != 3 {
		t.Fatalf("Duration 3")
	}
	h := NewHap(&span, span, "sd", map[string]any{"orbit": 2})
	if h.Context["orbit"] != 2 {
		t.Fatalf("orbit 2")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["delay"] = 0.2; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["delay"] != 0.2 {
		t.Fatalf("delay 0.2")
	}
}

func TestMJS_Port154_PatternFmapWithValueSecond(t *testing.T) {
	p := Pure(3).Fmap(func(v any) any { return v.(int) * 4 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 12 {
		t.Fatalf("Fmap 12")
	}
	q := FastCat(Pure(10), Pure(20)).Fmap(func(v any) any {
		switch x := v.(type) {
		case int:
			return x + 1
		case float64:
			return x + 1
		default:
			return v
		}
	})
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat Fmap 2")
	}
}
