package core

import "testing"

func TestMJS_SelectByFraction3(t *testing.T) {
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes")
	}
	o := Pure("a").SometimesBy(0.75, func(p Pattern) Pattern { return p })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.75")
	}
}

func TestMJS_PatStack3(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"), Pure("c"))
	haps := st.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Stack 3 expected 3 got %d", len(haps))
	}
}

func TestMJS_PatternDecay3(t *testing.T) {
	p := Pure("a").Set(Decay(0.1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Decay 0.1")
	}
}
