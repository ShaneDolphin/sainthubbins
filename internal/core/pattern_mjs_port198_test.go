package core

import "testing"

func TestMJS_Port198_PatternPureStackCatFourth(t *testing.T) {
	p := Pure("hello").Fmap(func(v any) any { return v.(string) + v.(string) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "hellohello" {
		t.Fatalf("hellohello")
	}
	q := Stack(Pure("bd"), Pure("sd"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	r := Stack(S("bd:1"), S("sd:2"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack bd:1 sd:2 2")
	}
	cat := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3")
	}
}

func TestMJS_Port198_SlowFastCompressFourth(t *testing.T) {
	s := Pure("a").Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Slow 2")
	}
	f := Pure("b").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastF 2 single empty")
	}
	c := Pure("c").Compress(FractionFromFloat(0), FractionFromFloat(0.5))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0-0.5")
	}
}

func TestMJS_Port198_ChooseWithRandFourth(t *testing.T) {
	ch := Pure(1).Choose([]any{"x", "y", "z"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose x y z")
	}
	r := Rand().Range(0, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,1")
	}
	s := Pure("a").Segment(4)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Segment 4")
	}
}
