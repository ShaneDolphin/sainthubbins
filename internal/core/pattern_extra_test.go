package core

import "testing"

func TestEuclid3_8(t *testing.T) {
	pat := Pure("x").Euclid(3, 8)
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Euclid 3,8 expected 3 haps, got %d", len(haps))
	}
}

func TestFastCatThree(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c"))
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("FastCat 3 expected 3, got %d", len(haps))
	}
}

func TestMiniBasic(t *testing.T) {
	// Mini is in internal/mini, but test core Reify stringParser hook
	SetStringParser(func(s string) Pattern {
		// Very simple: split spaces
		if s == "bd sd" {
			return FastCat(Pure("bd"), Pure("sd"))
		}
		return Pure(s)
	})
	pat := Reify("bd sd")
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Reify mini bd sd expected 2, got %d", len(haps))
	}
	SetStringParser(nil)
}

func TestAddRange(t *testing.T) {
	p := Pure(1.0).Range(0, 10)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Range no haps")
	}
	v := toFloat(haps[0].Value)
	if v < 0 || v > 10 {
		t.Fatalf("Range value out of range: %v", v)
	}
}
