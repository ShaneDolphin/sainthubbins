package core

import "testing"

func TestMJS_PatternControlSValue2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 || haps[0].Value.(map[string]any)["s"]!="bd" { t.Fatalf("s bd") }
	q := p.WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"]=0.8; return m
	})
	haps = q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value.(map[string]any)["gain"] != 0.8 { t.Fatalf("gain") }
}

func TestMJS_SignalPerlinWithSlow2(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Sine Slow") }
	r := Rand().Slow(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Rand Slow") }
}

func TestMJS_ArpWithFastSlow2(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("updown")
	f := p.FastF(FractionFromInt(2))
	if f.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Arp Fast") }
	s := p.Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Arp Slow") }
}
