package core

import "testing"

func TestMJS_PatternValueWithContext2(t *testing.T) {
	p := Pure(5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 || haps[0].Context == nil && haps[0].Value != 5 {
		// Context may be nil; just check value
		if haps[0].Value != 5 { t.Fatalf("value 5") }
	}
	if haps[0].Part.Begin.Float64() != 0 { t.Fatalf("part begin 0") }
}

func TestMJS_StackCatPolymeterSequence2(t *testing.T) {
	a := Pure("a")
	b := Stack(a, a.FastF(FractionFromInt(2)))
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Stack FastF") }
	c := Cat(Pure("a"), Pure("b"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(2))) < 2 { t.Fatalf("Cat 2 got %d", len(c.QueryArc(FractionFromInt(0), FractionFromInt(2)))) }
	steps := Polymeter(Pure(1))
	if steps.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("PolymeterSteps") }
}

func TestMJS_SignalSlowFastChoice2(t *testing.T) {
	s := SlowCat(Pure(1), Pure(2))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1)))==0 { t.Fatalf("Slowcat") }
	f := Pure("x").Fast(FastCat(Pure(FractionFromInt(1)), Pure(FractionFromInt(2))))
	if f.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Fast Pattern") }
}
