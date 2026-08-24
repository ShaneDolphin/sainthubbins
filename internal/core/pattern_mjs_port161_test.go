package core

import "testing"

func TestMJS_Port161_PatternControlSValueThird(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	q := p.WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.9; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.9 {
		t.Fatalf("gain 0.9")
	}
	r := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["orbit"] = 3; return m
	})
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("S bd orbit 3")
	}
	s := Stack(Pure("a"), Pure("b")).WithValue(func(v any) any { return v })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack WithValue 2")
	}
}

func TestMJS_Port161_SignalPerlinWithSlowThird(t *testing.T) {
	s := Sine().Slow(FractionFromInt(3)).Range(0, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 3")
	}
	p := Perlin().Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 2")
	}
	r := Rand().Range(0, 10).FastF(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,10 Fast 2")
	}
}

func TestMJS_Port161_ArpWithFastSlowThird(t *testing.T) {
	p := Pure("c3 e3 g3 a3").Arp("up").FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up Fast 2")
	}
	q := Pure("c3 e3 g3").Arp("down").Slow(FractionFromInt(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down Slow 2")
	}
	r := Pure("c4 e4 g4").Arp("converge")
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge")
	}
}
