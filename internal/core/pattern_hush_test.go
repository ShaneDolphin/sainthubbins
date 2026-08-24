package core

import "testing"

func TestHushBasic(t *testing.T) {
	p := Pure("a").Hush()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("Hush expected 0 got %d", len(haps))
	}
}

func TestSilenceBasic(t *testing.T) {
	p := Silence()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("Silence expected 0")
	}
}

func TestGapBasic(t *testing.T) {
	p := Pure("a").Gap(2)
	if p.Steps == nil || !p.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("Gap steps 2 expected")
	}
}
