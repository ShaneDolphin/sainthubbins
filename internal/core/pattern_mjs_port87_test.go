package core

import "testing"

func TestMJS_CatPrime2_2(t *testing.T) {
	s := SlowCat(Pure("a"), Pure("b"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("SlowCat 2")
	}
	f := FastCat(Pure("a"), Pure("b"))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastCat 2")
	}
}

func TestMJS_RevBrak2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	rev := p.Rev()
	if len(rev.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Rev 3")
	}
	b := Pure("a").Brak()
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
}

func TestMJS_PalindromePolyrhythm2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	pal := p.Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	pm := PolymeterSlowcat(Pure("a"), Pure("b"))
	if len(pm.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolymeterSlowcat")
	}
}
