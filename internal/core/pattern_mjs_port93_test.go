package core

import "testing"

func TestMJS_SuperimposeLayer2(t *testing.T) {
	s := Pure("a").Superimpose(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose FastF2")
	}
	l := Pure("a").Layer(func(p Pattern) Pattern { return p }, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Layer 2")
	}
}

func TestMJS_DiscreteOnsets2(t *testing.T) {
	p := Stack(Pure("a"), Silence())
	d := p.DiscreteOnly()
	if len(d.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DiscreteOnly 1")
	}
	o := p.OnsetsOnly()
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("OnsetsOnly 1")
	}
}

func TestMJS_RollWithChildren2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	ch := p.Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2")
	}
	// ChunkBack
	cb := p.ChunkBack(2, func(pat Pattern) Pattern { return pat })
	if len(cb.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ChunkBack 2")
	}
}
