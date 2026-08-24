package core

import "testing"

func TestMJS_Port169_StructWithBoolMaskThird(t *testing.T) {
	p := Pure("bd").Struct(FastCat(Pure(true), Pure(true), Pure(false), Pure(false)))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct tt ff")
	}
	q := Pure("a b").Mask(FastCat(Pure(true), Pure(false)))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Mask tf")
	}
	r := Pure("bd").Struct(Pure(false))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Struct false empty")
	}
}

func TestMJS_Port169_EveryOffWhenChunkThird(t *testing.T) {
	e := Pure("bd").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if e.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Every 2")
	}
	o := Pure("bd sd hh").Off(0.125, func(q Pattern) Pattern { return q.Rev() })
	if len(o.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Off 0.125 Rev")
	}
	w := Pure("a b c").When(true, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("When true Fast 2")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(4, func(pat Pattern) Pattern { return pat.Rev() })
	if ch.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 4 Rev")
	}
}

func TestMJS_Port169_PatternAddWithStructureThird(t *testing.T) {
	p := FastCat(Pure(1), Pure(2)).Add(Pure(5))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Add 5")
	}
	q := Pure("bd").Struct(Pure(true)).Add(Pure(0))
	if q.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Struct Add 0")
	}
	r := Stack(Pure(10), Pure(20)).Add(Pure(1))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack Add 1")
	}
}
