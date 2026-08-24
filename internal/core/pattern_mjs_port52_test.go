package core

import "testing"

func TestMJS_RevPalindrome(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	rev := p.Rev()
	haps := rev.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Rev 3 expected 3 got %d", len(haps))
	}
	pal := p.Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome expected non-empty")
	}
}

func TestMJS_JuxSuperimpose(t *testing.T) {
	j := Pure("a").Jux(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(j.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Jux expected non-empty")
	}
	sup := Pure("a").Superimpose(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(sup.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose expected >=2")
	}
}

func TestMJS_EuclidBjorklund(t *testing.T) {
	e := Pure("a").Euclid(3, 8)
	haps := e.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Euclid 3,8 expected 3 got %d", len(haps))
	}
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund 3,8 len 8 got %d", len(b))
	}
	count := 0
	for _, v := range b {
		if v != 0 {
			count++
		}
	}
	if count != 3 {
		t.Fatalf("Bjorklund 3,8 count 3 got %d", count)
	}
}
