package core

import "testing"

func TestMJS_Port165_SignalPerlinRandChoiceThird(t *testing.T) {
	s := Sine().Range(0, 10)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,10")
	}
	p := Perlin().Slow(FractionFromInt(3)).Range(0, 5)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 3 Range 0,5")
	}
	r := Rand().Range(-1, 1).FastF(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range Fast 2")
	}
}

func TestMJS_Port165_ArpWithChunkValuesThird(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	ch := arp.Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Chunk 2")
	}
	seq := Sequence(Pure("c3"), Pure("e3"), Pure("g3"), Pure("b3")).Arp("down")
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down seq")
	}
}

func TestMJS_Port165_PatternValuesStackCatThird(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"), Pure("c"), Pure("d"), Pure("e"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 5 {
		t.Fatalf("Stack 5")
	}
	q := Pure(5).Fmap(func(v any) any { return v.(int) * 2 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 10 {
		t.Fatalf("Fmap 10")
	}
	r := FastCat(Pure("x"), Pure("y")).WithValue(func(v any) any { return v })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat WithValue 2")
	}
}
