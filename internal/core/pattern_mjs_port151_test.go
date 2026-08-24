package core

import "testing"

func TestMJS_Port151_SignalSineTriSawSecond(t *testing.T) {
	s := Sine().Range(0, 1).FastF(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range Fast 2")
	}
	tri := Tri().Range(-1, 1).Slow(FractionFromInt(2))
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Slow 2")
	}
	saw := Saw().Range(5, 15)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 5,15")
	}
	perlin := Perlin().Slow(FractionFromInt(4)).Range(0, 10)
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 4 Range 0,10")
	}
}

func TestMJS_Port151_PatternWhenWithOffSecond(t *testing.T) {
	p := Pure(1).When(true, func(q Pattern) Pattern { return q.Add(Pure(10)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 11 {
		t.Fatalf("When true Add 10 ->11")
	}
	q := Pure("bd").Off(0.25, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 Fast 2")
	}
	r := Pure("a b c").When(false, func(q Pattern) Pattern { return q.Rev() })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When false")
	}
}

func TestMJS_Port151_HapContextWithValueSecond(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	h := NewHap(&span, span, "hh", map[string]any{"delay": 0.3})
	if h.Context["delay"] != 0.3 {
		t.Fatalf("delay 0.3")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["room"] = 0.8; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["room"] != 0.8 {
		t.Fatalf("room 0.8")
	}
	q := Pure(7).WithValue(func(v any) any { return v.(int) * 3 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 21 {
		t.Fatalf("WithValue 21")
	}
}
