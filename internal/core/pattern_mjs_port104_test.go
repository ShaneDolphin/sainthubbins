package core

import "testing"

func TestMJS_PatternPureStackCat2(t *testing.T) {
	p := Pure("a")
	s := Stack(p, p.FastF(FractionFromInt(2)))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Stack FastF") }
	c := Cat(Pure("x"), Pure("y"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(2)))<2 { t.Fatalf("Cat") }
}

func TestMJS_SlowFastCompress2(t *testing.T) {
	p := Pure("bd").Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Slow") }
	f := Pure("sd").FastF(FractionFromInt(2))
	if f.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("FastF") }
}

func TestMJS_ChooseWithRand2(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b"})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("ChooseCycles") }
	r := Rand().Segment(4)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Rand Segment") }
}
