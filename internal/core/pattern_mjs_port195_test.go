package core

import "testing"

func TestMJS_Port195_SuperimposeLayerPalBrakHurryFourth(t *testing.T) {
	p := Pure("a").Superimpose(func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Rev <2")
	}
	q := Stack(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("Stack 4")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	s := Pure("a b").Brak()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
	h := Pure("a b c").Hurry(2)
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 2")
	}
}

func TestMJS_Port195_FilterHapsMaskStructArangeFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FilterValues(func(v any) bool { return v == "a" || v == "c" })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FilterValues a c 2")
	}
	q := FastCat(Pure("a"), Pure("b")).Mask(Pure(true))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mask true 2")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Struct 2")
	}
	s := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat 3")
	}
}

func TestMJS_Port195_HapContextWithValueStackFourth(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 2})
	if h.Context["orbit"] != 2 {
		t.Fatalf("orbit 2")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["delay"] = 0.5; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["delay"] != 0.5 {
		t.Fatalf("delay 0.5")
	}
	q := Pure(5).WithValue(func(v any) any { return v.(int) * 4 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 20 {
		t.Fatalf("WithValue 20")
	}
}
