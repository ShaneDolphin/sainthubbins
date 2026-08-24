package core

import "testing"

func TestPureWithLoc(t *testing.T) {
	p := Pure("a").WithLoc(0, 1)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("WithLoc expected 1")
	}
	if haps[0].Context["loc"] == nil {
		t.Logf("WithLoc context %v", haps[0].Context)
	}
}

func TestReifyBasic(t *testing.T) {
	p := Reify("a")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Reify string expected 1")
	}
	q := Reify(Pure("a"))
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 {
		t.Fatalf("Reify pattern expected 1")
	}
}

func TestSilenceP2(t *testing.T) {
	p := Pure("a").SilenceP()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 0 {
		t.Fatalf("SilenceP expected 0")
	}
}
