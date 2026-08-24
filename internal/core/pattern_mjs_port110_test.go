package core

import "testing"

func TestMJS_ArpeggioUpDown2(t *testing.T) {
	up := Pure("c3 e3 g3").Arp("up")
	if len(up.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Arp up") }
	down := Pure("c3 e3 g3").Arp("down")
	if len(down.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Arp down") }
}

func TestMJS_DegradeBySometimesOften2(t *testing.T) {
	p := Pure("bd").DegradeBy(0.5)
	s := p.Sometimes(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Sometimes") }
	o := p.SometimesBy(0.5, func(q Pattern) Pattern { return q.Rev() })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("SometimesBy") }
}

func TestMJS_HapWithContextValue3(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"n": 5})
	if h.Context["n"] != 5 { t.Fatalf("n 5") }
	p := Pure(h.Value).WithContext(func(m map[string]any) map[string]any { m["s"]="bd"; return m })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("WithContext") }
}
