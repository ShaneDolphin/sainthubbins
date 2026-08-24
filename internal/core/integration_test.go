package core

import "testing"

func TestIntegration_StackFastCat(t *testing.T) {
	// Stack of two pures should give 2 haps per cycle
	p := Stack(Pure("bd"), Pure("sd"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("stack 2")
	}
	// FastCat 2 should interleave per half cycle
	fc := FastCat(Pure("a"), Pure("b"))
	haps2 := fc.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("fastcat 2")
	}
	if haps2[0].Part.Begin.Float64() != 0 || haps2[1].Part.Begin.Float64() != 0.5 {
		t.Fatalf("fastcat parts wrong %v %v", haps2[0].Part, haps2[1].Part)
	}
}

func TestIntegration_EarlyLate(t *testing.T) {
	p := Pure("x")
	early := p.EarlyF(MustParseFraction("1/4"))
	haps := early.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("early 1/4 expected 2 got %d", len(haps))
	}
	late := p.LateF(MustParseFraction("1/4"))
	haps2 := late.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("late 1/4 expected 2 got %d", len(haps2))
	}
}

func TestIntegration_SignalRange(t *testing.T) {
	s := Sine().Range(0, 1)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("signal range no haps")
	}
	v := toFloat(haps[0].Value)
	if v < 0 || v > 1 {
		t.Fatalf("signal range out of 0-1: %v", v)
	}
}

func TestIntegration_Pick(t *testing.T) {
	list := []any{"a", "b", "c"}
	pat := Pick(list, Pure(1))
	// Pick uses innerJoin, need to query
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("pick no haps")
	}
}
