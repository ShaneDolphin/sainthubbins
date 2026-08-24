package core

import "testing"

func TestMJS_Port163_SignalSineTriSawThird(t *testing.T) {
	s := Sine().Range(0, 10)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,10")
	}
	tri := Tri().Range(0, 1).Slow(FractionFromInt(2))
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range Slow 2")
	}
	saw := Saw().Range(-5, 5).FastF(FractionFromInt(2))
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range -5,5 Fast 2")
	}
	perlin := Perlin().Range(0, 100).Slow(FractionFromInt(1))
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range 0,100")
	}
}

func TestMJS_Port163_PatternWhenWithOffThird(t *testing.T) {
	p := Pure("a b").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev 2")
	}
	q := Pure("bd").Off(0.125, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125 Fast 2")
	}
	r := Pure("a").When(false, func(q Pattern) Pattern { return q.Add(Pure(5)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
}

func TestMJS_Port163_HapContextWithValueThird(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "sd", map[string]any{"orbit": 4})
	if h.Context["orbit"] != 4 {
		t.Fatalf("orbit 4")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["gain"] = 0.6; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["gain"] != 0.6 {
		t.Fatalf("gain 0.6")
	}
	q := Pure(map[string]any{"s": "hh", "n": 2}).WithValue(func(v any) any {
		m := v.(map[string]any); m["pan"] = 0.8; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["pan"] != 0.8 {
		t.Fatalf("pan 0.8")
	}
}
