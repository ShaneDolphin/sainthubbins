package core

import "testing"

func TestMJS_Port199_SignalPerlinRandChoiceFourth(t *testing.T) {
	s := Sine().Range(0, 1).Slow(FractionFromInt(4))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range Slow 4")
	}
	p := Perlin().Range(0, 10)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range 0,10")
	}
	r := Rand().Range(-5, 5)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range -5,5")
	}
	ch := Pure(1).Choose([]any{"a", "b"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b")
	}
}

func TestMJS_Port199_ArpWithChunkValuesFourth(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("updown")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
	ch := arp.Chunk(2, func(q Pattern) Pattern { return q.Rev() })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp Chunk 2 Rev")
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Arp("diverge")
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp diverge seq")
	}
}

func TestMJS_Port199_PatternValuesStackCatFourth(t *testing.T) {
	p := Stack(Pure("bd"), Pure("sd"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	q := Pure("hello").Fmap(func(v any) any { return v.(string) + v.(string) })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "hellohello" {
		t.Fatalf("hellohello")
	}
	r := Stack(S("bd:1"), S("sd:2"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack bd:1 sd:2 2")
	}
}
