package core

import "testing"

func TestMJS_FastGapCompress2(t *testing.T) {
	f := Pure("a").FastGap(0.5)
	haps := f.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("FastGap 0.5")
	}
	c := Pure("a").Compress(FractionFromInt(0), FractionFromInt(1))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0-1")
	}
	z := Pure("a").Zoom(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(z.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Zoom 0.25-0.75")
	}
}

func TestMJS_EveryWhen2(t *testing.T) {
	e := Pure("a").Every(3, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := e.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps) == 0 {
		t.Fatalf("Every 3 0-3")
	}
	w := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("When true 2")
	}
}

func TestMJS_IterBackChunk2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	it := p.Iter(2)
	_ = it
	ib := p.IterBack(2)
	if len(ib.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("IterBack 2")
	}
	ch := p.Chunk(2, func(pat Pattern) Pattern { return pat })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2")
	}
}
