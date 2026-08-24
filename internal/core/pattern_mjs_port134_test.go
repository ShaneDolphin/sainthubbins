package core

import "testing"

func TestMJS_Port134_SignalPerlinRand(t *testing.T) {
	s := Sine().Range(0, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,1")
	}
	p := Perlin().Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 2")
	}
	r := Rand().FastF(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Fast 2")
	}
}

func TestMJS_Port134_ArpWithChunk(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up")
	ch := p.Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Chunk 2 Fast 2")
	}
	q := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Arp("down")
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down Seq")
	}
}

func TestMJS_Port134_PatternValuesStack(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Stack 4")
	}
	q := Pure(42).Fmap(func(v any) any { return v.(int) * 2 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 84 {
		t.Fatalf("Fmap 84")
	}
	r := Stack(FastCat(Pure("x"), Pure("y")), Pure("a"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastCat Stack")
	}
}
