package core

import "testing"

func TestMJS_Port1107_SlowFastCompressZoomFourth(t *testing.T) {
	p := Pure("a b c").Slow(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 { t.Fatalf("Slow") }
	q := FastCat(Pure("x"), Pure("y")).FastF(FractionFromInt(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 { t.Fatalf("FastF2 4") }
	r := Pure("z").Compress(FractionFromFloat(0), FractionFromFloat(0.5))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Compress") }
	s := Pure("w").Zoom(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Zoom") }
}
func TestMJS_Port1107_EuclidSqueezeWithValueFourth(t *testing.T) {
	e := Pure("bd").Euclid(5, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Euclid") }
	sq := Pure(Pure("a")).SqueezeJoin()
	if len(sq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("SqueezeJoin") }
	p := Pure(3).WithValue(func(v any) any { return v.(int) + 2 })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value != 5 { t.Fatalf("WithValue") }
	q := S("bd").Set(Cutoff(800))
	v := q.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)
	if v["cutoff"] != 800 { t.Fatalf("cutoff") }
}
func TestMJS_Port1107_JuxRevBrakHurryFourth(t *testing.T) {
	p := Pure("a").JuxBy(0.25, func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 { t.Fatalf("JuxBy") }
	q := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 { t.Fatalf("Rev 3") }
	r := Pure("a b").Brak()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Brak") }
	h := Pure("a b c").Hurry(1.5)
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Hurry") }
}
