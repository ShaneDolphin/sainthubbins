package core

import "testing"

func TestMJS_Port141_StackCatPolymeterSequence(t *testing.T) {
	s := Stack(Pure("bd"), Pure("sd"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat 1")
	}
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Sequence 2")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"), Pure("c"))
	if pm.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("PolymeterSlowcat 3")
	}
}

func TestMJS_Port141_SlowFastChoiceWithValue(t *testing.T) {
	s := Pure("a b").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("x").FastF(FractionFromInt(4))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastF 4")
	}
	ch := Pure(0).Choose([]any{"a", "b"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b")
	}
}

func TestMJS_Port141_InsideOutsideChunk(t *testing.T) {
	p := Pure("bd sd hh").Inside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Inside 2 Fast 2")
	}
	q := Pure("a b c d").Outside(2, func(pat Pattern) Pattern { return pat.Rev() })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Outside 2 Rev")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2 Fast 2")
	}
}
