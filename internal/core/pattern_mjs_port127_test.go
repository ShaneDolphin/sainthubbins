package core

import "testing"

func TestMJS_Port127_PatternValueWithContext(t *testing.T) {
	p := Pure("bd").WithValue(func(v any) any { return v })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("WithValue bd")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.9; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.9 {
		t.Fatalf("WithValue gain 0.9")
	}
	r := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["orbit"] = 2; return m })
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["orbit"] != 2 {
		t.Fatalf("WithContext orbit 2")
	}
}

func TestMJS_Port127_SlowFastChoiceWithValue(t *testing.T) {
	s := Pure("a b").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2")
	}
	ch := Pure(0).Choose([]any{Pure("a"), Pure("b")})
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose Pattern a/b")
	}
}

func TestMJS_Port127_StackCatPolymeterSteps(t *testing.T) {
	st := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	cat := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("Cat 3 cycles 3")
	}
	pm := PolymeterSlowcat(Pure("x"), Pure("y"))
	if pm.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("PolymeterSlowcat 2")
	}
}
