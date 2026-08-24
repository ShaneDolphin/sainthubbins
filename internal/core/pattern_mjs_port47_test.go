package core

import "testing"

func TestMJS_BeatCeilFloor(t *testing.T) {
	p := Pure(2.3).Ceil()
	haps := p.FirstCycle()
	if len(haps) == 0 || toFloat(haps[0].Value) != 3 {
		t.Fatalf("Ceil 2.3 expected 3 got %v", haps)
	}
	p2 := Pure(2.7).Floor()
	haps2 := p2.FirstCycle()
	if len(haps2) == 0 || toFloat(haps2[0].Value) != 2 {
		t.Fatalf("Floor 2.7 expected 2 got %v", haps2)
	}
	// Beat should not be empty
	b := Pure("a").Beat(0, 4)
	haps3 := b.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) == 0 {
		t.Fatalf("Beat expected non-empty")
	}
}

func TestMJS_ChunkIntoLoopAtCps(t *testing.T) {
	// ChunkInto delegate to Chunk
	c := Pure("a").ChunkInto(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ChunkInto expected non-empty")
	}
	// chunkinto lower alias
	c2 := Pure("a").Chunkinto(2, func(p Pattern) Pattern { return p })
	if len(c2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunkinto expected non-empty")
	}
	cb := Pure("a").ChunkBackInto(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(cb.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ChunkBackInto expected non-empty")
	}
	// LoopAtCps delegate to LoopAt
	l := Pure("a").LoopAtCps(2, 0.5)
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("LoopAtCps expected non-empty")
	}
	l2 := Pure("a").Loopatcps(2, 0.5)
	if len(l2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Loopatcps expected non-empty")
	}
}

func TestMJS_FluxZoomArcStutWithApplyN(t *testing.T) {
	// Flux / JuxFlip
	f := Pure("a").Flux(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Flux expected non-empty")
	}
	jf := Pure("a").JuxFlip(func(p Pattern) Pattern { return p })
	if len(jf.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("JuxFlip expected non-empty")
	}
	jfb := Pure("a").JuxFlipBy(0.5, func(p Pattern) Pattern { return p })
	if len(jfb.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("JuxFlipBy expected non-empty")
	}
	fb := Pure("a").FluxBy(0.5, func(p Pattern) Pattern { return p })
	if len(fb.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FluxBy expected non-empty")
	}
	// ZoomArc
	za := Pure("a").ZoomArc(NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5)))
	if len(za.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ZoomArc expected non-empty")
	}
	// StutWith / Echo alias
	s := Pure("a").StutWith(2, FractionFromFloat(0.25), func(p Pattern, n int) Pattern { return p })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("StutWith expected non-empty")
	}
	// ApplyN
	an := Pure(1).ApplyN(3, func(p Pattern) Pattern {
		return p.Fmap(func(v any) any { return v.(int) + 1 })
	})
	haps := an.FirstCycle()
	if len(haps) == 0 || haps[0].Value.(int) != 4 {
		t.Fatalf("ApplyN 1+3 expected 4 got %v", haps)
	}
	// Reset / Restart
	r := Pure("a").Reset(Pure(1))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Reset expected non-empty")
	}
}
