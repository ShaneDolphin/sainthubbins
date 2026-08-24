package core

import "testing"

func TestMJS_SlowcatPrimeArp3(t *testing.T) {
	s := SlowCat(Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("SlowCat 2 cycles")
	}
	// SlowcatPrime via SlowCat alias (same)
	arp := Sequence(Pure("a"), Pure("b"), Pure("c")).Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
}

func TestMJS_RevSequence3(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	rev := p.Rev()
	haps := rev.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Rev 3 expected 3 got %d", len(haps))
	}
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Sequence 2")
	}
}

func TestMJS_PolyrhythmPolymeter3(t *testing.T) {
	pm := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat")
	}
	sc := StepCat(Pure("a"), Pure("b"), Pure("c"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("StepCat 3")
	}
}
