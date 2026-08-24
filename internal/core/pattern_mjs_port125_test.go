package core

import "testing"

func TestMJS_Port125_SuperimposeLayerPalBrak(t *testing.T) {
	s := Pure("bd sd").Superimpose(func(q Pattern) Pattern { return q.Rev() })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Rev")
	}
	l := Stack(Pure("bd"), Pure("sd")).Layer(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) }, func(p Pattern) Pattern { return p.Rev() })
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Layer Fast+Rev")
	}
	pal := Pure("a b c d").Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	if len(Pure("a b c").Brak().QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
}

func TestMJS_Port125_FilterHapsMaskStruct(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b"), Pure("a")).FilterValues(func(v any) bool { return v == "a" })
	haps := f.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("FilterValues a")
	}
	// Ensure filtered doesn't contain b
	for _, h := range haps {
		if h.Value == "b" {
			t.Fatalf("b should be filtered")
		}
	}
	m := Pure("a").Mask(FastCat(Pure(true), Pure(false)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tf")
	}
	st := Pure("a").Struct(FastCat(Pure(true), Pure(true), Pure(false)))
	if st.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct ttf")
	}
}

func TestMJS_Port125_HapContextWithValue(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "bd", map[string]any{"orbit": 3})
	if h.Context["orbit"] != 3 {
		t.Fatalf("orbit 3")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["delay"] = 0.5; return m })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Context["delay"] != 0.5 {
		t.Fatalf("WithContext delay 0.5")
	}
	q := Pure(42).WithValue(func(v any) any { return v.(int) * 2 })
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 84 {
		t.Fatalf("WithValue 84")
	}
}
