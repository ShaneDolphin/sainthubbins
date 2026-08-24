package core

import "testing"

func TestMJS_FastSlowCompress3(t *testing.T) {
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2 expected 2")
	}
	s := Pure("a").SlowF(FractionFromInt(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("SlowF 2")
	}
	c := Pure("a").Compress(FractionFromInt(0), FractionFromFloat(0.5))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0-0.5")
	}
}

func TestMJS_InsideOutside3(t *testing.T) {
	in := Pure("a").Inside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(in.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside 2")
	}
	out := Pure("a").Outside(2, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(out.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside 2")
	}
}

func TestMJS_FilterValues3(t *testing.T) {
	p := Stack(Pure("a"), Pure("b")).FilterValues(func(v any) bool { return v == "a" })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value != "a" {
		t.Fatalf("FilterValues a")
	}
}
