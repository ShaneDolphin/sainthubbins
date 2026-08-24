package core

import "testing"

func TestGapWithSteps(t *testing.T) {
	p := Gap(3)
	if p.Steps == nil || !p.Steps.Equals(FractionFromInt(3)) {
		t.Fatalf("Gap 3 steps")
	}
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("Gap 3 should be silence")
	}
}

func TestSilencePWithLoc(t *testing.T) {
	p := Pure("a").SilenceP()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("SilenceP should be silence")
	}
}
