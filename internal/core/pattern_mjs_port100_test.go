package core

import "testing"

func TestMJS_WhenOffSometimesDegrade2(t *testing.T) {
	p := Pure("bd")
	w := p.When(func(b bool) bool { return b }, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	o := w.Off(0.5, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("When Off") }
	s := p.Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Sometimes") }
}

func TestMJS_PatternTimeSpanHap4(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(1), FractionFromInt(2))
	if span.Duration().Float64() != 1 { t.Fatalf("Duration 1") }
	h := NewHap(&span, span, "x", map[string]any{"n": 1})
	if h.Value != "x" { t.Fatalf("Hap x") }
	if h.Context["n"] != 1 { t.Fatalf("Context n") }
}

func TestMJS_ClockAndSession3(t *testing.T) {
	c := NewClock(1.0)
	c.SetCPS(2.0)
	if c.CPS != 2 { t.Fatalf("Cps 2") }
	if c == nil { t.Fatalf("Clock nil") }
}
