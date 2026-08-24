package core

import "testing"

func TestMJS_SignalPerlinRand2(t *testing.T) {
	s := Sine().Slow(FractionFromInt(4))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Sine") }
	r := Rand().Segment(2)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Rand") }
}

func TestMJS_ArpWithChunk2(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up")
	c := p.Chunk(2, func(q Pattern) Pattern { return q.Rev() })
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Arp Chunk") }
}

func TestMJS_PatternValuesStack2(t *testing.T) {
	p := Stack(Pure(1), Pure(2), Pure(3))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)!=3 { t.Fatalf("Stack 3 got %d", len(haps)) }
}
