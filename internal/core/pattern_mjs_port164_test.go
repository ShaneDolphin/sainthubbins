package core

import "testing"

func TestMJS_Port164_PatternPureStackCatThird(t *testing.T) {
	p := Pure("hello").Fmap(func(v any) any { return v.(string) + "!" })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "hello!" {
		t.Fatalf("hello!")
	}
	s := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	c := Cat(Pure("x"), Pure("y"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
	fc := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastCat 4")
	}
}

func TestMJS_Port164_SlowFastCompressThird(t *testing.T) {
	s := Pure("a b c d").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(4))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastF 4")
	}
	c := Pure("bd sd").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0.25-0.75")
	}
}

func TestMJS_Port164_ChooseWithRandThird(t *testing.T) {
	ch := Pure(0).Choose([]any{"a", "b", "c", "d"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose 4")
	}
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand")
	}
	deg := Pure("bd sd").DegradeBy(0.5)
	if deg.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("DegradeBy 0.5 nil allow")
	}
	s := Pure("a").Sometimes(func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sometimes")
	}
}
