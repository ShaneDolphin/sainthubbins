package core

import "testing"

func TestMJS_RandChooseDegrade(t *testing.T) {
	r := Rand()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rand expected non-empty")
	}
	ch := Pure(0).Choose([]any{"a", "b", "c"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose expected non-empty")
	}
	d := Pure("a").Degrade()
	if len(d.QueryArc(FractionFromInt(0), FractionFromInt(1))) > 1 {
		t.Fatalf("Degrade should be <=1")
	}
}

func TestMJS_SometimesOftenRarely(t *testing.T) {
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes expected non-empty")
	}
	o := Pure("a").SometimesBy(0.75, func(p Pattern) Pattern { return p })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.75 expected non-empty")
	}
}

func TestMJS_StackCatSequence(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack a,b expected 2 got %d", len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	cat := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(3))) == 0 {
		t.Fatalf("Cat 3 cycles expected non-empty")
	}
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Sequence expected non-empty")
	}
}
