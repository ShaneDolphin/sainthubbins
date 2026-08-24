package core

import "testing"

func TestMJS_Port133_PatternPureStackCat(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v.(string) + "1" })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "a1" {
		t.Fatalf("Fmap a1")
	}
	s := Stack(Pure("bd"), Pure("sd"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	c := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3")
	}
}

func TestMJS_Port133_SlowFastCompress(t *testing.T) {
	s := Pure("a b c").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a b c d").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastF 2")
	}
	c := Pure("bd sd").Compress(FractionFromFloat(0), FractionFromFloat(0.5))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0-0.5")
	}
}

func TestMJS_Port133_ChooseWithRand(t *testing.T) {
	ch := Pure(0).Choose([]any{"x", "y", "z"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose x y z")
	}
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand nil")
	}
	seg := Rand().Segment(4)
	if seg.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Segment 4")
	}
}
