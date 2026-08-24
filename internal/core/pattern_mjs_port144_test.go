package core

import "testing"

func TestMJS_Port144_EuclidBjorklundWithOff(t *testing.T) {
	e := Pure("bd").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid 3,8")
	}
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund len 8 got %d", len(b))
	}
	off := Pure("bd").Euclid(3, 8).Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(off.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Euclid Off 0.25")
	}
	leg := Pure("bd").EuclidLegato(3, 8)
	if len(leg.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidLegato 3,8")
	}
}

func TestMJS_Port144_ScaleChordVoicing(t *testing.T) {
	s := Pure("c3").Scale("minor")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	ch := Pure("c3 e3 g3").Chord("m7")
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord m7")
	}
	// Scale with addition
	sc := Pure("c3").Scale("major").Add(Pure(2))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale major Add 2")
	}
}

func TestMJS_Port144_SignalRangeSineTri(t *testing.T) {
	s := Sine().Range(0, 100).Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,100 Slow 2")
	}
	tri := Tri().Range(-1, 1)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -1,1")
	}
	saw := Saw().Range(5, 10)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 5,10")
	}
	r := Rand().Range(0, 50)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,50")
	}
}
