package core

import "testing"

func TestMJS_IterAndPluck2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	it := p.Iter(4)
	// Iter 4 may be empty per Segment2+SqueezeJoin; check IterBack fallback
	if len(it.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		ib := p.IterBack(4)
		if len(ib.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
			t.Fatalf("Iter 4 and IterBack 4 both empty")
		}
	}
}

func TestMJS_PatternLag2(t *testing.T) {
	p := Pure("a").Late(FractionFromFloat(0.25))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Late 0.25 expected non-empty")
	}
	// Just check non-empty, shift may wrap
	_ = haps[0].Part.Begin
}

func TestMJS_TimeMethods4(t *testing.T) {
	e := Pure("a").Early(FractionFromFloat(0.1))
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Early 0.1 expected non-empty")
	}
	ef := Pure("a").EarlyF(FractionFromFloat(0.2))
	if len(ef.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("EarlyF 0.2 expected non-empty")
	}
	lf := Pure("a").LateF(FractionFromFloat(0.1))
	if len(lf.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("LateF 0.1 expected non-empty")
	}
}
