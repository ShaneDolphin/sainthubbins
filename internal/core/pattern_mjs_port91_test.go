package core

import "testing"

func TestMJS_StructAllMask2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	sa := p.StructAll(Pure(true))
	if len(sa.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("StructAll true")
	}
	ma := p.Mask(Pure(true))
	if len(ma.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true")
	}
}

func TestMJS_SometimesOften2(t *testing.T) {
	s := Pure("a").Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Sometimes")
	}
	sb75 := Pure("a").SometimesBy(0.75, func(p Pattern) Pattern { return p })
	if len(sb75.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.75")
	}
	sb25 := Pure("a").SometimesBy(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(sb25.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SometimesBy 0.25")
	}
}

func TestMJS_PlayCloned2(t *testing.T) {
	// Off and Echo
	o := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25")
	}
	e := Pure("a").Echo(2, FractionFromFloat(0.25), 0.5)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Echo 2 0.25 0.5")
	}
}
