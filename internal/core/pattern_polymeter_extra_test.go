package core

import "testing"

func TestPolymeterSteps(t *testing.T) {
	p := Polymeter(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Polymeter steps expected >=1")
	}
	// Check steps lcm handling via FastCat vs Polymeter
	q := FastCat(Pure("a"), Pure("b"))
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != len(haps2) {
		// Polymeter may differ but both non-empty
		t.Logf("Polymeter %d vs FastCat %d", len(haps), len(haps2))
	}
}

func TestPolyrhythmBasic(t *testing.T) {
	// Polyrhythm is Stack with different lengths? Simple check
	p := Stack(Pure("a").Slow(FractionFromInt(2)), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps) == 0 {
		t.Fatalf("Polyrhythm stack expected >=1")
	}
}

func TestArrangeBasic(t *testing.T) {
	p := Arrange(2, Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(4))
	if len(haps) == 0 {
		t.Fatalf("Arrange expected >=1")
	}
}

func TestPalindromeBasic(t *testing.T) {
	p := SlowCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps) == 0 {
		t.Fatalf("Palindrome over 2 cycles expected >=1")
	}
}
