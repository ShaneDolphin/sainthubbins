package core

import "testing"

func TestMJS_PatternCatStackArp2(t *testing.T) {
	p := Cat(Pure("a"), Pure("b"))
	s := Stack(p, p.FastF(FractionFromInt(2)))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) < 2 { t.Fatalf("Cat Stack") }
	a := Pure("c3 e3 g3").Arp("down")
	if a.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Arp down") }
}

func TestMJS_SignalRandChoose2(t *testing.T) {
	r := Rand().Segment(4)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Rand") }
	ch := Pure(0).Choose([]any{"a", "b", "c"})
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Choose") }
}

func TestMJS_ControlsNWithValue2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "n": 0})
	q := p.WithValue(func(v any) any {
		m := v.(map[string]any)
		m["n"] = 3
		return m
	})
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 || haps[0].Value.(map[string]any)["n"] != 3 { t.Fatalf("n 3") }
}
