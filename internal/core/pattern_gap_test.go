package core

import "testing"

func TestGapSilence(t *testing.T) {
	g := Pure("a").Gap(2)
	if len(g.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Gap should be silence")
	}
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("Gap steps expected 2 got %v", g.Steps)
	}
	s := Pure("a").SilenceP()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("SilenceP should be silence")
	}
}
