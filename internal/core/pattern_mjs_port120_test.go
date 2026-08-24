package core

import "testing"

func TestMJS_Port120_FastCatSlowCatSequenceArp(t *testing.T) {
	fc := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(fc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastCat 4")
	}
	sc := SlowCat(Pure("a"), Pure("b"))
	if len(sc.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("SlowCat 2 per cycle 1")
	}
	seq := Sequence(Pure("bd"), Pure("sd"), Pure("hh"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Sequence 3")
	}
	arp := Pure("c3 e3 g3 b3").Arp("updown")
	if len(arp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
}

func TestMJS_Port120_InsideOutsideWithSlowFast(t *testing.T) {
	p := Pure("bd sd hh oh").Inside(2, func(q Pattern) Pattern { return q.Slow(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Inside Slow 2")
	}
	o := Pure("bd sd hh oh").Outside(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if o.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Outside FastF 2")
	}
	f := Sine().Inside(4, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if f.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Inside 4")
	}
}

func TestMJS_Port120_ControlsNoteCutoffRoom(t *testing.T) {
	n := Note("c3").FirstCycle()
	if n[0].Value.(map[string]any)["note"] != "c3" {
		t.Fatalf("Note c3")
	}
	c := Cutoff(1200).FirstCycle()[0].Value.(map[string]any)
	if c["cutoff"] != 1200 {
		t.Fatalf("Cutoff 1200")
	}
	// Room alias?
	r := Room(0.5).FirstCycle()[0].Value.(map[string]any)
	if r["room"] != 0.5 {
		t.Fatalf("Room 0.5")
	}
	s := Stack(Note("c4"), Gain(0.8), Cutoff(800))
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Stack Note+Gain+Cutoff 3")
	}
}
