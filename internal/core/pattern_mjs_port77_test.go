package core

import "testing"

func TestMJS_FastCatSlowCatPrime3(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b"))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2 expected 2")
	}
	s := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps) == 0 {
		t.Fatalf("SlowCat 3 cycles")
	}
}

func TestMJS_ArpRev3(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	arp := p.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	rev := p.Rev()
	if len(rev.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Rev 3")
	}
}

func TestMJS_SequencePalindromePolymeter3(t *testing.T) {
	s := Sequence(Pure("a"), Pure("b"), Pure("c"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Sequence 3")
	}
	pal := s.Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat")
	}
}
