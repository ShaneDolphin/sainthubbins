package core

import "testing"

func TestWhenBasic(t *testing.T) {
	p := Pure("a").When(true, func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value.(string) != "b" {
		t.Fatalf("When true expected b")
	}
	q := Pure("a").When(false, func(p Pattern) Pattern { return Pure("b") })
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 || haps2[0].Value.(string) != "a" {
		t.Fatalf("When false expected a")
	}
}

func TestEveryBasic(t *testing.T) {
	p := Pure("a").Every(2, func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value.(string) != "b" {
		t.Fatalf("Every 2 first cycle b")
	}
	haps2 := p.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(haps2) != 1 || haps2[0].Value.(string) != "a" {
		t.Fatalf("Every 2 second cycle a")
	}
}

func TestOffBasic2(t *testing.T) {
	p := Pure(1).Off(0.5, func(p Pattern) Pattern { return p.Add(Pure(1)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Off with Add expected >=2")
	}
}
