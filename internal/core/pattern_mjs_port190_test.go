package core

import "testing"

func TestMJS_Port190_SuperimposeLayerPalBrakFourth(t *testing.T) {
	p := Pure("a").Superimpose(func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Rev <2")
	}
	q := Stack(Pure("a"), Pure("b"), Pure("c"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack 3")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	s := Pure("a b").Brak()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
}

func TestMJS_Port190_FilterHapsMaskStructArangeFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FilterValues(func(v any) bool { return v == "a" })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("FilterValues a 1")
	}
	q := FastCat(Pure("a"), Pure("b")).Mask(Pure(false))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Mask false 0")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(true), Pure(false), Pure(true)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Struct 3")
	}
	s := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat 3")
	}
}

func TestMJS_Port190_HapContextWithValueStackFourth(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"gain": 0.9})
	if h.Context["gain"] != 0.9 {
		t.Fatalf("gain 0.9")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["pan"] = 0.5; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["pan"] != 0.5 {
		t.Fatalf("pan 0.5")
	}
	q := Pure(10).WithValue(func(v any) any { return v.(int) + 5 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 15 {
		t.Fatalf("WithValue 15")
	}
	r := Stack(S("bd").Set(Gain(0.5)), S("sd").Set(Pan(0.3)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
}
