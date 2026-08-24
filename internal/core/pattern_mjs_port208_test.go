package core

import "testing"

func TestMJS_Port208_EveryOffWhenChunkFourth(t *testing.T) {
	p := Pure("a").Every(3, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(3))) == 0 {
		t.Fatalf("Every 3")
	}
	q := Pure("bd").Off(0.125, func(pat Pattern) Pattern { return pat.Rev() })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125 <2")
	}
	r := Pure("x").When(false, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When false")
	}
	s := Pure("a b c d e").Chunk(3, func(pat Pattern) Pattern { return pat.Rev() })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 3")
	}
}

func TestMJS_Port208_ScaleChordTransposeFourth(t *testing.T) {
	s := Pure("c3").Scale("minor")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	ch := Pure("c3").Chord("maj7")
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord maj7")
	}
	tr := Pure("c3").Transpose(2)
	if len(tr.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Transpose 2")
	}
}

func TestMJS_Port208_PatternAddMulDivModFourth(t *testing.T) {
	p := Pure(3).Add(Pure(4))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 7 {
		t.Fatalf("Add 3+4=7")
	}
	q := Pure(5).Mul(Pure(2))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 10 {
		t.Fatalf("Mul 5*2=10")
	}
	r := Pure(10).Div(Pure(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 5 {
		t.Fatalf("Div 10/2=5")
	}
	s := Pure(10).Mod(Pure(3))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 1 {
		t.Fatalf("Mod 10 Mod 3 is 1")
	}
}
