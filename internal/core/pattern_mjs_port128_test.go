package core

import "testing"

func TestMJS_Port128_InsideOutsideSignal(t *testing.T) {
	p := Sine().Inside(2, func(q Pattern) Pattern { return q.Range(0, 1) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Inside 2 Range")
	}
	q := Saw().Outside(2, func(pat Pattern) Pattern { return pat.Range(0, 10) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Outside 2 Range")
	}
	r := Pure("a b c").Inside(4, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Inside 4 Fast 2")
	}
}

func TestMJS_Port128_ArpWithMasksSignal(t *testing.T) {
	base := Sequence(Pure("c3"), Pure("e3"), Pure("g3"))
	arp := base.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	masked := arp.Mask(FastCat(Pure(true), Pure(false)))
	if masked.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Mask tf")
	}
	sig := Sine().Slow(FractionFromInt(4)).Range(0, 1)
	if sig.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 4 Range")
	}
}

func TestMJS_Port128_ChunkWithFastSlow(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2 Fast 2")
	}
	q := Pure("bd hh sd oh").Chunk(4, func(pat Pattern) Pattern { return pat.Rev() })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 4 Rev")
	}
	r := Pure("a b c d").Chunk(2, func(pat Pattern) Pattern { return pat.Slow(FractionFromInt(2)) })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2 Slow 2")
	}
}
