package core

import "testing"

func TestOffBasic(t *testing.T) {
	p := Pure("a").Off(0.25, func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Off expected >=2 got %d", len(haps))
	}
}

func TestOffWithAdd(t *testing.T) {
	p := Pure(1).Off(0.25, func(p Pattern) Pattern { return p.Add(Pure(1)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Off Add expected >=2")
	}
}
