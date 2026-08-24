package core

import "testing"

func TestMJS_StructWithBool2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	s := p.Struct(Pure(true))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true expected non-empty")
	}
	sF := p.Struct(Pure(false))
	if len(sF.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Struct false expected 0")
	}
}

func TestMJS_MaskWithPattern2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	m := p.Mask(Pure(true))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true expected non-empty")
	}
	mF := p.Mask(Pure(false))
	if len(mF.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Mask false expected 0")
	}
}

func TestMJS_EuclidLegato3(t *testing.T) {
	e := Pure("a").EuclidLegato(3, 8)
	haps := e.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("EuclidLegato 3,8 3 got %d", len(haps))
	}
}
