package core

import "testing"

func TestMJS_Port185_FastCatSlowCatSequenceArpFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 4 {
		t.Fatalf("FastCat 4")
	}
	q := SlowCat(Pure("a"), Pure("b"), Pure("c"))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(3))) != 3 {
		t.Fatalf("SlowCat 3")
	}
	r := Sequence(Pure("a"), Pure("b"), Pure("c"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("Sequence 3")
	}
	s := Pure("c3 e3 g3").Arp("updown")
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp updown")
	}
}

func TestMJS_Port185_InsideOutsideWithSlowFastFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Inside(FractionFromInt(2), func(q Pattern) Pattern { return q.Slow(FractionFromInt(2)) })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Inside Slow2")
	}
	q := FastCat(Pure("a"), Pure("b")).Outside(FractionFromInt(2), func(x Pattern) Pattern { return x.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Outside FastF2")
	}
	s := Sine().Inside(FractionFromInt(4), func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Inside 4 FastF2")
	}
}

func TestMJS_Port185_ControlsNoteCutoffRoomFourth(t *testing.T) {
	n := Note("c3")
	if n.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["note"] != "c3" {
		t.Fatalf("note c3")
	}
	c := Cutoff(1200)
	if c.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["cutoff"] != 1200 {
		t.Fatalf("cutoff 1200")
	}
	r := Room(0.5)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value.(map[string]any)["room"] != 0.5 {
		t.Fatalf("room 0.5")
	}
	stack := Stack(Note("c3").Set(Cutoff(800)), Note("e3").Set(Room(0.3)))
	if len(stack.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack Note Cutoff Room 2")
	}
}
