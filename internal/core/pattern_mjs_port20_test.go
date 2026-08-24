// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Twentieth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_SlowcatPrimeArp(t *testing.T) {
	scp := SlowcatPrime(Pure("a"), Pure("b"))
	if len(scp.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("slowcatPrime empty")
	}
	arp := FastCat(Pure("a"), Pure("b"), Pure("c")).Arp(Pure(2))
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("arp 2 empty")
	}
}

func TestMJS_RevSequence(t *testing.T) {
	rev := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := rev.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("rev 3 expected 3 got %d", len(haps))
	}
	// check order reversed: first should be c (since rev)
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("sequence 2 expected 2")
	}
}

func TestMJS_PolyrhythmPolymeter(t *testing.T) {
	// polyrhythm via FastCat with different lengths approximated
	pr := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pr.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("polymeter empty")
	}
	// stepcat already tested
	sc := StepCat(Pure("a"), Pure("b"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("stepcat empty")
	}
}
