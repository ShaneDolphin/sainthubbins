package core

import "testing"

func TestMJS_SignalSineTriSaw2(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Sine") }
	tr := Tri().Slow(FractionFromInt(2))
	if tr.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Tri") }
	sw := Saw().Slow(FractionFromInt(2))
	if sw.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Saw") }
}

func TestMJS_PatternWhenWithOff2(t *testing.T) {
	p := Pure("a b c")
	w := p.When(func(b bool) bool { return true }, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("When") }
	o := p.Off(0.25, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Off") }
}

func TestMJS_HapContextWithValue2(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, 42, map[string]any{"s": "bd"})
	if h.Value != 42 { t.Fatalf("42") }
	haps := Pure(h.Value).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("Pure hap") }
}
