package core

import "testing"

func TestMJS_Port150_SuperimposeLayerPalindBrak(t *testing.T) {
	s := Pure("bd").Superimpose(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose Fast 2")
	}
	l := Stack(Pure("a"), Pure("b")).Layer(func(p Pattern) Pattern { return p.Rev() })
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Layer Rev")
	}
	pal := Pure("a b c d").Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
	brak := Pure("bd sd hh oh").Brak()
	if len(brak.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak 4")
	}
}

func TestMJS_Port150_FilterHapsMaskStructPort(t *testing.T) {
	f := FastCat(Pure("a"), Pure("b"), Pure("c")).FilterValues(func(v any) bool { return v != "b" })
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FilterValues !=b 2")
	}
	m := Pure("a").Mask(FastCat(Pure(true), Pure(false), Pure(true)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tf t")
	}
	st := Pure("bd").Struct(FastCat(Pure(true), Pure(false)))
	if st.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf")
	}
}

func TestMJS_Port150_HapContextWithValuePort(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&span, span, "sd", map[string]any{"n": 5})
	if h.Context["n"] != 5 {
		t.Fatalf("n 5")
	}
	p := Pure("bd").WithContext(func(m map[string]any) map[string]any { m["orbit"] = 7; return m })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Context["orbit"] != 7 {
		t.Fatalf("orbit 7")
	}
	q := Pure(map[string]any{"s": "hh"}).WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.8; return m
	})
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["gain"] != 0.8 {
		t.Fatalf("gain 0.8")
	}
}
