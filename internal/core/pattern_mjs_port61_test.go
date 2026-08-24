package core

import "testing"

func TestMJS_ValueSteps3(t *testing.T) {
	p := Pure("a")
	if p.Steps == nil || !p.Steps.Equals(FractionFromInt(1)) {
		t.Fatalf("Pure steps 1 expected got %v", p.Steps)
	}
	ws := p.WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(2)) })
	if ws.Steps == nil || !ws.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("WithSteps *2 expected 2 got %v", ws.Steps)
	}
	g := Gap(4)
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(4)) {
		t.Fatalf("Gap 4 steps 4")
	}
	si := Silence()
	if si.Steps != nil {
		t.Fatalf("Silence steps nil expected got %v", si.Steps)
	}
}

func TestMJS_HapContext3(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "n": 1})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Pure map expected 1")
	}
	// SetContext
	pc := Pure("a").SetContext(map[string]any{"orbit": 1})
	haps2 := pc.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 || haps2[0].Context["orbit"] != 1 {
		t.Fatalf("SetContext orbit 1 got %v", haps2[0].Context)
	}
}

func TestMJS_PatternTypes3(t *testing.T) {
	if !IsPattern(Pure("a")) {
		t.Fatalf("IsPattern Pure true")
	}
	if IsPattern("a") {
		t.Fatalf("IsPattern string false")
	}
	r := Reify("a")
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Reify string -> Pure")
	}
	r2 := Reify(Pure("b"))
	if len(r2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Reify pattern -> pattern")
	}
}
