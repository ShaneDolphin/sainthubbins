package core

import "testing"

func TestStateSpan(t *testing.T) {
	s := NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil)
	if !s.Span.Begin.Equals(FractionFromInt(0)) {
		t.Fatalf("State Span Begin")
	}
	s2 := s.SetSpan(NewTimeSpan(FractionFromInt(1), FractionFromInt(2)))
	if !s2.Span.Begin.Equals(FractionFromInt(1)) {
		t.Fatalf("State SetSpan")
	}
}

func TestStateControls(t *testing.T) {
	s := NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), map[string]any{"gain": 0.5})
	if s.Controls["gain"] != 0.5 {
		t.Fatalf("State controls gain")
	}
}

func TestPatternFirstCycle(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"))
	first := p.FirstCycle()
	if len(first) != 3 {
		t.Fatalf("FirstCycle expected 3 got %d", len(first))
	}
}

func TestPatternWithValue2(t *testing.T) {
	p := Pure(1).WithValue(func(v any) any { return v.(int) * 2 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value.(int) != 2 {
		t.Fatalf("WithValue expected 2")
	}
}
