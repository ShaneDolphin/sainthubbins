package core

import "testing"

func TestFocusBasic(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Focus(0, 1)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Focus expected >=1")
	}
}

func TestFocusSpanBasic(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	p := Pure("a").FocusSpan(span)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("FocusSpan expected >=1")
	}
}
