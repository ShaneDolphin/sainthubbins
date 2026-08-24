package core

import "testing"

func TestMJS_Port116_RandChooseDegradeShuffle(t *testing.T) {
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand nil")
	}
	ch := Pure(0).Choose([]any{"a", "b", "c"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose")
	}
	d := Pure("bd").Degrade()
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Degrade nil")
	}
	sh := Pure("a b c d").Shuffle(4)
	if len(sh.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Shuffle 4")
	}
}

func TestMJS_Port116_MultiMapWithValueS(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	q := p.WithValue(func(v any) any {
		m := v.(map[string]any)
		m["n"] = 2
		return m
	})
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(map[string]any)["n"] != 2 {
		t.Fatalf("WithValue n:2")
	}
	r := S("hh").WithValue(func(v any) any {
		m := v.(map[string]any)
		m["gain"] = 0.8
		return m
	})
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("S WithValue gain")
	}
	s2 := Stack(S("bd"), S("sd")).WithValue(func(v any) any { return v })
	if len(s2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue 2")
	}
}

func TestMJS_Port116_HapSpanTimeValue(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Float64() != 2 {
		t.Fatalf("Duration 2")
	}
	cycles := span.SpanCycles()
	if len(cycles) != 2 {
		t.Fatalf("SpanCycles 2 got %d", len(cycles))
	}
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 1})
	if h.Context["orbit"] != 1 {
		t.Fatalf("orbit 1")
	}
	p := Pure(h.Value).WithContext(func(m map[string]any) map[string]any {
		m["s"] = "bd"
		return m
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("WithContext")
	}
	// Fraction roundtrip
	f := FractionFromFloat(1.5)
	if f.Float64() < 1.4 || f.Float64() > 1.6 {
		t.Fatalf("FractionFromFloat 1.5")
	}
}
