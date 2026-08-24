package core

import "testing"

func TestMJS_Port145_SuperimposeLayerPalBrakHurry(t *testing.T) {
	s := Pure("bd").Superimpose(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Fast 2")
	}
	l := Stack(Pure("a"), Pure("b")).Layer(func(p Pattern) Pattern { return p.Rev() }, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Layer Rev+Fast")
	}
	pal := Pure("a b c").Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	brak := Pure("a b c d").Brak()
	if len(brak.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
	if len(Pure("bd").Hurry(2).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 2")
	}
}

func TestMJS_Port145_FilterHapsMaskStructArange(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b"), Pure("a")).FilterValues(func(v any) bool { return v == "a" })
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FilterValues a 2 got %d", len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	m := Pure("bd").Mask(FastCat(Pure(true), Pure(false)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tf")
	}
	st := Pure("bd").Struct(FastCat(Pure(true), Pure(true), Pure(false)))
	if st.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct ttf")
	}
	// Arrange-like via SlowCat
	arr := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(arr.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat Arrange 3")
	}
}

func TestMJS_Port145_HapContextWithValueStack(t *testing.T) {
	h := NewHap(nil, NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), "bd", map[string]any{"orbit": 1})
	// With nil whole, but value check
	if h.Value != "bd" {
		t.Fatalf("Hap bd")
	}
	p := Pure("sd").WithContext(func(m map[string]any) map[string]any { m["gain"] = 0.9; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["gain"] != 0.9 {
		t.Fatalf("gain 0.9")
	}
	q := S("hh").WithValue(func(v any) any {
		m := v.(map[string]any); m["pan"] = 0.3; return m
	})
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("S hh pan 0.3")
	}
}
