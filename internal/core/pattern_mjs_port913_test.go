package core

import "testing"

func TestMJS_Port913_StackWithSStackAmpersandFourth(t *testing.T) {
	p := Stack(Pure("bd"), Pure("sd"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Stack 2") }
	q := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 { t.Fatalf("FastCat 3") }
	r := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 { t.Fatalf("SlowCat 3") }
}
func TestMJS_Port913_SignalTriSawRandPerlinFourth(t *testing.T) {
	tri := Tri().Range(-1, 1)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Tri") }
	saw := Saw().Range(0, 10)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Saw") }
	r := Rand().Range(0, 5)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Rand") }
	per := Perlin().Range(-2, 2)
	if per.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("Perlin") }
}
func TestMJS_Port913_DegradeOftenRarelySometimesFourth(t *testing.T) {
	d0 := Pure("x").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 { t.Fatalf("DegradeBy 0") }
	d1 := Pure("y").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 { t.Fatalf("DegradeBy 1") }
	s := Pure("a").SometimesBy(0, func(q Pattern) Pattern { return q.Rev() })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("SometimesBy 0") }
	s2 := Pure("b").SometimesBy(1, func(q Pattern) Pattern { return q.Rev() })
	if s2.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil { t.Fatalf("SometimesBy 1") }
}
