package core

import "testing"

func TestMJS_Port147_BinaryOpAddMulSubDivModPow(t *testing.T) {
	p := FastCat(Pure(8), Pure(4))
	d := p.Div(Pure(2))
	if len(d.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Div 2")
	}
	m := Pure(3).Pow(Pure(2))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(float64) != 9 {
		t.Fatalf("Pow 2 ->9")
	}
	mo := Pure(10).Mod(Pure(4))
	if mo.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mod 4")
	}
	band := Pure(5).Band(Pure(3))
	if band.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Band 3")
	}
	bor := Pure(5).Bor(Pure(2))
	if bor.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Bor 2")
	}
}

func TestMJS_Port147_StructureWithSignalSineTri(t *testing.T) {
	p := Pure("a b c").Struct(FastCat(Pure(true), Pure(false), Pure(true)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf t")
	}
	s := Sine().Range(-1, 1).Struct(Pure(true))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Struct true")
	}
	tri := Tri().Range(0, 5).Struct(FastCat(Pure(true), Pure(true), Pure(false)))
	if tri.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Tri Struct ttf")
	}
}

func TestMJS_Port147_EveryWhenOffSlowFastCompress(t *testing.T) {
	e := Pure("bd").Every(4, func(q Pattern) Pattern { return q.Rev() })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 4 Rev")
	}
	w := Pure("bd sd").When(false, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("When false 1")
	}
	o := Pure("a").Off(0.125, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125")
	}
	c := Pure("a b").Compress(FractionFromFloat(0), FractionFromFloat(0.5)).FastF(FractionFromInt(2))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0-0.5 Fast 2")
	}
}
