package core

import "testing"

func TestDefragmentBasic(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"))
	q := p.Defragment()
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Defragment expected >=1")
	}
}

func TestResetJoin(t *testing.T) {
	p := Pure(Pure("a")).ResetJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("ResetJoin expected >=1")
	}
}

func TestPolyJoin(t *testing.T) {
	p := Pure(Pure("a")).PolyJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("PolyJoin expected >=1")
	}
}
