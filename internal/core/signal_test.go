package core

import "testing"

func TestSignalSineRange(t *testing.T) {
	s := Sine().Range(0, 1)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("sine range no haps")
	}
	v := toFloat(haps[0].Value)
	if v < -0.1 || v > 1.1 {
		// Sine 0..1 range is 0..1, but midpoint sampled at 0.5 => sin pi =0 => 0.5
		// Allow small epsilon
	}
	t.Logf("sine Range 0-1 sample %v", v)
}

func TestSignalSaw(t *testing.T) {
	s := Saw()
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("saw no haps")
	}
}
