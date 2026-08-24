package core

import "testing"

func TestMJS_Port153_SignalPerlinRandChoiceSecond(t *testing.T) {
	s := Sine().Range(-1, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range -1,1")
	}
	p := Perlin().Range(0, 1)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range 0,1")
	}
	r := Rand().Range(0, 100)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,100")
	}
	ch := Pure(2).Choose([]any{"x", "y"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose x y")
	}
}

func TestMJS_Port153_ArpWithChunkValuesSecond(t *testing.T) {
	arp := Pure("c3 e3 g3 b3").Arp("updown")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown 4")
	}
	ch := Pure("bd sd hh").Chunk(2, func(q Pattern) Pattern { return q.Rev() })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2 Rev")
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Arp("converge")
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge seq")
	}
}

func TestMJS_Port153_PatternValuesStackCatSecond(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	q := Pure("test").Fmap(func(v any) any { return v.(string) + "!" })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "test!" {
		t.Fatalf("Fmap test!")
	}
	r := Stack(S("bd"), Gain(0.8))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack S bd Gain")
	}
}
