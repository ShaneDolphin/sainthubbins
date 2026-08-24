package core

import "testing"

func TestMJS_Port800_StackCatFastSlowFourth(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Stack 2") }
	q := FastCat(Pure("x"), Pure("y"), Pure("z"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 { t.Fatalf("FastCat 3") }
	f := Pure("c").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("FastF") }
	s := Pure("d").Slow(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 { t.Fatalf("Slow") }
}
func TestMJS_Port800_SignalRangePerlinFourth(t *testing.T) {
	s := Sine().Range(-1, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Sine") }
	per := Perlin().Range(0, 1).Slow(FractionFromInt(2))
	if per.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Perlin") }
	saw := Saw().Range(0, 100)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Saw") }
}
func TestMJS_Port800_ArpChunkDegradeFourth(t *testing.T) {
	arp := Pure("c3 e3 g3").Arp("converge")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Arp") }
	ch := Pure("a b c d").Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Chunk") }
	d0 := Pure("x").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 { t.Fatalf("DegradeBy 0") }
	d1 := Pure("y").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 { t.Fatalf("DegradeBy 1") }
}
