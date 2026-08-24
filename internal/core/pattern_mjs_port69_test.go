package core

import "testing"

func TestMJS_License2(t *testing.T) {
	p := Pure("a").FastF(FractionFromInt(2)).SlowF(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("FastF SlowF chain expected non-empty")
	}
}

func TestMJS_PolymeterEuclid2(t *testing.T) {
	// Euclid edge 0,8 => 0 haps
	e0 := Pure("a").Euclid(0, 8)
	if len(e0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Euclid 0,8 expected 0")
	}
	// Euclid 8,8 => 8 haps
	e8 := Pure("a").Euclid(8, 8)
	if len(e8.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 8 {
		t.Fatalf("Euclid 8,8 expected 8 got %d", len(e8.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	e3 := Pure("a").Euclid(3, 8)
	if len(e3.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Euclid 3,8 expected 3")
	}
}

func TestMJS_SignalBypass2(t *testing.T) {
	s := Sine().Range(0, 1)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Sine Range expected non-empty")
	}
	// Bypass
	b := Pure("a").Bypass(Pure(false))
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Bypass false expected non-empty")
	}
}
