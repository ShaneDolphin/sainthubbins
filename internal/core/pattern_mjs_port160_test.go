package core

import "testing"

func TestMJS_Port160_ArpeggioUpDownSecond(t *testing.T) {
	up := Pure("c3 e3 g3").Arp("up")
	if len(up.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	down := Pure("c3 e3 g3").Arp("down")
	if len(down.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down")
	}
	updown := Pure("c3 e3 g3 b3").Arp("updown")
	if len(updown.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
}

func TestMJS_Port160_DegradeBySometimesOftenSecond(t *testing.T) {
	p := Pure("bd").DegradeBy(0.5)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.5 nil allow")
	}
	s := Pure("bd").Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes Fast 2")
	}
	o := Pure("bd").SometimesBy(0.75, func(q Pattern) Pattern { return q.Rev() })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.75")
	}
}

func TestMJS_Port160_HapWithContextValueSecond(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 9})
	if h.Context["orbit"] != 9 {
		t.Fatalf("orbit 9")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["delay"] = 0.4; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["delay"] != 0.4 {
		t.Fatalf("delay 0.4")
	}
	q := Pure(20).WithValue(func(v any) any { return v.(int) - 5 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 15 {
		t.Fatalf("WithValue 15")
	}
}
