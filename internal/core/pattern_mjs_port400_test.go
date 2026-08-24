package core

import "testing"

func TestMJS_Port400_InsideOutsideRevHurryFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Inside(FractionFromInt(2), func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 { t.Fatalf("Inside") }
	q := FastCat(Pure("a"), Pure("b")).Outside(FractionFromInt(2), func(x Pattern) Pattern { return x.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Outside") }
	r := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 { t.Fatalf("Rev 3") }
	h := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Hurry(2)
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Hurry") }
}
func TestMJS_Port400_EuclidBjorklundStructFourth(t *testing.T) {
	e := Pure("bd").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Euclid") }
	b := Bjorklund(3, 8)
	c := 0; for _, v := range b { if v != 0 { c++ } }
	if c != 3 { t.Fatalf("Bjorklund") }
	s := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Struct") }
}
func TestMJS_Port400_ControlsGainPanCutoffFourth(t *testing.T) {
	p := S("bd").Set(Gain(0.9))
	v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["gain"] != 0.9 { t.Fatalf("gain") }
	q := S("hh").Set(Pan(0.3))
	v2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v2["pan"] != 0.3 { t.Fatalf("pan") }
	r := Stack(S("bd").Set(Cutoff(800)), S("sd").Set(Cutoff(1200)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Stack") }
}
