package core

import "testing"

func TestMJS_WithValueControls2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	q := p.WithValue(func(v any) any {
		m := v.(map[string]any)
		m["n"] = 2
		return m
	})
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("WithValue") }
	if haps[0].Value.(map[string]any)["n"] != 2 { t.Fatalf("n 2") }
}

func TestMJS_HapDurationSpan2(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Float64() != 2 { t.Fatalf("Duration 2") }
	cycles := span.SpanCycles()
	if len(cycles) != 2 { t.Fatalf("SpanCycles 2 got %d", len(cycles)) }
	h := NewHap(&span, span, "bd", nil)
	if h.Part.Duration().Float64() != 2 { t.Fatalf("Hap Duration") }
}

func TestMJS_PatternFmapWithValue3(t *testing.T) {
	p := Pure(10)
	q := p.Fmap(func(v any) any {
		switch x:=v.(type){case int: return x*2; case float64: return x*2; default: return v}
	})
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("Fmap") }
	v := haps[0].Value
	ok := false
	switch x:=v.(type){case int: ok = x==20; case float64: ok = x==20}
	if !ok { t.Fatalf("20 got %v", v) }
}
