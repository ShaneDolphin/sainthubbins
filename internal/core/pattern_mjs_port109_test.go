package core

import "testing"

func TestMJS_StructWithBoolMask2(t *testing.T) {
	p := Pure("bd sd hh")
	s := p.Struct(Pure(true).FastF(FractionFromInt(4)))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Struct") }
	m := p.Mask(Pure(true).FastF(FractionFromInt(2)))
	if m.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Mask") }
}

func TestMJS_EveryOffWhenChunk2(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.Rev() })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Every") }
	w := Pure("a b").When(func(b bool) bool { return true }, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("When") }
}

func TestMJS_PatternAddWithStructure3(t *testing.T) {
	p := Stack(Pure(1), Pure(2))
	added := p.Add(Pure(10))
	haps := added.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)!=2 { t.Fatalf("Add 2 got %d", len(haps)) }
}
