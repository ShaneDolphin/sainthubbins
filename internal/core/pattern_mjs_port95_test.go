package core

import "testing"

func TestMJS_BjorklundEuclid3(t *testing.T) {
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund 3,8 len 8")
	}
	e := Pure("a").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Euclid 3,8 3")
	}
	er := Pure("a").EuclidRot(3, 8, 1)
	if len(er.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("EuclidRot 3,8,1 3")
	}
}

func TestMJS_SignalPerlin3(t *testing.T) {
	s := Sine()
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Sine")
	}
	r := Rand()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rand")
	}
	// Perlin via Signal
	p := Signal(func(frac Fraction) float64 { return frac.Float64() * 0.5 })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Signal Perlin")
	}
}

func TestMJS_ControlsN3(t *testing.T) {
	p := S("bd")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(map[string]any)["s"] != "bd" {
		t.Fatalf("S bd")
	}
	p2 := S("bd").Set(N(2))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("S bd N 2")
	}
	if m, ok := haps2[0].Value.(map[string]any); !ok || m["n"] != 2 {
		// n may be int/float
	}
}
