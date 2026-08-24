package core

import "testing"

func TestMJS_Port119_EveryOffWhenChunk(t *testing.T) {
	e := Pure("bd").Every(3, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 3")
	}
	o := Pure("bd sd").Off(0.125, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125")
	}
	w := Pure("a b c").When(true, func(q Pattern) Pattern { return q.Rev() })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Rev")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(3, func(q Pattern) Pattern { return q.Rev() })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 3 Rev")
	}
}

func TestMJS_Port119_ScaleChordTranspose(t *testing.T) {
	p := Pure("c3").Scale("minor")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Scale minor")
	}
	q := Pure("c3").Chord("maj7")
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chord maj7")
	}
	// Chord stores chord in context, degree passthrough - just check non-empty
	r := Pure("c4").Transpose(2)
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Transpose 2")
	}
	s := Pure(60).Scale("major")
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Scale major on 60")
	}
}

func TestMJS_Port119_PatternAddMulDivMod(t *testing.T) {
	a := FastCat(Pure(2), Pure(4))
	b := a.Add(Pure(3))
	haps := b.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Add 3 len 2 got %d", len(haps))
	}
	m := FastCat(Pure(2), Pure(3)).Mul(Pure(2))
	if len(m.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Mul 2")
	}
	d := Pure(10).Div(Pure(2))
	if d.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Div")
	}
	mod := Pure(5).Mod(Pure(3))
	if mod.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mod")
	}
}
