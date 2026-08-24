package core

import "testing"

func TestMJS_Port243_InsideOutsideRevHurryFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Inside(FractionFromInt(2), func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Inside 2 Rev <2")
	}
	q := FastCat(Pure("a"), Pure("b")).Outside(FractionFromInt(2), func(x Pattern) Pattern { return x.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside 2 FastF2")
	}
	r := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := r.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Rev 3 got %d", len(haps))
	}
	h := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Hurry(2)
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 2")
	}
}

func TestMJS_Port243_EuclidBjorklundStructFourth(t *testing.T) {
	e := Pure("bd").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid 3,8")
	}
	b := Bjorklund(3, 8)
	cnt := 0
	for _, v := range b {
		if v != 0 {
			cnt++
		}
	}
	if cnt != 3 {
		t.Fatalf("Bjorklund 3 !=3 %d", cnt)
	}
	s := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(false), Pure(true), Pure(false)))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Struct t f true 2")
	}
}

func TestMJS_Port243_ControlsGainPanCutoffFourth(t *testing.T) {
	p := S("bd").Set(Gain(0.9))
	v := p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["gain"] != 0.9 {
		t.Fatalf("gain 0.9 got %v", v["gain"])
	}
	q := S("hh").Set(Pan(0.3))
	v2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v2["pan"] != 0.3 {
		t.Fatalf("pan 0.3 got %v", v2["pan"])
	}
	r := Stack(S("bd").Set(Cutoff(800)), S("sd").Set(Cutoff(1200)))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack cutoff 2")
	}
}
