package core

import "testing"

func TestMJS_Port118_StackWithSStackAmpersand(t *testing.T) {
	s := Stack(S("bd"), S("sd hh"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack S bd + sd hh 2 got %d", len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	c := Cat(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 3 -> SlowCat 1")
	}
	fc := FastCat(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3")
	}
	sc := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat 3 cycles 3")
	}
}

func TestMJS_Port118_SignalTriSawRandPerlin(t *testing.T) {
	tri := Tri().Range(-1, 1)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -1,1")
	}
	saw := Saw().Range(0, 10)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 0,10")
	}
	rand := Rand().Range(0, 1)
	if rand.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range")
	}
	perlin := Perlin().Range(0, 5).Slow(FractionFromInt(2))
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Range Slow")
	}
}

func TestMJS_Port118_DegradeOftenRarelySometimes(t *testing.T) {
	d0 := Pure("bd").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0")
	}
	d1 := Pure("bd").DegradeBy(1)
	if len(d1.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("DegradeBy 1 empty")
	}
	s := Pure("bd").SometimesBy(0, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("SometimesBy 0 ->1")
	}
	o := Pure("bd").SometimesBy(1, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("SometimesBy 1 ->2")
	}
}
