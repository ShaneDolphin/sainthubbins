package core

import "testing"

func TestMJS_Port139_StackCatSlowcat(t *testing.T) {
	s := Stack(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Stack 4")
	}
	cat := Cat(Pure("a"), Pure("b"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
	fast := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(fast.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3")
	}
	slow := SlowCat(Pure("a"), Pure("b"))
	if len(slow.QueryArc(FractionFromInt(0), FractionFromInt(2))) != 2 {
		t.Fatalf("SlowCat 2 cycles 2")
	}
}

func TestMJS_Port139_SlowFastCompressZoom(t *testing.T) {
	s := Pure("a").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a b c").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastF 2")
	}
	c := Pure("bd").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0.25-0.75")
	}
	z := Pure("a b").Zoom(FractionFromFloat(0), FractionFromFloat(0.5))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0-0.5")
	}
}

func TestMJS_Port139_ChooseWithRandSegment(t *testing.T) {
	ch := Pure(0).Choose([]any{1, 2, 3})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose 1,2,3")
	}
	r := Rand().Segment(4)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Segment 4")
	}
	p := Pure("a").Choose([]any{Pure("b"), Pure("c")})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose Pattern b,c")
	}
}
