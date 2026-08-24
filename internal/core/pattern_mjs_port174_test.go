package core

import "testing"

func TestMJS_Port174_PatternPureStackCatFourth(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v.(string) + v.(string) + v.(string) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "aaa" {
		t.Fatalf("Fmap aaa")
	}
	s := Stack(Pure("bd"), Pure("sd"), Pure("hh"), Pure("oh"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Stack 4")
	}
	c := FastCat(Pure("a"), Pure("b"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2")
	}
}

func TestMJS_Port174_SlowFastCompressFourth(t *testing.T) {
	s := Pure("a b").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("bd sd hh").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastF 2")
	}
	c := Pure("a b c").Compress(FractionFromFloat(0.5), FractionFromFloat(1))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0.5-1")
	}
	z := Pure("a b c d").Zoom(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0.25-0.75")
	}
}

func TestMJS_Port174_ChooseWithRandFourth(t *testing.T) {
	ch := Pure(0).Choose([]any{"x", "y"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose x y")
	}
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand")
	}
	deg := Pure("a b c").DegradeBy(0.2)
	if deg.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.2 nil allow")
	}
	s := Rand().Segment(3)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Segment 3")
	}
}
