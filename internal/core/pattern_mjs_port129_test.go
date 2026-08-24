package core

import "testing"

func TestMJS_Port129_PatternAddMulStructure(t *testing.T) {
	p := FastCat(Pure(1), Pure(2), Pure(3))
	mul := p.Mul(Pure(10))
	haps := mul.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Mul 10 len 3 got %d", len(haps))
	}
	add := FastCat(Pure(10), Pure(20)).Add(Pure(5))
	if len(add.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 5")
	}
	sub := Pure(100).Sub(Pure(40))
	if sub.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 60 {
		t.Fatalf("Sub 40 ->60")
	}
}

func TestMJS_Port129_WhenOffSometimesDegrade(t *testing.T) {
	w := Pure("bd").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true Fast 2")
	}
	o := Pure("bd").Off(0.25, func(q Pattern) Pattern { return q.Add(Pure(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25")
	}
	s := Pure("bd").Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes Rev")
	}
	d := Pure("bd").DegradeBy(0.5)
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.5 nil but may be empty? allow nil")
	}
}

func TestMJS_Port129_PatternTimeSpanHap(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if span.Duration().Float64() != 2 {
		t.Fatalf("Duration 2")
	}
	h := NewHap(&span, span, "hh", map[string]any{"gain": 0.8})
	if h.Value != "hh" || h.Context["gain"] != 0.8 {
		t.Fatalf("Hap hh gain 0.8")
	}
	p := Pure("a").Early(FractionFromFloat(0.5))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Early 0.5")
	}
	q := Pure("a").Late(FractionFromFloat(0.5))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Late 0.5")
	}
}
