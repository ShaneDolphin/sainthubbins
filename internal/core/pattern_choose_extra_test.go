package core

import "testing"

func TestPickBasic(t *testing.T) {
	p := Pick("a", Pure(1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Pick expected >=1")
	}
}

func TestPickModBasic(t *testing.T) {
	p := PickMod("a", Pure(1))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("PickMod expected >=1")
	}
}

func TestSqueezeJoinBasic(t *testing.T) {
	p := Pure(Pure("a")).SqueezeJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("SqueezeJoin expected >=1")
	}
}

func TestInnerJoinBasic(t *testing.T) {
	p := Pure(Pure("a")).InnerJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("InnerJoin expected >=1")
	}
}
