package core

import "testing"

func TestMJS_Port159_StructWithBoolMaskSecond(t *testing.T) {
	p := Pure("bd").Struct(FastCat(Pure(true), Pure(false), Pure(true)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tf t")
	}
	q := Pure("a").Mask(FastCat(Pure(true), Pure(true), Pure(false)))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tt f")
	}
	r := Pure("a b c").Struct(Pure(true))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true")
	}
}

func TestMJS_Port159_EveryOffWhenChunkSecond(t *testing.T) {
	e := Pure("bd").Every(3, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 3")
	}
	o := Pure("bd").Off(0.25, func(q Pattern) Pattern { return q.Rev() })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.25 Rev")
	}
	w := Pure("a").When(true, func(q Pattern) Pattern { return q.Add(Pure(1)) })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("When true Add 1")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2 Fast 2")
	}
}

func TestMJS_Port159_PatternAddWithStructureSecond(t *testing.T) {
	p := Stack(Pure(10), Pure(20))
	q := p.Add(Pure(5))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack Add 5 2")
	}
	r := FastCat(Pure(1), Pure(2), Pure(3)).Add(Pure(10))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("FastCat Add 10 3")
	}
	s := Pure("a").Struct(Pure(true)).Add(Pure(0))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct Add 0")
	}
}
