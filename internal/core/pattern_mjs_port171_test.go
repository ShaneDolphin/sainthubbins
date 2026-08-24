package core

import "testing"

func TestMJS_Port171_PatternControlSValueFourth(t *testing.T) {
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["orbit"] = 4; return m
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("S bd orbit 4")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.3; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.3 {
		t.Fatalf("gain 0.3")
	}
	r := Stack(S("bd"), Gain(0.7), Pan(0.2))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3 controls")
	}
}

func TestMJS_Port171_SignalPerlinWithSlowFourth(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2)).Range(-1, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 2 -1,1")
	}
	p := Perlin().Slow(FractionFromInt(3)).Range(0, 5)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 3")
	}
	r := Rand().Range(0, 100).FastF(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range Fast 2")
	}
}

func TestMJS_Port171_ArpWithFastSlowFourth(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up").Slow(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up Slow 2")
	}
	q := Pure("c3 e3 g3 b3").Arp("down").FastF(FractionFromInt(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down Fast 2")
	}
	r := Pure("c2 e2 g2").Arp("converge")
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge")
	}
}
