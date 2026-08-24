package core

import "testing"

func TestMJS_InsideOutsideWithSignal3(t *testing.T) {
	s := Sine().Range(0, 1)
	p := s.Inside(FractionFromInt(2), func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside2 with signal")
	}
	q := s.Outside(FractionFromInt(2), func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside2 with signal")
	}
}

func TestMJS_ArpWithMasks3(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	arp := p.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	mT := p.Mask(Pure(true))
	if len(mT.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true")
	}
}

func TestMJS_ChunkWithFast3(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	ch := p.Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2 FastF2")
	}
}
