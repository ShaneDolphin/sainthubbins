package core

import "testing"

func TestPatternWithValue(t *testing.T) {
	p := Pure(1).Fmap(func(v any) any { return toFloat(v) * 2 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps[0].Value) != 2 {
		t.Fatalf("Fmap *2")
	}
}

func TestPatternStackCat(t *testing.T) {
	p1 := Pure("a")
	p2 := Pure("b")
	stack := Stack(p1, p2)
	if len(stack.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("stack 2")
	}
	cat := FastCat(p1, p2)
	haps := cat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("fastcat 2 got %d", len(haps))
	}
}

func TestPatternFastSlow(t *testing.T) {
	p := Pure("x")
	fast := p.FastF(FractionFromInt(2))
	if len(fast.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fast2")
	}
	slow := p.SlowF(FractionFromInt(2))
	// Slow 2 over 1 cycle gives 1 hap? Actually Slow 2 over 1 cycle queries 0.5 cycles
	haps := slow.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("slow2 no haps")
	}
}

func TestPatternEuclidPoly(t *testing.T) {
	pat := Pure("bd").Euclid(5, 8)
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 5 {
		t.Fatalf("euclid 5,8 expected 5 got %d", len(haps))
	}
}
