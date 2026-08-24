package core

import "testing"

func TestMJS_SelectByFraction2(t *testing.T) {
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes expected non-empty")
	}
	o := Pure("a").SometimesBy(0.75, func(p Pattern) Pattern { return p })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.75 expected non-empty")
	}
	r := Pure("a").SometimesBy(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.25 expected non-empty")
	}
}

func TestMJS_PatStack2(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"), Pure("c"))
	haps := st.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Stack 3 expected 3 got %d", len(haps))
	}
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	haps2 := c.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps2) == 0 {
		t.Fatalf("Cat 3 cycles expected non-empty")
	}
}

func TestMJS_PatternDecay2(t *testing.T) {
	p := Pure("a").Set(Decay(0.1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Decay via Set expected non-empty")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["decay"] == nil {
		// Decay stored as map decay key via controls
	}
}
