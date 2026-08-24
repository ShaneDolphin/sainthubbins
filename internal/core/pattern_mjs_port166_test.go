package core

import "testing"

func TestMJS_Port166_WithValueControlsThird(t *testing.T) {
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.8; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.8 {
		t.Fatalf("gain 0.8")
	}
	q := S("hh:2").FirstCycle()
	if len(q) == 0 {
		t.Fatalf("S hh:2")
	}
	r := Stack(S("bd:1"), Gain(0.6))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack bd:1 Gain 2")
	}
}

func TestMJS_Port166_HapDurationSpanThird(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 5})
	if h.Context["orbit"] != 5 {
		t.Fatalf("orbit 5")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["pan"] = 0.2; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["pan"] != 0.2 {
		t.Fatalf("pan 0.2")
	}
	q := Pure(10).WithValue(func(v any) any { return v.(int) * 3 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 30 {
		t.Fatalf("WithValue 30")
	}
}

func TestMJS_Port166_PatternFmapWithValueThird(t *testing.T) {
	p := Pure("test").Fmap(func(v any) any { return v.(string) + "!" })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "test!" {
		t.Fatalf("test!")
	}
	q := FastCat(Pure(1), Pure(2)).Fmap(func(v any) any {
		switch x := v.(type) {
		case int:
			return x + 10
		case float64:
			return x + 10
		default:
			return v
		}
	})
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat Fmap 2")
	}
}
