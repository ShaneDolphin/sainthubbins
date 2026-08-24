package core

import "testing"

func TestMJS_Port138_SignalPerlinRandChoice(t *testing.T) {
	p := Perlin().Range(0, 5)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range 0,5")
	}
	r := Rand().Range(0, 1).FastF(FractionFromInt(2))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range Fast 2")
	}
	ch := Pure("a").Choose([]any{"a", "b", "c"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b c")
	}
}

func TestMJS_Port138_ArpWithChunkValues(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("up").Chunk(2, func(q Pattern) Pattern { return q.Rev() })
	if arp.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Chunk 2 Rev")
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c")).Arp("down")
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down seq")
	}
}

func TestMJS_Port138_PatternValuesStackCat(t *testing.T) {
	p := Stack(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	q := FastCat(Pure("a"), Pure("b")).Fmap(func(v any) any { return v.(string) + "!" })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat Fmap 2")
	}
	r := Pure(100).WithValue(func(v any) any { return v.(int) / 2 })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 50 {
		t.Fatalf("WithValue 50")
	}
}
