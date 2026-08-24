package core

import "testing"

func TestMJS_SignalRange3(t *testing.T) {
	s := Signal(func(frac Fraction) float64 { return frac.Float64() }).Range(0, 1)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Signal Range expected non-empty")
	}
	// Sine range 0-1
	si := Sine().Range(0, 1)
	haps2 := si.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Sine Range expected non-empty")
	}
	// Saw range
	sw := Saw().Range(10, 20)
	haps3 := sw.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) == 0 {
		t.Fatalf("Saw Range expected non-empty")
	}
}

func TestMJS_HapState3(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "a", map[string]any{"orbit": 1})
	if h.Value != "a" {
		t.Fatalf("Hap value a")
	}
	state := NewState(span, map[string]any{"cps": 0.5})
	if state.Span.Begin.Float64() != 0 {
		t.Fatalf("State span 0")
	}
}

func TestMJS_AddWithStructure3(t *testing.T) {
	// FastCat with Add pattern
	p := FastCat(Pure(1), Pure(2)).Fmap(func(v any) any { return v.(int) + 10 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("FastCat Add expected 2 got %d", len(haps))
	}
	if haps[0].Value.(int) != 11 || haps[1].Value.(int) != 12 {
		t.Fatalf("FastCat Add values 11,12 got %v %v", haps[0].Value, haps[1].Value)
	}
}
