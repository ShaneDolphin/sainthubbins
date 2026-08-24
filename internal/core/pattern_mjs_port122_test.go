package core

import "testing"

func TestMJS_Port122_StackCatSequencePolymeter(t *testing.T) {
	s := Stack(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	cat := Cat(Pure("a"), Pure("b"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat SlowCat 1")
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Sequence 4")
	}
	pm := PolymeterSlowcat(Pure("bd"), Pure("sd"))
	if pm.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("PolymeterSlowcat 2")
	}
}

func TestMJS_Port122_SignalSineRandChoose(t *testing.T) {
	s := Sine().Range(0, 1).Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 2")
	}
	r := Rand().Range(-1, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range -1,1")
	}
	ch := Pure(0).Choose([]any{"a", "b", "c", "d"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose 4")
	}
}

func TestMJS_Port122_ControlsNWithValue(t *testing.T) {
	s := S("bd:1")
	// bd:1 parsed via mini? S("bd:1") may be s with n? Check map
	h := s.FirstCycle()
	if len(h) == 0 {
		t.Fatalf("S bd:1 empty")
	}
	n := N(2).FirstCycle()[0].Value.(map[string]any)
	if n["n"] != 2 {
		t.Fatalf("N 2")
	}
	p := Pure(map[string]any{"s": "bd"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["orbit"] = 2; return m
	})
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["orbit"] != 2 {
		t.Fatalf("WithValue orbit 2")
	}
}
