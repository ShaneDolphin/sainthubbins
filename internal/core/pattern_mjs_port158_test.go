package core

import "testing"

func TestMJS_Port158_PatternCatStackArpSecond(t *testing.T) {
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
	s := Stack(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Stack 4")
	}
	arp := Pure("c3 e3 g3").Arp("down")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down")
	}
}

func TestMJS_Port158_SignalRandChooseSecond(t *testing.T) {
	r := Rand().FastF(FractionFromInt(2)).Range(0, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Fast 2 Range")
	}
	ch := Pure(0).Choose([]any{"a", "b"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b")
	}
	p := Perlin().Slow(FractionFromInt(2))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 2")
	}
}

func TestMJS_Port158_ControlsNWithValueSecond(t *testing.T) {
	n := N(5).FirstCycle()[0].Value.(map[string]any)
	if n["n"] != 5 {
		t.Fatalf("N 5")
	}
	p := S("bd:2").FirstCycle()
	if len(p) == 0 {
		t.Fatalf("S bd:2")
	}
	q := Pure(map[string]any{"s": "bd"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["delay"] = 0.5; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["delay"] != 0.5 {
		t.Fatalf("delay 0.5")
	}
}
