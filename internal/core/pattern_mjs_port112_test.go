package core

import "testing"

func TestMJS_Port112_StackCatFastSlow(t *testing.T) {
	s := Stack(Pure("bd"), Pure("sd"))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	c := Cat(Pure("bd"), Pure("sd"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Cat SlowCat 1 per cycle")
	}
	f := Pure("bd").FastF(FractionFromInt(2))
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("FastF 2")
	}
	slow := Pure("bd sd").Slow(FractionFromInt(2))
	if slow.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Slow nil")
	}
}

func TestMJS_Port112_SignalRangePerlin(t *testing.T) {
	sine := Sine().Range(0, 1)
	haps := sine.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps == nil || len(haps) == 0 {
		t.Fatalf("Sine Range nil")
	}
	perlin := Perlin().Slow(FractionFromInt(2))
	if perlin.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow")
	}
	saw := Saw().Range(10, 20)
	if saw.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Saw Range")
	}
}

func TestMJS_Port112_ArpChunkDegrade(t *testing.T) {
	p := Pure("c3 e3 g3")
	arp := p.Arp("up")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	chunk := Pure("bd sd hh oh").Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if chunk.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Chunk 2")
	}
	deg := Pure("bd").DegradeBy(0)
	if len(deg.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0")
	}
	if len(Pure("bd").DegradeBy(1).QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("DegradeBy 1 should be 0")
	}
}
