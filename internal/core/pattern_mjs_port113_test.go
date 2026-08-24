package core

import "testing"

func TestMJS_Port113_InsideOutsideRevHurry(t *testing.T) {
	p := Sine().Inside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Inside 2")
	}
	o := Sine().Outside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Outside 2")
	}
	rev := Pure("a b c").Rev()
	if len(rev.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rev")
	}
	if len(Pure("bd sd").Hurry(2).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 2")
	}
	if len(Pure("bd sd").Brak().QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
}

func TestMJS_Port113_EuclidBjorklundStruct(t *testing.T) {
	p := Pure("bd").Euclid(3, 8)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid 3,8")
	}
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("Bjorklund len 8 got %d", len(b))
	}
	s := Pure("a").Struct(Pure(true))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true")
	}
	s2 := Pure("a").Struct(Pure(false))
	if len(s2.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Struct false should be empty")
	}
}

func TestMJS_Port113_ControlsGainPanCutoff(t *testing.T) {
	h := S("bd").FirstCycle()
	if len(h) == 0 {
		t.Fatalf("S bd empty")
	}
	m, ok := h[0].Value.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Fatalf("S bd map")
	}
	g := Gain(0.5).FirstCycle()
	if g[0].Value.(map[string]any)["gain"] != 0.5 {
		t.Fatalf("Gain 0.5")
	}
	pan := Pan(0.7).FirstCycle()[0].Value.(map[string]any)
	if pan["pan"] != 0.7 {
		t.Fatalf("Pan 0.7")
	}
	stack := Stack(S("bd"), Gain(0.9))
	if len(stack.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack S+Gain 2")
	}
}
