package core

import "testing"

func TestMJS_Port137_SignalSineTriSawSlow(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2)).Range(-1, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 2 Range -1,1")
	}
	tri := Tri().FastF(FractionFromInt(2)).Range(0, 1)
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Fast 2")
	}
	saw := Saw().Range(10, 20).Slow(FractionFromInt(1))
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range 10,20 Slow 1")
	}
}

func TestMJS_Port137_PatternPureStackCatSlow(t *testing.T) {
	p := Pure("x").Fmap(func(v any) any { return v.(string) + v.(string) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != "xx" {
		t.Fatalf("Fmap xx")
	}
	s := Stack(Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	c := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat 3")
	}
}

func TestMJS_Port137_ChooseWithRandSegment(t *testing.T) {
	ch := Pure(1).Choose([]any{10, 20, 30})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose 10,20,30")
	}
	r := Rand()
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand")
	}
	rs := Rand().Segment(2)
	if rs.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Segment 2")
	}
}
