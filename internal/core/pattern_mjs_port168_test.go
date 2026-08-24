package core

import "testing"

func TestMJS_Port168_PatternCatStackArpThird(t *testing.T) {
	c := Cat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
	s := Stack(Pure("x"), Pure("y"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	arp := Pure("c3 e3 g3 b3").Arp("converge")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp converge")
	}
}

func TestMJS_Port168_SignalRandChooseThird(t *testing.T) {
	r := Rand().Range(0, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,1")
	}
	ch := Pure(1).Choose([]any{"a", "b", "c"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b c")
	}
	p := Perlin().FastF(FractionFromInt(2)).Range(0, 5)
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Fast 2 Range 0,5")
	}
}

func TestMJS_Port168_ControlsNWithValueThird(t *testing.T) {
	n := N(3).FirstCycle()[0].Value.(map[string]any)
	if n["n"] != 3 {
		t.Fatalf("N 3")
	}
	s := S("bd:1*2").FirstCycle()
	if len(s) == 0 {
		t.Fatalf("S bd:1*2 empty")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["pan"] = -0.5; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["pan"] != -0.5 {
		t.Fatalf("pan -0.5")
	}
}
