package core

import "testing"

func TestMJS_Port187_StackCatSequencePolymeterFourth(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	q := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat 3")
	}
	r := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Sequence 4")
	}
	s := PolymeterSlowcat(Pure("a b"), Pure("c d e"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat 2")
	}
}

func TestMJS_Port187_SignalSineRandChooseFourth(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2)).Range(0, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow2 Range 0,1")
	}
	r := Rand().Range(0, 10)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Range 0,10")
	}
	ch := Pure(1).Choose([]any{"a", "b", "c", "d"})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose a b c d")
	}
}

func TestMJS_Port187_ControlsNWithValueFourth(t *testing.T) {
	p := S("bd")
	v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["s"] != "bd" {
		t.Fatalf("s bd got %v", v["s"])
	}
	q := S("bd").Set(N(2))
	v2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v2["n"] != 2 {
		t.Fatalf("n 2 got %v", v2["n"])
	}
	r := Pure(5).WithValue(func(v any) any { return v.(int) * 2 })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 10 {
		t.Fatalf("WithValue 10")
	}
	s := S("bd").Set(Gain(0.5)).Set(Pan(0.2))
	vm := s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if vm["gain"] != 0.5 || vm["pan"] != 0.2 {
		t.Fatalf("gain pan got %v", vm)
	}
}
