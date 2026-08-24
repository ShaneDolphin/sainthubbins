package core

import "testing"

func TestMJS_Port170_ArpeggioUpDownThird(t *testing.T) {
	up := Pure("c3 e3 g3 b3").Arp("up")
	if len(up.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	down := Pure("c3 e3 g3").Arp("down")
	if len(down.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down")
	}
	updown := Pure("c3 e3 g3").Arp("updown")
	if len(updown.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
	diverge := Pure("c3 e3 g3 b3").Arp("diverge")
	if len(diverge.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp diverge")
	}
}

func TestMJS_Port170_DegradeBySometimesOftenThird(t *testing.T) {
	p := Pure("bd").DegradeBy(0.2)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.2 nil allow")
	}
	s := Pure("bd").SometimesBy(0.3, func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SometimesBy 0.3")
	}
	o := Pure("bd").Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes")
	}
}

func TestMJS_Port170_HapWithContextValueThird(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"velocity": 0.8})
	if h.Context["velocity"] != 0.8 {
		t.Fatalf("velocity 0.8")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["orbit"] = 2; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["orbit"] != 2 {
		t.Fatalf("orbit 2")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["pan"] = 0.5; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["pan"] != 0.5 {
		t.Fatalf("pan 0.5")
	}
}
