package core

import "testing"

func TestMJS_ExpandRange2(t *testing.T) {
	p := Pure(0.5).Range(0, 100)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Range 0.5 0-100")
	}
	v := toFloat(haps[0].Value)
	if v < 49 || v > 51 {
		t.Fatalf("Range 0.5 0-100 expected ~50 got %v", v)
	}
}

func TestMJS_ArpPolychord2(t *testing.T) {
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

func TestMJS_StackWithSignal2(t *testing.T) {
	s := Saw().Range(0, 1)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Saw Range")
	}
	si := Sine().Range(0, 1)
	haps2 := si.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Sine Range")
	}
}
