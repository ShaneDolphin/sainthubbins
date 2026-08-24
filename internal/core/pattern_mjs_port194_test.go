package core

import "testing"

func TestMJS_Port194_EuclidBjorklundWithOffFourth(t *testing.T) {
	e := Pure("bd").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid 3,8")
	}
	b := Bjorklund(3, 8)
	cnt := 0
	for _, v := range b {
		if v != 0 {
			cnt++
		}
	}
	if cnt != 3 {
		t.Fatalf("Bjorklund 3 !=3 %d", cnt)
	}
	o := Pure("bd").Euclid(3, 8).Off(0.25, func(p Pattern) Pattern { return p.Rev() })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Euclid Off 0.25 <2")
	}
	l := Pure("bd").EuclidLegato(3, 8)
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EuclidLegato 3,8")
	}
}

func TestMJS_Port194_ScaleChordVoicingFourth(t *testing.T) {
	s := Pure("c3").Scale("minor")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	ch := Pure("c3").Chord("m7")
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord m7")
	}
	a := Pure("c3").Add(Pure(2))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Add 2 Scale test")
	}
}

func TestMJS_Port194_SignalRangeSineTriFourth(t *testing.T) {
	s := Sine().Range(0, 10)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,10")
	}
	tri := Tri().Range(-5, 5)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -5,5")
	}
	saw := Saw().Range(0, 100)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,100")
	}
	r := Rand().Range(-1, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range -1,1")
	}
}
