package core

import "testing"

func TestMJS_ClockDispose2(t *testing.T) {
	c := NewClock(0.5)
	if c.CPS != 0.5 {
		t.Fatalf("Clock CPS 0.5 got %v", c.CPS)
	}
	c.SetCPS(1.0)
	if c.CPS != 1.0 {
		t.Fatalf("SetCPS 1.0 got %v", c.CPS)
	}
	if c.Duration <= 0 {
		t.Fatalf("Duration >0")
	}
	_ = c.Interval
}

func TestMJS_PatternAddSub2(t *testing.T) {
	p := Pure(5).Add(Pure(7))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || toFloat(haps[0].Value) != 12 {
		t.Fatalf("Add 5+7=12 got %v", haps[0].Value)
	}
	s := Pure(10).Sub(Pure(4))
	haps2 := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps2[0].Value) != 6 {
		t.Fatalf("Sub 10-4=6 got %v", haps2[0].Value)
	}
}

func TestMJS_RangeSignal2(t *testing.T) {
	s := Saw().Range(10, 20)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Saw Range 10-20 expected non-empty")
	}
	// Check range within 10-20
	for _, h := range haps {
		v := toFloat(h.Value)
		if v < 9 || v > 21 {
			t.Fatalf("Saw Range 10-20 got %v", v)
		}
	}
}
