package core

import "testing"

func TestMJS_Port149_EuclidBjorklundSignalPort(t *testing.T) {
	e := Pure("bd").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Euclid 3,8 =>3 got %d", len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	b := Bjorklund(3, 8)
	count := 0
	for _, v := range b {
		if v != 0 {
			count++
		}
	}
	if count != 3 {
		t.Fatalf("Bjorklund 3,8 count 3 got %d", count)
	}
	s := Pure("bd").EuclidLegato(3, 8)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidLegato 3,8")
	}
	rot := Pure("bd").EuclidRot(3, 8, 2)
	if len(rot.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidRot 3,8,2")
	}
}

func TestMJS_Port149_ScaleVoicingChordPort(t *testing.T) {
	s := Pure("c3").Scale("minor")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	ch := Pure("c4").Chord("minor")
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord minor")
	}
	v := Pure("c3 e3 g3").Chord("major")
	if len(v.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord major seq")
	}
}

func TestMJS_Port149_SignalRangeSinePort(t *testing.T) {
	s := Sine().Range(0, 100)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Sine Range 0,100")
	}
	tri := Tri().Range(-10, 10)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -10,10")
	}
	perlin := Perlin().Range(0, 1).Slow(FractionFromInt(2))
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range 0,1 Slow 2")
	}
}
