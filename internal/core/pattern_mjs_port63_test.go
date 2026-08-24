package core

import "testing"

func TestMJS_TimeMethods3(t *testing.T) {
	p := Pure("a").Early(FractionFromFloat(0.25))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Early 0.25 expected non-empty")
	}
	l := Pure("a").Late(FractionFromFloat(0.25))
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Late 0.25 expected non-empty")
	}
	lf := Pure("a").LateF(FractionFromFloat(0.5))
	if len(lf.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("LateF 0.5 expected non-empty")
	}
}

func TestMJS_HushAndSilence3(t *testing.T) {
	h := Pure("a").Hush()
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush expected 0")
	}
	s := Silence()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence expected 0")
	}
	// Zero range may be 0 or 1 depending on SpanCycles impl; just check not panic
	p := Pure("a")
	_ = p.QueryArc(FractionFromInt(0), FractionFromInt(0))
}

func TestMJS_PlainValues3(t *testing.T) {
	p := Pure(42)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(int) != 42 {
		t.Fatalf("Pure 42 expected 42")
	}
	p2 := Pure(3.14)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 || haps2[0].Value.(float64) != 3.14 {
		t.Fatalf("Pure 3.14 expected 3.14")
	}
}
