package core

import "testing"

func TestMJS_Port132_SignalSineTriSaw(t *testing.T) {
	s := Sine().Range(0, 10).Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,10 Slow 2")
	}
	tri := Tri().Range(-5, 5)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -5,5")
	}
	saw := Saw().Range(0, 1).FastF(FractionFromInt(2))
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,1 Fast 2")
	}
}

func TestMJS_Port132_PatternWhenWithOff(t *testing.T) {
	p := Pure("a b c").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	q := Pure("bd").Off(0.5, func(pat Pattern) Pattern { return pat.Add(Pure(5)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.5")
	}
	r := Pure("a").When(false, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
}

func TestMJS_Port132_HapContextWithValue(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"n": 2})
	if h.Value != "bd" || h.Context["n"] != 2 {
		t.Fatalf("Hap bd n 2")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["gain"] = 0.7; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["gain"] != 0.7 {
		t.Fatalf("WithContext gain 0.7")
	}
	q := Pure(10).WithValue(func(v any) any { return v.(int) + 5 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 15 {
		t.Fatalf("WithValue 15")
	}
}
