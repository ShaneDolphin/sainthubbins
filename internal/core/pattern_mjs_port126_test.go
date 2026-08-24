package core

import "testing"

func TestMJS_Port126_BinaryOpAddMulSubDiv(t *testing.T) {
	p := FastCat(Pure(10), Pure(20))
	q := p.Add(Pure(5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 5")
	}
	// Check values with float64 tolerance (Add uses float)
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value.(float64) < 14.9 || haps[0].Value.(float64) > 15.1 {
		t.Fatalf("Add value 15 got %v", haps[0].Value)
	}
	m := FastCat(Pure(3), Pure(4)).Mul(Pure(2))
	hm := m.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if hm[0].Value.(float64) < 5.9 || hm[1].Value.(float64) < 7.9 {
		t.Fatalf("Mul 2 got %v %v", hm[0].Value, hm[1].Value)
	}
	d := Pure(20).Div(Pure(4))
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 5 {
		t.Fatalf("Div 4 ->5")
	}
	s := Pure(10).Sub(Pure(3))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 7 {
		t.Fatalf("Sub 3 ->7")
	}
}

func TestMJS_Port126_TimeSpanHapWithState(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(3))
	cycles := span.SpanCycles()
	if len(cycles) != 3 {
		t.Fatalf("SpanCycles 3 got %d", len(cycles))
	}
	if span.Duration().Float64() != 3 {
		t.Fatalf("Duration 3")
	}
	h := NewHap(&span, span, "sd", map[string]any{"n": 1})
	if h.Value != "sd" {
		t.Fatalf("Hap sd")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["orbit"] = 5; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["orbit"] != 5 {
		t.Fatalf("WithContext orbit 5")
	}
}

func TestMJS_Port126_FastSlowChunkArp(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b")).FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Fast a b *2 =>4 got %d", len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	s := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Slow 2 =>2")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(q Pattern) Pattern { return q.Rev() })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2 Rev")
	}
	arp := Pure("c3 e3 g3").Arp("updown").FastF(FractionFromInt(2))
	if arp.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp updown Fast 2")
	}
}
