package core

import "testing"

func TestMJS_Splice2(t *testing.T) {
	// Splice via Slice alias
	sl := Slice(2, Pure(0), Pure("a"))
	if len(sl.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Slice 2 expected non-empty")
	}
}

func TestMJS_PatternCFuncs2(t *testing.T) {
	c := Cat(Pure("a"), Pure("b"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("Cat 2 cycles")
	}
	st := Stack(Pure("a"), Pure("b"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
	si := Silence()
	if len(si.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence 0")
	}
}

func TestMJS_ArpMode2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"))
	arpUp := p.Arp("up")
	if len(arpUp.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up")
	}
	arpDown := p.Arp("down")
	if len(arpDown.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down")
	}
}
