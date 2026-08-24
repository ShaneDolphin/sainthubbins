package core

import "testing"

func TestMJS_ArpWithSlowInside2(t *testing.T) {
	p := Pure("c3 e3 g3").Slow(FractionFromInt(2))
	a := p.Arp("up")
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Arp up slow") }
}

func TestMJS_ChunkWithFastSlow2(t *testing.T) {
	p := Pure("bd sd").Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Chunk Fast") }
	s := Pure("a b c d").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Slow") }
}

func TestMJS_PatternAddMulStructure2(t *testing.T) {
	p := FastCat(Pure(1), Pure(2))
	mul := p.Mul(Pure(2))
	haps := mul.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)!=2 { t.Fatalf("Mul 2 got %d", len(haps)) }
}
