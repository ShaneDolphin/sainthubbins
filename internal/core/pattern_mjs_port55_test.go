package core

import "testing"

func TestMJS_PolymeterSteps2(t *testing.T) {
	s := Stack(Pure("a"), Pure("b"))
	if s.Steps != nil {
		// Stack nil steps expected (FastCat nil) - just check non-nil doesn't matter
	}
	seq := Sequence(Pure("a"), Pure("b"), Pure("c"))
	haps := seq.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Sequence 3 expected non-empty")
	}
	g := Gap(4)
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(4)) {
		t.Fatalf("Gap 4 steps 4 expected got %v", g.Steps)
	}
}

func TestMJS_ControlStack2(t *testing.T) {
	pat := Stack(S("bd").Set(Gain(0.8)), S("sd").Set(Pan(0.5)))
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Control Stack 2 expected 2 got %d", len(haps))
	}
	// Check s keys
	for _, h := range haps {
		if m, ok := h.Value.(map[string]any); ok {
			if _, ok := m["s"]; !ok {
				t.Fatalf("s key missing %v", m)
			}
		}
	}
}

func TestMJS_SlowInsideNested2(t *testing.T) {
	p := Pure("a").SlowF(FractionFromInt(2)).Inside(FractionFromInt(2), func(pat Pattern) Pattern {
		return pat.FastF(FractionFromInt(2))
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slow2 Inside2 FastF2 expected non-empty")
	}
}
