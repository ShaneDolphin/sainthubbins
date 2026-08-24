package core

import "testing"

func TestMJS_HushGapSilence2(t *testing.T) {
	h := Pure("a").Hush()
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush 0")
	}
	g := Gap(2)
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("Gap 2 steps 2")
	}
	si := Silence()
	if len(si.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence 0")
	}
}

func TestMJS_DefragResetPoly2(t *testing.T) {
	p := Pure(Pure("a"))
	pj := p.PolyJoin()
	if len(pj.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("PolyJoin")
	}
	rj := p.ResetJoin()
	_ = rj
	d := Pure("a").Defragment()
	_ = d
}

func TestMJS_BinaryOps2(t *testing.T) {
	a := Pure(3).Band(Pure(1))
	if len(a.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Band")
	}
	// Range
	r := Pure(0.5).Range(0, 100)
	haps := r.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 || toFloat(haps[0].Value) < 49 {
		t.Fatalf("Range 0.5 0-100")
	}
}
