package core

import "testing"

func TestMJS_ChordVoicing2(t *testing.T) {
	p := Pure("a").Off(0.25, func(pat Pattern) Pattern { return pat })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 expected >=2")
	}
	// Scale
	s := Pure("a").Scale("minor")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor expected non-empty")
	}
}

func TestMJS_ScaleQuantize2(t *testing.T) {
	p := Pure(2).Scale("major")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale major with 2 expected non-empty")
	}
}

func TestMJS_EuclidWithOff2(t *testing.T) {
	e := Pure("a").Euclid(3, 8)
	haps := e.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Euclid 3,8 3")
	}
	eOff := Pure("a").Euclid(3, 8).Off(0.25, func(p Pattern) Pattern { return p })
	if len(eOff.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid Off expected non-empty")
	}
}
