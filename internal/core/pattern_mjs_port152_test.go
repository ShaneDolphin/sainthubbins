package core

import "testing"

func TestMJS_Port152_PatternPureStackCatSecond(t *testing.T) {
	p := Pure(5).Fmap(func(v any) any { return v.(int) * 2 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 10 {
		t.Fatalf("Fmap 10")
	}
	s := Stack(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	c := FastCat(Pure("a"), Pure("b"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2")
	}
}

func TestMJS_Port152_SlowFastCompressSecond(t *testing.T) {
	s := Pure("a b c d").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("bd").FastF(FractionFromInt(3))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastF 3")
	}
	c := Pure("bd sd").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0.25-0.75")
	}
	z := Pure("a b").Zoom(FractionFromFloat(0.5), FractionFromFloat(1))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0.5-1")
	}
}

func TestMJS_Port152_ChooseWithRandSecond(t *testing.T) {
	ch := Pure(0).Choose([]any{"a", "b", "c"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b c")
	}
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand")
	}
	deg := Pure("bd").DegradeBy(0.3)
	if deg.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.3 nil allow")
	}
}
