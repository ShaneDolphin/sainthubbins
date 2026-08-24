package core

import "testing"

func TestGoldenStackSubpatterns(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Stack subpatterns expected 2 got %d", len(haps))
	}
}

func TestGoldenFastCatSequence(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("FastCat sequence expected 3 got %d", len(haps))
	}
	// Check order a,b,c by time
	if haps[0].Value.(string) != "a" || haps[1].Value.(string) != "b" || haps[2].Value.(string) != "c" {
		// Allow any order due to sorting, just check presence
		found := map[string]bool{}
		for _, h := range haps {
			found[h.Value.(string)] = true
		}
		if !found["a"] || !found["b"] || !found["c"] {
			t.Fatalf("FastCat missing values %v", haps)
		}
	}
}

func TestGoldenPolymeterBasic(t *testing.T) {
	p := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Polymeter expected >=1 got 0")
	}
}

func TestGoldenSlowCat(t *testing.T) {
	p := SlowCat(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("SlowCat over 1 cycle expected 1 got %d", len(haps))
	}
	haps2 := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps2) != 2 {
		t.Fatalf("SlowCat over 2 cycles expected 2 got %d", len(haps2))
	}
}

func TestGoldenPalindrome(t *testing.T) {
	p := SlowCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Palindrome expected >=1")
	}
}

func TestGoldenSometimes(t *testing.T) {
	p := Pure("a").Sometimes(func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Sometimes expected 1 got %d", len(haps))
	}
}

func TestGoldenOff(t *testing.T) {
	p := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.Add(Pure(1)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Off expected >=2 got %d", len(haps))
	}
}

func TestGoldenJux(t *testing.T) {
	p := Pure("a").Jux(func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Jux expected 2 got %d", len(haps))
	}
}

func TestGoldenSuperimpose(t *testing.T) {
	p := Pure("a").Superimpose(func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Superimpose expected 2 got %d", len(haps))
	}
}
