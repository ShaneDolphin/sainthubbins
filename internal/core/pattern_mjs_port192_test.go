package core

import "testing"

func TestMJS_Port192_BinaryOpAddMulSubDivModPowFourth(t *testing.T) {
	p := Pure(2).Pow(Pure(3))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 8 {
		t.Fatalf("Pow 2^3=8")
	}
	q := Pure(7).Mod(Pure(3))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 1 {
		t.Fatalf("Mod 7 mod 3 is 1")
	}
	r := Pure(5).Band(Pure(3))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 1 {
		t.Fatalf("Band 5 & 3 =1")
	}
	s := Pure(5).Bor(Pure(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 7 {
		t.Fatalf("Bor 5 |2=7")
	}
}

func TestMJS_Port192_StructureWithSignalSineTriFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Struct(FastCat(Pure(true), Pure(false), Pure(true)))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Struct true false true 2")
	}
	q := Sine().Range(0, 1)
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Range 0,1")
	}
	r := Tri().Range(-1, 1)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Range -1,1")
	}
}

func TestMJS_Port192_EveryWhenOffSlowFastCompressFourth(t *testing.T) {
	p := Pure("a").Every(4, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(4))) == 0 {
		t.Fatalf("Every 4")
	}
	q := Pure("x").When(false, func(q Pattern) Pattern { return q.Rev() })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
	r := Pure("bd").Off(0.125, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125 <2")
	}
	s := Pure("c").Compress(FractionFromFloat(0), FractionFromFloat(0.5))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Compress 0-0.5")
	}
}
