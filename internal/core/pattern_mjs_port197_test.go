package core

import "testing"

func TestMJS_Port197_SignalSineTriSawFourth(t *testing.T) {
	s := Sine().Range(-5, 5).Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range -5,5 Slow 2")
	}
	tri := Tri().Range(0, 1).FastF(FractionFromInt(2))
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range 0,1 Fast 2")
	}
	saw := Saw().Range(0, 100)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,100")
	}
	perlin := Perlin().Range(-1, 1).Slow(FractionFromInt(3))
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range -1,1 Slow 3")
	}
}

func TestMJS_Port197_PatternWhenWithOffFourth(t *testing.T) {
	p := Pure(1).When(true, func(q Pattern) Pattern { return q.Add(Pure(5)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 6 {
		t.Fatalf("When true Add 5 ->6")
	}
	q := Pure("bd").Off(0.5, func(pat Pattern) Pattern { return pat.Rev() })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.5 Rev")
	}
	r := Pure("a b c").When(false, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When false")
	}
}

func TestMJS_Port197_HapContextWithValueFourth(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(1), FractionFromInt(3))
	h := NewHap(&span, span, "sd", map[string]any{"n": 7})
	if h.Context["n"] != 7 {
		t.Fatalf("n 7")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["orbit"] = 8; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["orbit"] != 8 {
		t.Fatalf("orbit 8")
	}
	q := Pure(5).WithValue(func(v any) any { return v.(int) + 10 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 15 {
		t.Fatalf("WithValue 15")
	}
}
