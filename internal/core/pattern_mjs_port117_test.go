package core

import "testing"

func TestMJS_Port117_SlowFastCompressZoom(t *testing.T) {
	s := Pure("a b c").Slow(FractionFromInt(2))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow 2")
	}
	f := Pure("a").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2")
	}
	c := Pure("bd sd").Compress(FractionFromFloat(0), FractionFromFloat(0.5))
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Compress 0-0.5")
	}
	z := Pure("bd sd").Zoom(FractionFromFloat(0.25), FractionFromFloat(0.75))
	if z.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Zoom 0.25-0.75")
	}
}

func TestMJS_Port117_EuclidSqueezeWithValue(t *testing.T) {
	e := Pure("bd").Euclid(5, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Euclid 5,8")
	}
	sj := Pure(Pure("bd")).SqueezeJoin()
	if sj.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SqueezeJoin")
	}
	p := S("bd").WithValue(func(v any) any {
		m := v.(map[string]any); m["gain"] = 0.6; return m
	})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("S WithValue gain 0.6")
	}
}

func TestMJS_Port117_JuxRevBrakHurry(t *testing.T) {
	j := Pure("bd sd cp").JuxBy(0.25, func(q Pattern) Pattern { return q.Rev() })
	if len(j.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("JuxBy 0.25")
	}
	rev := Pure("a b c d").Rev()
	if len(rev.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rev 4")
	}
	if len(Pure("bd sd").Brak().QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Brak")
	}
	if len(Pure("bd").Hurry(1.5).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Hurry 1.5")
	}
}
