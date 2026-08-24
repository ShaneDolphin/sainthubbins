package core

import "testing"

func TestMJS_InsideOutsideWithSignal2(t *testing.T) {
	s := Sine().Range(0, 1)
	p := s.Inside(FractionFromInt(2), func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside2 with signal expected non-empty")
	}
	q := s.Outside(FractionFromInt(2), func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside2 with signal expected non-empty")
	}
}

func TestMJS_ArpWithMasks2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	arp := p.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up expected non-empty")
	}
	// Mask true vs false
	mT := p.Mask(Pure(true))
	if len(mT.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true expected non-empty")
	}
	mF := p.Mask(Pure(false))
	if len(mF.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Mask false expected 0")
	}
}

func TestMJS_ChunkWithFast2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	ch := p.Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2 FastF2 expected non-empty")
	}
}
